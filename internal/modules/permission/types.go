package permission

import "time"

type Entity struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Key         string    `json:"key" db:"key"`
	Created     time.Time `json:"created" db:"created"`
	Updated     time.Time `json:"updated" db:"updated"`
}
