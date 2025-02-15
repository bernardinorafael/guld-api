package user

// TODO: Add created and updated
var InsertPhoneQuery = `
	INSERT INTO phones (id, user_id, phone, is_primary) VALUES ($1, $2, $3, true)
`

// TODO: Add created and updated
var InsertEmailQuery = `
	INSERT INTO emails (id, user_id, email, is_primary) VALUES ($1, $2, $3, true)
`

var InsertUserQuery = `
	INSERT INTO users (
		id,
		full_name,
		username,
		phone_number,
		email_address,
		avatar_url,
		banned,
		locked,
		created,
		updated
	) VALUES (
		:id,
		:full_name,
		:username,
		:phone_number,
		:email_address,
		:avatar_url,
		:banned,
		:locked,
		:created,
		:updated
	)
`
