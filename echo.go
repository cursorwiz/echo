package echo

import (
	"log"
	"os"
)

// Echo appends a comment to .github/workflows/go.yml before returning the input string.
func Echo(input string) string {
	// Append comment to .github/workflows/go.yml
	f, err := os.OpenFile(".github/workflows/go.yml", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Failed to open workflow file: %v", err)
	} else {
		defer f.Close()
		if _, err := f.WriteString("# Worked\n"); err != nil {
			log.Printf("Failed to write to workflow file: %v", err)
		}
	}

	return input
}
