package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	flutterExeName     = "flutter"
	flutterExeNameWin  = "flutter.bat"
	windowsOS          = "windows"
)

// FlutterSDKFinder handles Flutter SDK detection and validation
type FlutterSDKFinder struct{}

// FindFlutterSDK locates the Flutter SDK installation
func (f *FlutterSDKFinder) FindFlutterSDK() (string, error) {
	// First try to get Flutter SDK path from flutter doctor (most reliable)
	if sdkPath := f.getFlutterSDKFromDoctor(); sdkPath != "" {
		if f.isValidFlutterSDK(sdkPath) {
			return sdkPath, nil
		}
	}

	// Check environment variable
	if flutterRoot := os.Getenv("FLUTTER_ROOT"); flutterRoot != "" {
		if f.isValidFlutterSDK(flutterRoot) {
			return flutterRoot, nil
		}
	}

	// Check PATH for flutter executable
	if flutterPath := f.findFlutterInPath(); flutterPath != "" {
		sdkPath := f.extractSDKFromExecutable(flutterPath)
		if f.isValidFlutterSDK(sdkPath) {
			return sdkPath, nil
		}
	}

	// Check common installation locations
	commonPaths := f.getCommonFlutterPaths()
	for _, path := range commonPaths {
		if f.isValidFlutterSDK(path) {
			return path, nil
		}
	}

	return "", fmt.Errorf("flutter SDK not found. Please ensure Flutter is installed and either:\n1. Add Flutter to your PATH\n2. Set FLUTTER_ROOT environment variable\n3. Install Flutter in a standard location")
}

// findFlutterInPath searches for flutter executable in PATH
func (f *FlutterSDKFinder) findFlutterInPath() string {
	flutterExe := flutterExeName
	if runtime.GOOS == windowsOS {
		flutterExe = flutterExeNameWin
	}

	path, err := exec.LookPath(flutterExe)
	if err != nil {
		return ""
	}
	return path
}

// extractSDKFromExecutable extracts SDK path from flutter executable path
func (f *FlutterSDKFinder) extractSDKFromExecutable(execPath string) string {
	// Flutter executable is typically in bin/flutter within the SDK
	binDir := filepath.Dir(execPath)
	sdkPath := filepath.Dir(binDir)
	return sdkPath
}

// getCommonFlutterPaths returns common Flutter installation locations
func (f *FlutterSDKFinder) getCommonFlutterPaths() []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	paths := []string{}
	
	// macOS
	if runtime.GOOS == "darwin" {
		paths = append(paths, 
			filepath.Join(homeDir, "flutter"),
			filepath.Join(homeDir, "development", "flutter"),
			"/usr/local/flutter",
			"/opt/flutter",
		)
	}
	
	// Linux
	if runtime.GOOS == "linux" {
		paths = append(paths,
			filepath.Join(homeDir, "flutter"),
			filepath.Join(homeDir, "development", "flutter"),
			"/usr/local/flutter",
			"/opt/flutter",
		)
	}
	
	// Windows
	if runtime.GOOS == "windows" {
		paths = append(paths,
			filepath.Join(homeDir, "flutter"),
			filepath.Join(homeDir, "development", "flutter"),
			"C:\\flutter",
			"C:\\src\\flutter",
		)
	}
	
	return paths
}

// getFlutterSDKFromDoctor extracts Flutter SDK path from flutter doctor output
func (f *FlutterSDKFinder) getFlutterSDKFromDoctor() string {
	// First try machine-readable output
	if sdkPath := f.tryFlutterDoctorMachineForSDK(); sdkPath != "" {
		return sdkPath
	}
	
	// Fallback to verbose output parsing
	return f.tryFlutterDoctorVerboseForSDK()
}

// findFlutterExecutable finds the Flutter executable path
func (f *FlutterSDKFinder) findFlutterExecutable() string {
	// Try PATH first
	if flutterPath := f.findFlutterInPath(); flutterPath != "" {
		return flutterPath
	}
	
	// Try common installation locations
	commonPaths := f.getCommonFlutterPaths()
	for _, sdkPath := range commonPaths {
		if f.isValidFlutterSDK(sdkPath) {
			flutterExe := filepath.Join(sdkPath, "bin", flutterExeName)
			if runtime.GOOS == windowsOS {
				flutterExe = filepath.Join(sdkPath, "bin", flutterExeNameWin)
			}
			if _, err := os.Stat(flutterExe); err == nil {
				return flutterExe
			}
		}
	}
	
	return ""
}

// tryFlutterDoctorMachineForSDK tries to get Flutter SDK from machine-readable flutter doctor output
func (f *FlutterSDKFinder) tryFlutterDoctorMachineForSDK() string {
	// Note: --machine flag doesn't work with flutter doctor
	// We'll use the verbose output instead
	return ""
}

// tryFlutterDoctorVerboseForSDK tries to get Flutter SDK from verbose flutter doctor output
func (f *FlutterSDKFinder) tryFlutterDoctorVerboseForSDK() string {
	// Find Flutter executable
	flutterExe := f.findFlutterExecutable()
	if flutterExe == "" {
		return ""
	}
	
	// nolint:gosec // flutter executable path is validated and controlled
	cmd := exec.Command(flutterExe, "doctor", "-v")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Look for Flutter SDK path patterns
		// Pattern: "• Flutter version X.X.X on channel stable at /path/to/flutter"
		if strings.Contains(line, "Flutter version") && strings.Contains(line, "at") {
			// Extract path after "at "
			if atIndex := strings.Index(line, "at "); atIndex != -1 {
				pathPart := line[atIndex+3:] // Skip "at "
				// Remove any trailing text after the path
				if spaceIndex := strings.Index(pathPart, " "); spaceIndex > 0 {
					pathPart = pathPart[:spaceIndex]
				}
				if f.isValidFlutterSDK(pathPart) {
					return pathPart
				}
			}
		}
		
		// Also check for other patterns
		patterns := []string{
			"Flutter SDK at ",
			"Flutter SDK location: ",
			"FLUTTER_ROOT: ",
		}
		
		for _, pattern := range patterns {
			if strings.Contains(line, pattern) {
				parts := strings.Split(line, pattern)
				if len(parts) > 1 {
					sdkPath := strings.TrimSpace(parts[1])
					// Remove any trailing text after the path
					if spaceIndex := strings.Index(sdkPath, " "); spaceIndex > 0 {
						sdkPath = sdkPath[:spaceIndex]
					}
					if f.isValidFlutterSDK(sdkPath) {
						return sdkPath
					}
				}
			}
		}
	}
	
	return ""
}

// GetFlutterSDKFromDoctor extracts Flutter SDK path from flutter doctor output
func (f *FlutterSDKFinder) GetFlutterSDKFromDoctor() string {
	return f.getFlutterSDKFromDoctor()
}

// FindFlutterInPath searches for flutter executable in PATH
func (f *FlutterSDKFinder) FindFlutterInPath() string {
	return f.findFlutterInPath()
}

// ExtractSDKFromExecutable extracts SDK path from flutter executable path
func (f *FlutterSDKFinder) ExtractSDKFromExecutable(execPath string) string {
	return f.extractSDKFromExecutable(execPath)
}

// GetCommonFlutterPaths returns common Flutter installation locations
func (f *FlutterSDKFinder) GetCommonFlutterPaths() []string {
	return f.getCommonFlutterPaths()
}

// IsValidFlutterSDK checks if the path contains a valid Flutter SDK
func (f *FlutterSDKFinder) IsValidFlutterSDK(path string) bool {
	return f.isValidFlutterSDK(path)
}

// isValidFlutterSDK checks if the path contains a valid Flutter SDK
func (f *FlutterSDKFinder) isValidFlutterSDK(path string) bool {
	if path == "" {
		return false
	}
	
	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	
	// Check for key Flutter SDK files and directories
	requiredItems := []string{
		"bin",
		"packages",
		"version",
	}
	
	for _, item := range requiredItems {
		itemPath := filepath.Join(path, item)
		if _, err := os.Stat(itemPath); os.IsNotExist(err) {
			return false
		}
	}
	
	return true
}

// GetFlutterInfo returns comprehensive Flutter SDK information
func (f *FlutterSDKFinder) GetFlutterInfo(sdkPath string) map[string]string {
	info := make(map[string]string)
	
	if sdkPath == "" {
		return info
	}
	
	// Get Flutter version
	if version, err := f.getFlutterVersion(sdkPath); err == nil {
		info["version"] = version
	}
	
	// Get Flutter channel
	if channel, err := f.getFlutterChannel(sdkPath); err == nil {
		info["channel"] = channel
	}
	
	// Get Dart version
	if dartVersion, err := f.getDartVersion(sdkPath); err == nil {
		info["dart_version"] = dartVersion
	}
	
	info["sdk_path"] = sdkPath
	return info
}

// getFlutterVersion gets Flutter version from SDK
func (f *FlutterSDKFinder) getFlutterVersion(sdkPath string) (string, error) {
	versionFile := filepath.Join(sdkPath, "version")
	// nolint:gosec // version file path is validated and controlled
	content, err := os.ReadFile(versionFile)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

// getFlutterChannel gets Flutter channel from SDK
func (f *FlutterSDKFinder) getFlutterChannel(sdkPath string) (string, error) {
	// nolint:gosec // flutter executable path is validated and controlled
	cmd := exec.Command(filepath.Join(sdkPath, "bin", "flutter"), "channel")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*") {
			// Extract channel name from "* channel_name" format
			parts := strings.Fields(line)
			if len(parts) > 1 {
				return parts[1], nil
			}
		}
	}
	
	return "unknown", nil
}

// getDartVersion gets Dart version from Flutter SDK
func (f *FlutterSDKFinder) getDartVersion(sdkPath string) (string, error) {
	// nolint:gosec // dart executable path is validated and controlled
	cmd := exec.Command(filepath.Join(sdkPath, "bin", "dart"), "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	// Parse "Dart version X.X.X" format
	outputStr := string(output)
	if strings.Contains(outputStr, "Dart version") {
		parts := strings.Fields(outputStr)
		for i, part := range parts {
			if part == "version" && i+1 < len(parts) {
				return parts[i+1], nil
			}
		}
	}
	
	return "unknown", nil
}

// GetAndroidSDKFromFlutterDoctor uses flutter doctor to find Android SDK
func (f *FlutterSDKFinder) GetAndroidSDKFromFlutterDoctor(sdkPath string) string {
	// First try machine-readable output
	if androidSDK := f.tryFlutterDoctorMachine(sdkPath); androidSDK != "" {
		return androidSDK
	}
	
	// Fallback to verbose output parsing
	return f.tryFlutterDoctorVerbose(sdkPath)
}

// tryFlutterDoctorMachine tries to get Android SDK from machine-readable flutter doctor output
func (f *FlutterSDKFinder) tryFlutterDoctorMachine(sdkPath string) string {
	flutterExe := filepath.Join(sdkPath, "bin", "flutter")
	// nolint:gosec // flutter executable path is validated and controlled
	cmd := exec.Command(flutterExe, "doctor", "--machine")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	var doctorOutput []map[string]interface{}
	if err := json.Unmarshal(output, &doctorOutput); err != nil {
		return ""
	}
	
	// Look for Android toolchain information
	for _, item := range doctorOutput {
		if itemType, ok := item["type"].(string); ok && itemType == "android" {
			if statusInfo, ok := item["statusInfo"].(map[string]interface{}); ok {
				if message, exists := statusInfo["message"].(string); exists {
					// Extract Android SDK path from message
					if sdkPath := f.extractAndroidSDKPath(message); sdkPath != "" {
						return sdkPath
					}
				}
			}
		}
	}
	
	return ""
}

// tryFlutterDoctorVerbose tries to get Android SDK from verbose flutter doctor output
func (f *FlutterSDKFinder) tryFlutterDoctorVerbose(sdkPath string) string {
	flutterExe := filepath.Join(sdkPath, "bin", "flutter")
	// nolint:gosec // flutter executable path is validated and controlled
	cmd := exec.Command(flutterExe, "doctor", "-v")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Look for Android SDK path patterns
		patterns := []string{
			"Android SDK at ",
			"Android SDK location: ",
			"ANDROID_SDK_ROOT: ",
			"ANDROID_HOME: ",
		}
		
		for _, pattern := range patterns {
			if strings.Contains(line, pattern) {
				parts := strings.Split(line, pattern)
				if len(parts) > 1 {
					sdkPath := strings.TrimSpace(parts[1])
					// Remove any trailing text after the path
					if spaceIndex := strings.Index(sdkPath, " "); spaceIndex > 0 {
						sdkPath = sdkPath[:spaceIndex]
					}
					if f.IsValidAndroidSDK(sdkPath) {
						return sdkPath
					}
				}
			}
		}
	}
	
	return ""
}

// extractAndroidSDKPath extracts Android SDK path from flutter doctor message
func (f *FlutterSDKFinder) extractAndroidSDKPath(message string) string {
	// Common patterns in flutter doctor messages
	patterns := []string{
		"Android SDK at ",
		"located at ",
		"path: ",
	}
	
	messageLower := strings.ToLower(message)
	for _, pattern := range patterns {
		if idx := strings.Index(messageLower, pattern); idx != -1 {
			pathStart := idx + len(pattern)
			pathPart := message[pathStart:]
			
			// Find the end of the path (usually a space, comma, or newline)
			endChars := []string{" ", ",", "\n", "\r", "("}
			endIdx := len(pathPart)
			for _, endChar := range endChars {
				if idx := strings.Index(pathPart, endChar); idx != -1 && idx < endIdx {
					endIdx = idx
				}
			}
			
			sdkPath := strings.TrimSpace(pathPart[:endIdx])
			if f.IsValidAndroidSDK(sdkPath) {
				return sdkPath
			}
		}
	}
	
	return ""
}

// IsValidAndroidSDK checks if the path contains a valid Android SDK
func (f *FlutterSDKFinder) IsValidAndroidSDK(path string) bool {
	if path == "" {
		return false
	}
	
	// Check if the path exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	
	// Check for key Android SDK directories
	requiredDirs := []string{
		"platforms",
		"platform-tools",
		"build-tools",
	}
	
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(path, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return false
		}
	}
	
	return true
}