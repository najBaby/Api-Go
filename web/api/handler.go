package api

import (
	"net/http"
	"reflect"
)

type Handler struct {
	Entities map[string]reflect.Type
}

func (h *Handler) entities(entities []interface{}) {
	mp := make(map[string]reflect.Type)
	for _, entity := range entities {
		etype := reflect.TypeOf(entity)
		if etype.Kind() == reflect.Ptr {
			etype = etype.Elem()
		}
		if etype.Kind() != reflect.Struct {
			panic("the entities must be of type Struct")
		}
		mp[etype.Name()] = etype
	}
	h.Entities = mp
}

func (h *Handler) handle(req *Request, res *Response) {
	switch {
	case req.Mutation != nil:
		res.Operation = "mutation"
		req.Mutation.ParseMutation(h.Entities)
		if errors := req.Mutation.Errors; errors != nil {
			res.Errors = errors
		} else {
			res.Data = req.Mutation.Filter
		}
	case req.Query != nil:
		res.Operation = "query"
		req.Query.ParseQuery(h.Entities)
		if errors := req.Query.Errors; errors != nil {
			res.Errors = errors
		} else {
			res.Data = req.Query.Serve()
		}
	default:
	}
	res.ServeJSON()
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.handle(NewRequest(r), NewResponse(w))
}
