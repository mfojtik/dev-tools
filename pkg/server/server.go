package server

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type State string

const (
	_             = iota
	Stopped State = "stopped"
	Running State = "running"
)

var binaryPath, _ = exec.LookPath("openshift")

type Server struct {
	Address       string
	WorkingDir    string
	CurrentState  State
	stop, stopped chan struct{}
}

func NewOpenShift(netInterface, dir string) (*Server, error) {
	if len(binaryPath) == 0 {
		return nil, fmt.Errorf("unable to locate 'openshift' binary in PATH")
	}
	if _, err := os.Stat(dir); err != nil {
		return nil, fmt.Errorf("the working directory %q must exist", dir)
	}
	addr, err := exec.Command("ip", "-o", "-4", "addr", "list", netInterface).Output()
	if err != nil {
		return nil, fmt.Errorf("unable to get IP address from %q interface: %v", netInterface, err)
	}
	if len(addr) < 6 {
		return nil, fmt.Errorf("parsing IP address failed from %q", string(addr))
	}
	ipAddr := strings.Split(strings.Split(string(addr), " ")[6], "/")[0]
	stop := make(chan struct{}, 1)
	stopped := make(chan struct{}, 1)
	return &Server{Address: ipAddr, WorkingDir: dir, CurrentState: Stopped, stop: stop, stopped: stopped}, nil
}

func (s *Server) Execute(cmds ...string) error {
	cmd := exec.Command("sudo", cmds...)
	cmd.Dir = s.WorkingDir
	return cmd.Run()
}

func (s *Server) Stop() {
	close(s.stop)
	<-s.stopped
}

func (s *Server) Run() (*Server, error) {
	// Refresh directories
	s.Execute("rm", "-rf", s.WorkingDir)
	s.Execute("mkdir", "-p", s.WorkingDir)

	path, _ := exec.LookPath("openshift")
	c := exec.Command("sudo", path, "start", "--master",
		fmt.Sprintf("https://%s:8443", s.Address), "--etcd-dir", "etcd", "--latest-images",
		"--volume-dir", "volumes",
	)
	c.Env = append(c.Env, "PATH="+os.Getenv("PATH"))
	c.Dir = s.WorkingDir
	go func() {
		defer func() { s.CurrentState = Stopped }()
		if err := c.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Server failed to start: %v\n", err)
		}
		s.CurrentState = Running
		if err := c.Wait(); err != nil {
			if strings.Contains(err.Error(), "143") {
				return
			}
			fmt.Fprintf(os.Stderr, "Server failed to run: %v\n", err)
		}
	}()
	go func() {
		<-s.stop
		fmt.Fprintf(os.Stdout, "Stopping server ...\n")
		if _, err := exec.Command("sudo", "kill", fmt.Sprintf("%d", c.Process.Pid)).Output(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to terminate server: %v\n", err)
		}
		c.Process.Wait()
		close(s.stopped)
	}()
	fmt.Fprintf(os.Stdout, "Waiting for server to be available ...\n")
	if err := waitForURL(fmt.Sprintf("https://%s:8443/healthz", s.Address)); err != nil {
		close(s.stop)
		return nil, err
	}
	return s, nil
}

func (s *Server) FixPermissions() {
	s.Execute("chmod", "a+rwX", "openshift.local.config/master/admin.kubeconfig")
	s.Execute("chmod", "+r", "openshift.local.config/master/openshift-registry.kubeconfig")
	s.Execute("chmod", "+r", "openshift.local.config/master/openshift-router.kubeconfig")
}

func waitForURL(url string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	for counter := 120; counter >= 0; counter-- {
		req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
		if err != nil {
			return err
		}
		client := &http.Client{Transport: tr}
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timeout while waiting for server up")
}
