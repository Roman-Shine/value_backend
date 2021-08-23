package service

type EmailService struct {
	//sender emailProvider.Sender
	//config config.EmailConfig
}

// Structures used for templates.
type verificationEmailInput struct {
	VerificationLink string
}

func (s *EmailService) SendUserVerificationEmail(input verificationEmailInput) error {
	// todo implement
	return nil
}
