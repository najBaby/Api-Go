package types

import (
	"web/orm"
)

type Activity struct {
	Id    int
	Title string
	Todo  string
}

func init() {
	orm.RegisterModel(&Activity{})
}
