package main

import (
	"log"
	"wall-server/app"
)


func main() {
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
