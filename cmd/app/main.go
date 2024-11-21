package main

import (
	"jedyEvgeny/online-music-library/internal/pkg/app"
	"log"
)

func main() {
	a := app.New()
	err := a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
