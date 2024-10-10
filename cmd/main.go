package main

import (
	_ "practice/docs/swagger"
	"practice/internal/app"
	"practice/internal/pkg"
)

// @title Practice
// @version 1.0
// @description Practice API
// @host localhost:8080
// @BasePath  /
func main() {
	app.New(pkg.Module).Run()
}
