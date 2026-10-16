package internal

import (
	"bufio"
	"fmt"
	"os"
)

// get from file one password at a time.
func GetPasswords(pwdListFile string) (chan string, error) {
	pwdChan := make(chan string)

	file, err := os.Open(pwdListFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open password lists file")
	}

	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			pwdChan <- scanner.Text()
		}
		_ = file.Close()
		close(pwdChan)
	}()

	return pwdChan, nil
}
