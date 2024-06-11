package main

import (
	"github.com/islombay/cms/api"
	"github.com/islombay/cms/db"
)

func main() {
	db := db.Init()

	r := api.Init(db)

	r.Run("0.0.0.0:10000")
}
