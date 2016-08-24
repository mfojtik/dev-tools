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
	"github.com/mfojtik/dev-tools/pkg/util"
)

func RebuildOpenShiftImage(name string) error {
	tmpDir, err := ioutils.TempDir("", "otp-rebuild-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	if err := pullIfNotExists(name); err != nil {
		return err
	}
	if err := locateAndCopyOpenShiftBinary(tmpDir); err != nil {
		return err
	}
	layerID, err := findBinaryCopyLayer(name)
	if err != nil {
		return err
	}
	tmpBaseImageName := fmt.Sprintf("base-%d", time.Now().Unix())
	util.Debugf("tagging image %q as %q", layerID, tmpBaseImageName)
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

func pullIfNotExists(image string) error {
	_, err := exec.Command("docker", "inspect", image).CombinedOutput()
	if err == nil {
		return nil
	}
	util.Infof("pulling image %q ...", image)
	_, err = exec.Command("docker", "pull", image).CombinedOutput()
	return err
}

func rebuildImage(baseImage, imageName, buildDir string) error {
	tmpImageName := fmt.Sprintf("image-%d", time.Now().Unix())
	util.Debugf("building %q", tmpImageName)
	cmd := exec.Command("docker", "build", "-t", tmpImageName, ".")
	cmd.Dir = buildDir
	if out, err := cmd.Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	util.Debugf("removing image %q", baseImage)
	if out, err := exec.Command("docker", "rmi", "-f", baseImage).Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	util.Debugf("removing image %q", imageName)
	if out, err := exec.Command("docker", "rmi", "-f", imageName).Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	util.Debugf("tagging image %q as %q", tmpImageName, imageName)
	if out, err := exec.Command("docker", "tag", tmpImageName, imageName).Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	util.Debugf("removing image %q", tmpImageName)
	if out, err := exec.Command("docker", "rmi", "-f", tmpImageName).Output(); err != nil {
		return fmt.Errorf("%s (%v)", string(out), err)
	}
	return nil
}

func generateDockerfile(from, dstDir string) error {
	content := []string{
		fmt.Sprintf("FROM %s", from),
		"ADD openshift /usr/bin/openshift",
	}
	outputFile := filepath.Join(dstDir, "Dockerfile")
	err := ioutil.WriteFile(outputFile, []byte(strings.Join(content, "\n")+"\n"), 0666)
	if err != nil {
		return err
	}
	util.Debugf("written Dockerfile to %q", outputFile)
	return nil
}

func findBinaryCopyLayer(name string) (string, error) {
	util.Debugf("searching for the ADD openshift layer in %q image", name)
	historyOut, err := exec.Command("docker", "history", "--no-trunc", name).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %v", string(historyOut), err)
	}
	result := ""
	for _, line := range strings.Split(string(historyOut), "\n") {
		if strings.Contains(line, "ADD") && strings.Contains(line, "/usr/bin/openshift") {
			parts := strings.Split(line, " ")
			if len(parts) > 0 {
				result = strings.TrimSpace(parts[0])
			}
			util.Debugf("found ADD layer in %q", result)
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
	util.Debugf("located openshift binary in %q", binaryPath)
	if cpResult, err := exec.Command("cp", "-f", binaryPath, dst).Output(); err != nil {
		return fmt.Errorf("%s: %v", string(cpResult), err)
	}
	util.Debugf("copied openshift binary to %q", dst)
	return nil
}
