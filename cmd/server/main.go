package main

import "github.com/MuriloAbranches/Go-Expert-API/configs"

func main() {
	config, _ := configs.LoadConfig(".")
	println(config.DBDriver)
}
