package main

import (
	"gitlab.glaske.net/mglaske/observerip"
)

func main() {
	x, err := observerip.New(5080)
	if err != nil {
		panic(err)
	}
	x.ListenAndServe()
}
