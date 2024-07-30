package worker

import (
	"fmt"

	event "github.com/sejamuchhal/taskhub/notification/pb"
)

func prepareEmailContent(message event.TaskUpdateEvent) (string, string) {
	mailSubject := fmt.Sprintf("%s: %s", message.Title, message.Status)
	mailBody := fmt.Sprintf("Task Title: %s\nStatus: %s", message.Title, message.Status)
	return mailSubject, mailBody
}
