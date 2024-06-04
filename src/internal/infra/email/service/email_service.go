package email_service

import (
	"auth-api/src/internal/domain/email"
	"auth-api/src/pkg/logger"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type EmailService struct {
	sesClient *ses.Client
	logger    logger.Logger
}

func NewEmailService(sesClient *ses.Client, logger logger.Logger) email.EmailService {
	return &EmailService{
		sesClient: sesClient,
		logger:    logger,
	}
}

func (h *EmailService) SendEmail(input email.Email) error {
	if err := input.Validate(); err != nil {
		return err
	}

	sesSendEmailInput := &ses.SendEmailInput{
		Source: aws.String("motacartmarpaulo@gmail.com"), //TODO: Get this from config
		Destination: &types.Destination{
			ToAddresses: []string{
				input.To,
			},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(input.Body),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(input.Subject),
			},
		},
	}

	_, err := h.sesClient.SendEmail(context.TODO(), sesSendEmailInput)
	if err != nil {
		h.logger.Error("failed to send email: %v", err)
		return err
	}

	h.logger.Info("email sent successfully")
	return nil
}

// func (h *EmailService) Handle(event events.Event) error {
// 	emailEvent, ok := event.(events.EmailEvent)
// 	if !ok {
// 		h.logger.Error("invalid event type")
// 		return errors.New("invalid event type")
// 	}

// 	input := &ses.SendEmailInput{
// 		Source: aws.String("motacartmarpaulo@gmail.com"), //TODO: Get this from config
// 		Destination: &types.Destination{
// 			ToAddresses: []string{
// 				emailEvent.To,
// 			},
// 		},
// 		Message: &types.Message{
// 			Body: &types.Body{
// 				Text: &types.Content{
// 					Charset: aws.String("UTF-8"),
// 					Data:    aws.String(emailEvent.Body),
// 				},
// 			},
// 			Subject: &types.Content{
// 				Charset: aws.String("UTF-8"),
// 				Data:    aws.String(emailEvent.Subject),
// 			},
// 		},
// 	}

// 	_, err := h.sesClient.SendEmail(context.TODO(), input)
// 	if err != nil {
// 		h.logger.Error("failed to send email: %v", err)
// 		return err
// 	}

// 	h.logger.Info("email sent successfully")
// 	return nil
// }
