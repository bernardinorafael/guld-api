package team

// Expects a team entity as parameter
var InsertTeamQuery = `
	INSERT INTO teams (
		id,
		name,
		slug,
		owner_id,
		org_id,
		logo,
		created,
		members_count,
		updated
	)
	VALUES (
		:id,
		:name,
		:slug,
		:owner_id,
		:org_id,
		:logo,
		:created,
		:members_count,
		:updated
	)
`

// Expects owner_id and org_id as parameters
var FindAllTeamsQuery = `
	SELECT * FROM teams WHERE owner_id = $1 AND org_id = $2 ORDER BY created DESC
`
