package main

type Args map[string]interface{}

type Credentials struct {
	Username string
	Password string
}

type Result struct {
	Credentials
	Success bool
}
