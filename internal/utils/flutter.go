package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"tilokit/pkg/constants"
)

type FlutterSDKFinder struct{}

func (f *FlutterSDKFinder) FindFlutterSDK() (string, error) {
	if sdkPath := f.getFlutterSDKFromDoctor(); sdkPath != "" {
		if f.isValidFlutterSDK(sdkPath) {
			return sdkPath, nil
		}
	}

	if flutterRoot := os.Getenv("FLUTTER_ROOT"); flutterRoot != "" {
		if f.isValidFlutterSDK(flutterRoot) {
			return flutterRoot, nil
		}
	}

	if flutterPath := f.findFlutterInPath(); flutterPath != "" {
		sdkPath := f.extractSDKFromExecutable(flutterPath)
		if f.isValidFlutterSDK(sdkPath) {
			return sdkPath, nil
		}
	}

	commonPaths := f.getCommonFlutterPaths()
	for _, path := range commonPaths {
		if f.isValidFlutterSDK(path) {
			return path, nil
		}
	}

	return "", fmt.Errorf("flutter SDK not found. Please ensure Flutter is installed and either:\n1. Add Flutter to your PATH\n2. Set FLUTTER_ROOT environment variable\n3. Install Flutter in a standard location")
}

func (f *FlutterSDKFinder) findFlutterInPath() string {
	return ""
}

func (f *FlutterSDKFinder) extractSDKFromExecutable(execPath string) string {
	binDir := filepath.Dir(execPath)
	sdkPath := filepath.Dir(binDir)
	return sdkPath
}

func (f *FlutterSDKFinder) getCommonFlutterPaths() []string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	paths := []string{}

	if runtime.GOOS == constants.OSDarwin {
		paths = append(paths,
			filepath.Join(homeDir, "flutter"),
			filepath.Join(homeDir, "development", "flutter"),
			"/usr/local/flutter",
			"/opt/flutter",
		)
	}

	if runtime.GOOS == constants.OSLinux {
		paths = append(paths,
			filepath.Join(homeDir, "flutter"),
			filepath.Join(homeDir, "development", "flutter"),
			"/usr/local/flutter",
			"/opt/flutter",
		)
	}

	if runtime.GOOS == constants.OSWindows {
		paths = append(paths,
			filepath.Join(homeDir, "flutter"),
			filepath.Join(homeDir, "development", "flutter"),
			"C:\\flutter",
			"C:\\src\\flutter",
		)
	}

	return paths
}

func (f *FlutterSDKFinder) getFlutterSDKFromDoctor() string {
	if sdkPath := f.tryFlutterDoctorMachineForSDK(); sdkPath != "" {
		return sdkPath
	}
	return f.tryFlutterDoctorVerboseForSDK()
}

func (f *FlutterSDKFinder) tryFlutterDoctorMachineForSDK() string {
	return ""
}

func (f *FlutterSDKFinder) tryFlutterDoctorVerboseForSDK() string {
	return ""
}

func (f *FlutterSDKFinder) isValidFlutterSDK(path string) bool {
	if path == "" {
		return false
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

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

func (f *FlutterSDKFinder) GetFlutterInfo(sdkPath string) map[string]string {
	info := make(map[string]string)

	if sdkPath == "" {
		return info
	}

	if version, err := f.getFlutterVersion(sdkPath); err == nil {
		info["version"] = version
	}

	if channel, err := f.getFlutterChannel(sdkPath); err == nil {
		info["channel"] = channel
	}

	if dartVersion, err := f.getDartVersion(sdkPath); err == nil {
		info["dart_version"] = dartVersion
	}

	info["sdk_path"] = sdkPath
	return info
}

func (f *FlutterSDKFinder) getFlutterVersion(_ string) (string, error) {
	return "3.16.0", nil
}

func (f *FlutterSDKFinder) getFlutterChannel(_ string) (string, error) {
	return "stable", nil
}

func (f *FlutterSDKFinder) getDartVersion(_ string) (string, error) {
	return "3.2.0", nil
}

func (f *FlutterSDKFinder) PromptUserToContinue() bool {
	return true
}
