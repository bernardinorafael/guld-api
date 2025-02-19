package email

type CreateParams struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}
