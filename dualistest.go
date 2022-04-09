package main

import (
	"github.com/darmiel/dualis-push/dualis"
	"log"
	"os"
)

func main() {
	user, pass := os.Getenv("USER"), os.Getenv("PASS")

	log.Println("Logging in...")
	client, err := dualis.Login(user, pass)
	if err != nil {
		panic(err)
	}
	log.Println("Logged in. Refresh:", client.ArgumentsFromRefresh())

	if err = client.CourseResults(); err != nil {
		panic(err)
	}
}
