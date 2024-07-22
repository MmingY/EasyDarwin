package utils

import (
	"os"
	// "syscall"
)

// RedirectStderr to the file passed in
func RedirectStderr() (err error) {
	logFile, err := os.OpenFile(ErrorLogFilename(), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer logFile.Close() // Ensure the file is closed when function returns

	// Redirect stderr to the log file
	// Save the original stderr
	// oldStderr := os.Stderr
	
	// Set os.Stderr to the log file
	os.Stderr = logFile
	return
}
