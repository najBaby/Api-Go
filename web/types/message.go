package types

import (
	"web/orm"
)

type Message struct {
	Id      int    `json:",omitempty"`
	Title   string `json:",omitempty"`
	Content string `json:",omitempty"`
	User    *User  `json:",omitempty" orm:"rel(fk)"`
}

func init() {
	orm.RegisterModel(&Message{})
}
