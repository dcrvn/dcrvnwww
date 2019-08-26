package models

import (
	"github.com/asdine/storm/q"
	"time"
)

type (
	User struct {
		ID			int
		Name		string
		DisplayName	string
		Password	string
		Email		string
		CreatedAt	time.Time
		LastLogin	time.Time
	}
	UserFilter struct {
		Pagination
		Search	string
	}
)

func (u *User) LoginTrack() error {
	u.LastLogin = time.Now()
	return db.UpdateField(u,"LastLogin", u.LastLogin)
}

func (u *User) Create() error {
	u.CreatedAt = time.Now()
	return db.Save(u)
}

func (u *User) Save() error {
	return db.Save(u)
}

func (u *User) Update() error {
	return db.Update(u)
}

func (u *User) GetByID() error {
	return  db.One("ID", u.ID, u)
}

func (f *UserFilter) GetList(count bool) (users []User, err error) {
	matchers := f.generateMatcher()

	query := db.Select(matchers...)
	if count {
		limit, skip := f.Paginate()
		query = query.Limit(limit).Skip(skip)
		var user User
		f.Total, _ = db.Select(matchers...).Count(&user)
	}
	if f.SortKey != "" {
		query = query.OrderBy(f.SortKey)
		if f.SortVal < 0 {
			query = query.Reverse()
		}
	} else {
		query = query.OrderBy("LastLogin").Reverse()
	}

	err = query.Each(new(User), func(record interface{}) error {
		if user, ok := record.(*User); ok {
			users = append(users, *user)
		}
		return nil
	})
	return
}


func (f UserFilter) generateMatcher() (matchers []q.Matcher) {
	if f.Search != "" {
		search := "(?i)"+f.Search
		matchers = append(matchers, q.Or(
			q.Re("Name", search), q.Re("DisplayName", search), q.Re("Email", search),
		))
	}
	return
}