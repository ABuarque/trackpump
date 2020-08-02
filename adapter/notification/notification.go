package notification

import (
	"fmt"
	"trackpump/email"
	"trackpump/usecase/service"
)

type notificationService struct {
	Email        string
	Password     string
	emailService *email.Client
}

// NewNotificationService returns a new notification service
func NewNotificationService(e, p string) service.Notification {
	return &notificationService{
		Email:        e,
		Password:     p,
		emailService: email.New(e, p),
	}
}

func (n *notificationService) SendWeeklyReport(payload *service.WeeklyReportPayload) error {
	msg := "From: " + n.Email + "\n" +
		"To: " + payload.Email + "\n" +
		"Subject: Weekly workout report\n\n" +
		payload.Report
	if err := n.emailService.Send(n.Email, []string{payload.Email}, "Weekly workout report", msg); err != nil {
		return fmt.Errorf("failed to send email to %s, erro %q", payload.Email, err)
	}
	return nil
}
