package main

import (
	"context"
	"fmt"

	"github.com/RohanKhatua/go-chi/application"
)

func main() {
	app := application.New()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
}
