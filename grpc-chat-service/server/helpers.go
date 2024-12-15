package server

import "strings"

// Helper untuk memeriksa apakah recipient adalah grup
func isGroup(recipient string) bool {
	// Aturan sederhana untuk menentukan apakah recipient adalah grup:
	// Misalnya, grup memiliki tanda khusus seperti '@'
	return strings.Contains(recipient, "@")
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
