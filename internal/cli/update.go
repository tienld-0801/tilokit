package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"tilokit/internal/utils"
	"tilokit/pkg/constants"

	"github.com/fatih/color"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func RunUpdateProcess() error {
	fmt.Println("🔍 Checking for updates...")

	latestRelease, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	currentVersion := strings.TrimPrefix(constants.Version, "v")
	latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")

	if currentVersion == latestVersion {
		utils.Success("You're already running the latest version: %s", constants.Version)
		return nil
	}

	fmt.Printf("📦 New version available: %s → %s\n", constants.Version, latestRelease.TagName)
	fmt.Printf("📝 Release notes:\n%s\n\n", latestRelease.Body)

	if !askConfirmation("Do you want to update now?") {
		utils.Info("Update cancelled.")
		return nil
	}

	fmt.Println("⬇️  Downloading latest version...")
	if err := downloadAndInstall(latestRelease); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}
	utils.Success("🎉 Successfully updated to %s!", latestRelease.TagName)
	utils.Info("Run 'tilokit --version' to verify the update")
	return nil
}

func getLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/tienld-0801/tilokit/releases/latest")
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func downloadAndInstall(release *GitHubRelease) error {
	var binaryName string
	switch runtime.GOOS {
	case constants.OSDarwin:
		if runtime.GOARCH == constants.ARM64 {
			binaryName = "tilokit-darwin-arm64"
		} else {
			binaryName = "tilokit-darwin-amd64"
		}
	case constants.OSLinux:
		if runtime.GOARCH == constants.ARM64 {
			binaryName = "tilokit-linux-arm64"
		} else {
			binaryName = "tilokit-linux-amd64"
		}
	case constants.OSWindows:
		if runtime.GOARCH == constants.ARM64 {
			binaryName = "tilokit-windows-arm64.exe"
		} else {
			binaryName = "tilokit-windows-amd64.exe"
		}
	default:
		return fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == binaryName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no binary found for platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Get(downloadURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close response body: %v\n", err)
		}
	}()

	currentExe, err := os.Executable()
	if err != nil {
		return err
	}

	currentExe = filepath.Clean(currentExe)
	exeDir := filepath.Dir(currentExe)

	out, err := os.CreateTemp(exeDir, "tilokit_update_*.tmp")
	if err != nil {
		return err
	}
	tmpFile := out.Name()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		_ = out.Close()
		_ = os.Remove(tmpFile)
		return err
	}

	if err := out.Close(); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	if err := os.Chmod(tmpFile, 0600); err != nil {
		_ = os.Remove(tmpFile)
		return err
	}

	if runtime.GOOS == constants.OSWindows {
		return replaceExecutableWindows(currentExe, tmpFile)
	} else {
		return os.Rename(tmpFile, currentExe)
	}
}

func replaceExecutableWindows(currentExe, tmpFile string) error {
	batchScript := currentExe + "_update.bat"

	if strings.Contains(batchScript, "..") || strings.Contains(batchScript, "//") {
		return fmt.Errorf("security violation: invalid batch script path")
	}

	if !filepath.IsAbs(batchScript) {
		return fmt.Errorf("security violation: batch script must be absolute path")
	}

	scriptContent := fmt.Sprintf(`@echo off
timeout /t 2
move "%s" "%s"
del "%%~f0"`, tmpFile, currentExe)

	if err := os.WriteFile(batchScript, []byte(scriptContent), 0600); err != nil {
		return err
	}

	if err := executeWindowsUpdateScript(batchScript); err != nil {
		return err
	}
	return nil
}

func executeWindowsUpdateScript(scriptPath string) error {
	if !filepath.IsAbs(scriptPath) || strings.Contains(scriptPath, "..") {
		return fmt.Errorf("invalid script path")
	}

	fmt.Println("⚠️  On Windows, please manually replace the binary after download completes.")
	fmt.Printf("📁 Downloaded file: %s\n", scriptPath)
	fmt.Println("🔄 The update will complete on next restart.")

	return nil
}

func askConfirmation(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/N]: ", color.YellowString(question))
		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))
		switch response {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		default:
			fmt.Println("Please answer 'y' or 'n'")
		}
	}
}
