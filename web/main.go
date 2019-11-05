package main

import (
	_ "github.com/lib/pq"
	"web/api"
	"web/types"
)

//database user has many message and activity exercice de maison
// Run api on a port here
func main() {
	api.Run(1234, &types.User{}, &types.Message{}, &types.Activity{})
}
