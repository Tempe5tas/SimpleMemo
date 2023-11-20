package main

import (
	"SimpleMemo/conf"
	"SimpleMemo/model"
	"SimpleMemo/router"
)

func main() {
	conf.Init()
	model.Init()
	r := router.Init()

	r.Run(conf.ServiceAddr)
}
