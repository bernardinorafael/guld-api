package user

var InsertEmailQuery = `
	INSERT INTO emails (id, user_id, email, is_primary, is_verified) VALUES ($1, $2, $3, true, true)
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
