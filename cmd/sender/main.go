package main

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

func main() {
	apiKey := "re_Bcc5tfNo_JfhsF7wJBR8dHfR6tsKrpbHj"
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Acme <onboarding@resend.dev>",
		To:      []string{"rafaelferreirab2@gmail.com"},
		Html:    "<strong>hello world</strong>",
		Subject: "Hello from Golang",
	}

	sent, err := client.Emails.Send(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(sent.Id)
}
