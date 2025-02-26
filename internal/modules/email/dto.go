package email

type CreateEmailDTO struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	SendCode bool   `json:"send_code"`
}

type ValidateEmailDTO struct {
	EmailID string `json:"email_id"`
	Code    string `json:"code"`
}

type GenerateEmailValidationDTO struct {
	EmailID string `json:"email_id"`
	UserID  string `json:"user_id"`
}
