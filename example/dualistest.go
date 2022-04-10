package main

import (
	"github.com/darmiel/dualis-push/dualis"
	"log"
	"os"
	"time"
)

func main() {
	user, pass := os.Getenv("USER"), os.Getenv("PASS")

	log.Println("Logging in...")
	client, err := dualis.Login(user, pass)
	if err != nil {
		panic(err)
	}
	log.Println("Logged in. Refresh:", client.ArgumentsFromRefresh())

	swStart := time.Now()
	var exams []*dualis.Grade
	if exams, err = client.CourseResults(); err != nil {
		panic(err)
	}
	swStop := time.Now()

	log.Println("Got", len(exams), "Exams! Took", swStop.Sub(swStart))
}
