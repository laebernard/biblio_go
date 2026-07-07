package main

import "os"

var API_URL string

func init() {
	API_URL = os.Getenv("API_URL")
	if API_URL == "" {
		API_URL = "http://localhost:8080"
	}
}
