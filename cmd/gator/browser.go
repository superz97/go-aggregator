package main

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func openInBrowser(url string) error {
	if runtime.GOOS == "linux" && isWSL() {
		return exec.Command("explorer.exe", url).Start()
	}

	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", url).Start()
	case "windows":
		return exec.Command("cmd", "/c", "start", url).Start()
	default:
		return exec.Command("xdg-open", url).Start()
	}
}

func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}
