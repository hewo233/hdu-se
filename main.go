package main

import (
	"github.com/hewo233/hdu-se/Init"
	"github.com/hewo233/hdu-se/route"
)

func main() {
	Init.AllInit()
	route.InitRoute()

	err := route.R.Run(":8080")
	if err != nil {
		panic(err)
	}
}
