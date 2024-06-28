package main

import (
	"fmt"
	"os"

	"github.com/toasterllc/sesgo"
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

	to := os.Args[1]
	from := os.Args[2]
	subject := os.Args[3]
	body := os.Args[4]
	err := sesgo.SendEmail(awsAccessKey, awsSecretKey, from, []string{to}, nil, nil, subject, body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "SendEmail failed: %v\n\n", err)
		os.Exit(1)
	}

	fmt.Printf("Sent\n")
}
