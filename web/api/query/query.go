package query

import (
	"fmt"
	"reflect"
	"web/bd"
)

type Query struct {
	Method string
	Limits int64
	Entity interface{}
	Errors []string
	Fields []string
	Filter interface{}
}

func (query *Query) ParseQuery(entitiesmap map[string]reflect.Type) {
	if etype, ok := entitiesmap[query.Entity.(string)]; ok {
		query.Entity = etype
		ParseFields(query, etype)
		ParseFilter(query, etype)
	} else {
		query.Errors = append(query.Errors, fmt.Sprintf("unknow entity '%v'", query.Entity))
	}
}

func (query *Query) Serve() interface{} {
	if query.Method == "GET" || query.Method == "HEAD" {
		return bd.Read(query.Entity.(reflect.Type), query.Fields, query.Filter.(map[string]interface{}), query.Limits)
	}
	return nil
}
