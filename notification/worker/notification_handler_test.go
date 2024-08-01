package worker_test

import (
	"encoding/json"
	"testing"

	"github.com/mailersend/mailersend-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sejamuchhal/taskhub/notification/mocks/mock_email_sender"
	"github.com/sejamuchhal/taskhub/notification/mocks/mock_rabbitmq"
	event "github.com/sejamuchhal/taskhub/notification/pb"
	"github.com/sejamuchhal/taskhub/notification/worker"
	"go.uber.org/mock/gomock"
)

func TestNotificationHandler_EmailSentSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailSender := mock_email_sender.NewMockEmailSenderInterface(ctrl)
	mockRabbitMQ := mock_rabbitmq.NewMockRabbitMQBrokerInterface(ctrl)
	w := worker.Worker{
		EmailSender:    mockEmailSender,
		RabbitMQBroker: mockRabbitMQ,
	}

	taskUpdateEvent := event.TaskUpdateEvent{
		Email: "test@example.com",
		Title: "Test Task",
	}
	msgBody, _ := json.Marshal(taskUpdateEvent)
	msg := amqp.Delivery{Body: msgBody}

	mockEmailSender.EXPECT().SendEmail(taskUpdateEvent.Email, gomock.Any(), gomock.Any()).Return(&mailersend.Response{}, nil).Times(1)

	w.NotificationHandler("test-queue", msg, nil)

}

func TestNotificationHandler_UnmarshallingError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEmailSender := mock_email_sender.NewMockEmailSenderInterface(ctrl)
	mockRabbitMQ := mock_rabbitmq.NewMockRabbitMQBrokerInterface(ctrl)
	w := worker.Worker{
		EmailSender:    mockEmailSender,
		RabbitMQBroker: mockRabbitMQ,
	}

	invalidJson := []byte(`{"invalid_json"}`)

	msg := amqp.Delivery{Body: invalidJson}
	w.NotificationHandler("test-queue", msg, nil)
}
