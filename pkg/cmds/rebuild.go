package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/pkg/ioutils"
)

func RebuildOpenShiftImage(name string) error {
	tmpDir, err := ioutils.TempDir("", "otp-rebuild-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	if err := locateAndCopyOpenShiftBinary(tmpDir); err != nil {
		return err
	}
	layerID, err := findBinaryCopyLayer(name)
	if err != nil {
		return err
	}
	tmpBaseImageName := fmt.Sprintf("base-%d", time.Now().Unix())
	if out, err := exec.Command("docker", "tag", layerID, tmpBaseImageName).Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	if err := generateDockerfile(tmpBaseImageName, tmpDir); err != nil {
		return err
	}
	if err := rebuildImage(tmpBaseImageName, name, tmpDir); err != nil {
		return err
	}
	return nil
}

func rebuildImage(baseImage, imageName, buildDir string) error {
	tmpImageName := fmt.Sprintf("image-%d", time.Now().Unix())
	cmd := exec.Command("docker", "build", "-t", tmpImageName, ".")
	cmd.Dir = buildDir
	if out, err := cmd.Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	if out, err := exec.Command("docker", "rmi", baseImage).Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	if out, err := exec.Command("docker", "rmi", imageName).Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	if out, err := exec.Command("docker", "tag", tmpImageName, imageName).Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	if out, err := exec.Command("docker", "rmi", tmpImageName).Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	return nil
}

func generateDockerfile(from, dstDir string) error {
	content := []string{
		fmt.Sprintf("FROM %s", from),
		"COPY openshift /usr/bin/openshift",
	}
	outputFile := filepath.Join(dstDir, "Dockerfile")
	err := ioutil.WriteFile(outputFile, []byte(strings.Join(content, "\n")+"\n"), 0666)
	if err != nil {
		return err
	}
	return nil
}

func findBinaryCopyLayer(name string) (string, error) {
	historyOut, err := exec.Command("docker", "history", "--no-trunc", name).Output()
	if err != nil {
		return "", fmt.Errorf("%s: %v", string(historyOut), err)
	}
	result := ""
	for _, line := range strings.Split(string(historyOut), "\n") {
		if strings.Contains(line, "COPY") && strings.Contains(line, "/usr/bin/openshift") {
			parts := strings.Split(line, " ")
			if len(parts) > 0 {
				result = strings.TrimSpace(parts[0])
			}
		}
	}
	if len(result) == 0 {
		return "", fmt.Errorf("unable to detect COPY layer in %q", name)
	}
	return result, nil
}

func locateAndCopyOpenShiftBinary(dst string) error {
	binaryPath, err := exec.LookPath("openshift")
	if err != nil {
		return err
	}
	if cpResult, err := exec.Command("cp", "-f", binaryPath, dst).Output(); err != nil {
		return fmt.Errorf("%s: %v", string(cpResult), err)
	}
	return nil
}
