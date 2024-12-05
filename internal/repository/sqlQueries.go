package repository

const (
	querySetUserInfo = `
		INSERT INTO users_info (id, ip, hashed_refresh_token)
		VALUES ($1, $2, $3)
	`
)