package service

// WeeklyReportPayload is the email payload
type WeeklyReportPayload struct {
	Email  string
	Name   string
	Report string
}

// Notification defines how this app comunicates with user
type Notification interface {
	SendWeeklyReport(payload *WeeklyReportPayload) error
}
