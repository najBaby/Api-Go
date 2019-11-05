package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"web/api/mutation"
	"web/api/query"
)

type Request struct {
	Query    *query.Query
	Mutation *mutation.Mutation
}

func NewRequest(httprequest *http.Request) *Request {
	req := Request{}
	switch httprequest.Header.Get("Content-Type") {
	case "application/x-www-form-urlencoded":
		query := query.Query{
			Method: httprequest.Method,
			Entity: httprequest.FormValue("entity"),
			Limits: 1,
			Fields: nil,
		}
		if limits := httprequest.FormValue("limits"); limits != "" {
			u, err := strconv.ParseInt(limits, 10, 64)
			if err != nil {
				query.Errors = append(query.Errors, err.Error())
			} else {
				query.Limits = u
			}
		}

		if fields := httprequest.FormValue("fields"); fields != "" {
			query.Fields = strings.Split(strings.ReplaceAll(fields, " ", ""), ",")
		}
		if err := httprequest.ParseForm(); err != nil {
			query.Errors = []string{err.Error()}
		}
		for _, v := range []string{"entity", "fields", "orders", "limits"} {
			delete(httprequest.Form, v)
		}
		query.Filter = httprequest.Form
		req.Query = &query
	case "application/json":
		mutation := mutation.Mutation{
			Method: httprequest.Method,
		}
		decoder := json.NewDecoder(httprequest.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&mutation); err != nil {
			mutation.Errors = append(mutation.Errors, err.Error())
		}
		req.Mutation = &mutation
	}
	return &req
}
