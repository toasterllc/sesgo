package main

import (
	"fmt"
	"os"

	"github.com/heytoaster/sesgo"
)

func main() {
	if len(os.Args) != 5 {
		fmt.Fprintf(os.Stderr, "Usage:\n  email <to> <from> <subject> <body>\n\n")
		os.Exit(1)
	}

	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	if awsAccessKey == "" || awsSecretKey == "" {
		fmt.Fprintf(os.Stderr, "The AWS_ACCESS_KEY_ID / AWS_SECRET_ACCESS_KEY environment variables aren't set.\n")
		os.Exit(1)
	}

	err := sesgo.SendEmail(awsAccessKey, awsSecretKey, os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	if err != nil {
		fmt.Fprintf(os.Stderr, "SendEmail failed: %v\n\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sent\n")
}
