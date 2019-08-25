package models

import "time"

type (
	User struct {
		ID			int
		Name		string
		Password	string
		Email		string
		CreatedAt	time.Time
		LastLogin	time.Time
	}
	UserFilter struct {
		Search	string
	}
)

func (u *User) LoginTrack() error {
	u.LastLogin = time.Now()
	return db.UpdateField(u,"LastLogin", u.LastLogin)
}
