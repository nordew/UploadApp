package main

import (
	"github.com/nordew/UploadApp/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
