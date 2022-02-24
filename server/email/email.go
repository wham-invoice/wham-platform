package email

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"github.com/juju/errors"
	"github.com/rstorr/wham-platform/db"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func SendInvoice(ctx context.Context, user *db.User, invoice *db.Invoice) error {
	b, err := ioutil.ReadFile("/Users/work/go/src/github.com/rstorr/wham-platform/secrets/google_web_client_credentials.json")
	if err != nil {
		return errors.Trace(err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailComposeScope)
	if err != nil {
		return errors.Trace(err)
	}

	httpClient := config.Client(context.Background(), &user.OAuth)
	service, err := gmail.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return errors.Trace(err)
	}

	invoiceURL := fmt.Sprintf("https://whaminvoice.co.nz/invoice/%s", invoice.ID)
	body := fmt.Sprintf("Hi %s,\n\n"+
		"Your invoice is ready.\n\n"+
		"To view and download it please visit: %s "+
		"Thanks.\n"+
		"%s", invoice.Client.FirstName, invoiceURL, user.FirstName)

	return email(service, "me", invoice.Client.Email, "Invoice", body)
}

func email(service *gmail.Service, from, to, subject, body string) error {

	message := &gmail.Message{}
	messageStr := []byte(fmt.Sprintf(
		"From: %s \r\n"+
			"To: %s \r\n"+
			"Subject: %s\r\n\r\n"+
			"%s", from, to, subject, body))
	message.Raw = base64.URLEncoding.EncodeToString(messageStr)

	_, err := service.Users.Messages.Send(from, message).Do()

	return errors.Trace(err)
}
