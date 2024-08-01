package worker

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
	event "github.com/sejamuchhal/taskhub/notification/pb"
	log "github.com/sirupsen/logrus"
)


func (w *Worker) NotificationHandler(queue string, msg amqp.Delivery, err error) {
	logger := log.WithFields(log.Fields{"method": "NotificationHandler"})

	if err != nil {
		logger.WithError(err).Error("Error occurred in RMQ consumer")
	}

	logger.Infof("Message received on '%s' queue: %s", queue, string(msg.Body))
	var message event.TaskUpdateEvent
	err = json.Unmarshal(msg.Body, &message)
	if err != nil {
		logger.WithError(err).Error("Error while unmarshalling reminder")
		return
	}

	mailSubject, mailBody := prepareEmailContent(message)

	_, err = w.EmailSender.SendEmail(message.Email, mailSubject, mailBody)
	if err != nil {
		logger.WithFields(log.Fields{
			"task_title": message.Title,
			"email":      message.Email,
		}).WithError(err).Error("Failed to send email reminder")
	}

	logger.WithFields(log.Fields{"reminder_id": message.Title, "email": message.Email}).Info("Email sent successfully")
}
