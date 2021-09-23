package main

import (
	"flag"
)

/**
 *  parse command-line arguments
 */
func parseCommandlineArgs() Args {
	url := flag.String("url", "https://www.wordpress-site.com", "URL to WordPress website to target (without trailing /)")
	pwdList := flag.String("dict", "dict.txt", "A text dictionary containing passwords to try")
	username := flag.String("user", "admin", "Username to brute-force")
	maxWorkers := flag.Int("routines", 10, "Maximum number of Goroutines to spawn")

	flag.Parse()

	return Args{
		"url":        *url,
		"dict":       *pwdList,
		"username":   *username,
		"maxWorkers": *maxWorkers,
	}
}
