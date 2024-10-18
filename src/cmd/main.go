package main

import (
	"github.com/betterreads/internal/application"
)

func main() {
	r := application.NewRouter(":8080")
	r.Run()
}
