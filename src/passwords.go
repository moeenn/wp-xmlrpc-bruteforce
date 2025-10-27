package main

import (
	"bufio"
	"fmt"
	"os"
)

// get from file one password at a time.
func getPasswords(pwdListFile string) (chan string, error) {
	pwdChan := make(chan string)

	file, err := os.Open(pwdListFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open password lists file")
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pwdChan <- scanner.Text()
	}

	return pwdChan, nil
}
