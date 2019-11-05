package bd

import (
	"fmt"
	"reflect"
	"web/orm"
)

func Read(entity reflect.Type, fields []string, filter map[string]interface{}, limits int64) interface{} {
	o := orm.NewOrm()
	c := orm.NewCondition()
	result := reflect.New(reflect.SliceOf(entity)).Interface()
	for k, v := range filter {
		c.And(k, v)
		fmt.Println(c)
	}
	fmt.Println(fields)
	o.QueryTable(entity.Name()).Limit(limits).RelatedSel().SetCond(c).All(result, fields)
	return result
}
