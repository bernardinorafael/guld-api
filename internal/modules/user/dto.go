package user

type UserRegisterParams struct {
	FullName     string `json:"full_name"`
	Username     string `json:"username"`
	EmailAddress string `json:"email_address"`
	PhoneNumber  string `json:"phone_number"`
}

type CreatePhoneData struct {
	Phone     string `json:"phone"`
	IsPrimary bool   `json:"is_primary"`
}

type CreateEmailParams struct {
	Email     string `json:"email"`
	UserID    string `json:"user_id"`
	IsPrimary bool   `json:"is_primary"`
}

type UserProfileUpdate struct {
	FullName string `json:"full_name" db:"full_name"`
	Username string `json:"username" db:"username"`
}

type UserSearchParams struct {
	Sort  string `json:"sort"`
	Query string `json:"q"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}
