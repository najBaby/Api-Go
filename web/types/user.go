package types

import (
	"web/orm"
)

type User struct {
	Id       int        `json:",omitempty"`
	Email    string     `json:",omitempty"`
	Password string     `json:",omitempty"`
	Messages []*Message `json:",omitempty" orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(&User{})
}

//docker run -d -p 5432:5432 -v postgres-data:/var/lib/postgresql/data `
//  --name postgres1 postgres
