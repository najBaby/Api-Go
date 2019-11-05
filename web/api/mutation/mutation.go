package mutation

import (
	"fmt"
	"reflect"
	"web/bd"
)

type Mutation struct {
	Entity string      `json:"entity"`
	Fields interface{} `json:"fields"`
	Filter interface{} `json:"filter"`
	Errors []string    `json:"-"`
	Method string      `json:"-"`
}

func (mutation *Mutation) ParseMutation(entitiesmap map[string]reflect.Type) {
	if etype, ok := entitiesmap[mutation.Entity]; ok {
		ParseFields(mutation, etype)
		ParseFilter(mutation, etype)
	} else {
		mutation.Errors = append(mutation.Errors, fmt.Sprintf("unknow entity '%v'", mutation.Entity))
	}
}

func (mutation *Mutation) POST() {
	bd.Create(mutation.Fields, mutation.Filter)
}

func (mutation *Mutation) PUT() {

}

func (mutation *Mutation) PATCH() {

}

func (mutation *Mutation) DELETE() {

}

func (mutation *Mutation) Serve() {
	switch mutation.Method {
	case "POST":
		mutation.POST()
	case "PUT":
		mutation.PUT()
	case "PATCH":
		mutation.PATCH()
	case "DELETE":
		mutation.DELETE()
	default:
	}
}
