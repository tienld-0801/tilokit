package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/ti-lo/tilokit/internal/utils"
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// RunUpdateProcess handles the update logic
func RunUpdateProcess() error {
	// No banner for update command - only init has banner
	fmt.Println("🔍 Checking for updates...")

	// Get latest release info
	latestRelease, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	// Compare versions
	currentVersion := strings.TrimPrefix(Version, "v")
	latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")

	if currentVersion == latestVersion {
		utils.Success("You're already running the latest version: %s", Version)
		return nil
	}

	fmt.Printf("📦 New version available: %s → %s\n", Version, latestRelease.TagName)
	fmt.Printf("📝 Release notes:\n%s\n\n", latestRelease.Body)

	// Ask for confirmation
	if !askConfirmation("Do you want to update now?") {
		utils.Info("Update cancelled.")
		return nil
	}

	// Download and install
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
	case "darwin":
		if runtime.GOARCH == "arm64" {
			binaryName = "tilokit-darwin-arm64"
		} else {
			binaryName = "tilokit-darwin-amd64"
		}
	case "linux":
		if runtime.GOARCH == "arm64" {
			binaryName = "tilokit-linux-arm64"
		} else {
			binaryName = "tilokit-linux-amd64"
		}
	case "windows":
		if runtime.GOARCH == "arm64" {
			binaryName = "tilokit-windows-arm64.exe"
		} else {
			binaryName = "tilokit-windows-amd64.exe"
		}
	default:
		return fmt.Errorf("unsupported platform: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	// Find the download URL for our platform
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

	// Download the binary
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

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return err
	}

	// Create temporary file
	tmpFile := currentExe + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := out.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", err)
		}
	}()

	// Copy downloaded content
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		// #nosec G104 - cleanup error ignored
		_ = os.Remove(tmpFile)
		return err
	}

	// Make executable
	// #nosec G302 - executables need 0700 permissions
	if err := os.Chmod(tmpFile, 0700); err != nil {
		// #nosec G104 - cleanup error ignored
		_ = os.Remove(tmpFile)
		return err
	}

	// Replace current executable
	if runtime.GOOS == "windows" {
		// On Windows, we can't replace a running executable
		// So we'll use a helper script approach
		return replaceExecutableWindows(currentExe, tmpFile)
	} else {
		// On Unix systems, we can replace the file
		return os.Rename(tmpFile, currentExe)
	}
}

func replaceExecutableWindows(currentExe, tmpFile string) error {
	// Create batch script for replacement
	batchScript := currentExe + "_update.bat"
	scriptContent := fmt.Sprintf(`@echo off
timeout /t 2
move "%s" "%s"
del "%%~f0"`, tmpFile, currentExe)

	// Write batch script
	if err := os.WriteFile(batchScript, []byte(scriptContent), 0600); err != nil {
		return err
	}

	// Execute batch script
	// #nosec G204 - safe script content
	cmd := exec.Command("cmd", "/C", "start", "/B", batchScript)
	return cmd.Start()
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
