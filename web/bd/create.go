package bd

import (
	"fmt"
	"web/orm"
)

func Create(fields, filter interface{}) {
	o := orm.NewOrm()
	id, err := o.Insert(filter)
	if err != nil {
		panic(err)
	}
	fmt.Println(id)
}
