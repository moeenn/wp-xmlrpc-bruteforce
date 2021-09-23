package main

import (
	"bufio"
	"log"
	"os"
)

/**
 *  get from file one password at a time
 */
func getPasswords(pwdList string, passwords chan<- string) {
	file, err := os.Open(pwdList)
	if err != nil {
		log.Fatalln("Failed to open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// send one password at a time, into the passwords channel
	for scanner.Scan() {
		passwords <- scanner.Text()
	}

	close(passwords)
}
