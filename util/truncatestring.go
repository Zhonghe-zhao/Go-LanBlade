package util

// TruncateString truncates a string to the specified length and adds "..." if truncated.
func TruncateString(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}
