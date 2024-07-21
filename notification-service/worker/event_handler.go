package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (w *Worker) ReminderHandler(queue string, msg amqp.Delivery, err error) {
	logger := log.WithFields(log.Fields{"method": "ReminderHandler"})

	if err != nil {
		logger.WithError(err).Error("Error occurred in RMQ consumer")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	logger.Infof("Message received on '%s' queue: %s", queue, string(msg.Body))
	var reminder storage.Reminder
	err = json.Unmarshal(msg.Body, &reminder)
	if err != nil {
		logger.WithError(err).Error("Error while unmarshalling reminder")
		return
	}

	_, err = w.EmailSender.SendEmail(reminder.Email, "Reminder", reminder.Message)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"reminder_id": reminder.Id,
			"email":       reminder.Email,
		}).WithError(err).Error("Failed to send email reminder")
	}

	logger.WithFields(logrus.Fields{"reminder_id": reminder.Id, "email": reminder.Email}).Info("Email sent successfully")

}