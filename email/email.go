package email

import (
	"encoding/base64"
	"fmt"

	"github.com/juju/errors"
	"google.golang.org/api/gmail/v1"
)

func GmailSend(service *gmail.Service, from, to, subject, body string) error {
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
