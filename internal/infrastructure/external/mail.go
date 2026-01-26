package external

import "doan/internal/external"

type mailService struct {
}

func NewMailService() external.MailService {
	return &mailService{}
}

func (m *mailService) SendEmail(to string, subject string, body string) error {
	return nil
}
