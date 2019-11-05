package api

import (
	"fmt"
	"web/orm"
	"log"
	"net/http"
)

type api struct {
	srv *http.Server
}

func Connectdb() {
	err := orm.RegisterDriver("postgres", orm.DRPostgres)
	if err != nil {
		panic(err)
	}
	err = orm.RegisterDataBase("default", "postgres", "postgres://postgres:@172.17.0.1:5432/postgres?sslmode=disable")
	if err != nil {
		panic(err)
	}
	orm.RunSyncdb("default", false, true)
}

func Run(port uint32, entities ...interface{}) {
	Connectdb()
	handler := &Handler{}
	handler.entities(entities)
	s := &api{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: handler,
		},
	}
	log.Printf("Server listens on port %d", port)
	log.Fatal(s.srv.ListenAndServe())
}
