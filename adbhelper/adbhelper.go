package adbhelper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/cheggaaa/pb/v3"
	"github.com/mholt/archiver"
)

type AdbHelper struct {
	adbPath string
	adbExec string
}

func New(adbPath string) *AdbHelper {
	return &AdbHelper{adbPath: adbPath, adbExec: adbPath + "/adb"}
}

func (adbHelper *AdbHelper) SendCommand(command string) error {
	cmd := exec.Command(adbHelper.adbExec, command)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (adbHelper *AdbHelper) SendCommandWithResults(command string) (string, error) {
	cmd := exec.Command(adbHelper.adbExec, command)
	output, err := cmd.Output()
	if err != nil {

		return "", err
	}
	return string(output), nil
}

func (adbHelper *AdbHelper) SetupADB(adbUrl string) error {
	if _, err := os.Stat(adbHelper.adbPath); err == nil {
		return nil
	}
	var adbURL string
	switch runtime.GOOS {
	case "windows":
		adbURL = adbUrl + "windows.zip"
	case "linux":
		adbURL = adbUrl + "linux.zip"
	case "darwin":
		adbURL = adbUrl + "darwin.zip"
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	// Download ADB
	//plataform-tools is the default folder provided by google
	if err := downloadFile("platform-tools.zip", adbURL); err != nil {
		return err
	}

	// Extract ADB
	if err := archiver.Unarchive("platform-tools.zip", "."); err != nil {
		return err
	}

	os.Rename("platform-tools", adbHelper.adbPath)

	// Make adb executable on Unix systems
	if runtime.GOOS != "windows" {
		if err := os.Chmod(adbHelper.adbPath, 0755); err != nil {
			return err
		}
	}
	os.Remove("platform-tools.zip")
	return nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	bar := pb.Full.Start64(resp.ContentLength)
	barReader := bar.NewProxyReader(resp.Body)
	_, err = io.Copy(out, barReader)
	bar.Finish()
	if err != nil {
		return err
	}
	return nil
}
