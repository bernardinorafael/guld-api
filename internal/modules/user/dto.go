package user

type UserRegisterDTO struct {
	FullName     string `json:"full_name"`
	Username     string `json:"username"`
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number"`
}

type UpdateProfileDTO struct {
	FullName             string `json:"full_name" db:"full_name"`
	Username             string `json:"username" db:"username"`
	IgnorePasswordPolicy bool   `json:"ignore_password_policy" db:"ignore_password_policy"`
}
