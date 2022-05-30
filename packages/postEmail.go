package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bugsnag/bugsnag-go/v2"
	mailgun "github.com/mailgun/mailgun-go/v4"
)

type Request struct {
	Email string `json:"email"`
}

type Response struct {
	StatusCode int               `json:"statusCode,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
	Body       string            `json:"body,omitempty"`
}

func Main(r Request) (*Response, error) {
	bugsnag.Configure(bugsnag.Configuration{
		// Your Bugsnag project API key, required unless set as environment
		// variable $BUGSNAG_API_KEY
		APIKey: "65ece1598afa0c809103b5c260ecdc6a",
		// The development stage of your application build, like "alpha" or
		// "production"
		ReleaseStage: "production",
		// The import paths for the Go packages containing your source files
		ProjectPackages: []string{"main", "github.com/org/myapp"},
		// more configuration options
	})

	if r.Email == "" {
		return nil, errors.New("email must be passed")
	}
	bugsnag.Notify(fmt.Errorf("Test error"))

	domain := os.Getenv("DOMAIN")
	apiKey := os.Getenv("API_KEY")
	id, err := SendMessage(domain, apiKey)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Response{
		Body: fmt.Sprintf("Email sent for id: %s", id),
	}, nil
}

func SendMessage(domain, apiKey string) (string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	m := mg.NewMessage(
		"Excited User <mailgun@YOUR_DOMAIN_NAME>",
		"Subject line",
		"Testing the body!",
		"test1234566789@mailinator.com",
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	// Send the message with a 10 second timeout
	resp, id, err := mg.Send(ctx, m)
	log.Println(resp)
	return id, err
}
