package apputils

import (
	"strings"

	"github.com/mileusna/useragent"
)

// Returns a concise summary like "BrowserName vX.Y on OS X.Y"
func SummarizeUserAgent(uaString string) string {
	ua := useragent.Parse(uaString)
	name := ua.Name
	version := ua.Version
	osName := ua.OS
	osVersion := ua.OSVersion

	// Ambil major.minor version browser
	majorMinor := ""
	if version != "" {
		parts := strings.Split(version, ".")
		if len(parts) >= 2 {
			majorMinor = parts[0] + "." + parts[1]
		} else {
			majorMinor = parts[0]
		}
	}

	// Ringkas nama OS populer dan ambil versi utama OS jika ada
	// macOS: "Mac OS X" → "macOS"
	if strings.Contains(osName, "Mac OS X") || strings.Contains(osName, "Intel Mac OS X") {
		osName = "macOS"
	}
	// Windows: "Windows NT 10.0" → "Windows 10"
	if strings.Contains(osName, "Windows") && osVersion == "10.0" {
		osName = "Windows"
		osVersion = "10"
	}
	if strings.Contains(osName, "Windows") && osVersion == "11.0" {
		osName = "Windows"
		osVersion = "11"
	}
	// iOS: "iPhone OS" → "iOS"
	if strings.Contains(osName, "iPhone OS") {
		osName = "iOS"
	}
	// Android: "Android"
	if strings.Contains(osName, "Android") {
		osName = "Android"
	}

	if name == "" && osName == "" {
		return "Unknown"
	}
	result := name
	if majorMinor != "" {
		result += " v" + majorMinor
	}
	if osName != "" {
		result += " on " + osName
		if osVersion != "" && osVersion != "0" {
			// Only append version if not empty or "0"
			result += " " + osVersion
		}
	}
	return strings.TrimSpace(result)
}
