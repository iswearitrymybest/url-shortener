package main

import (
	"fmt"

	"github.com/ghostvoid/url-shortener/internal/config"
)

func main() {
	//cleanevn
	cfg := config.MustLoad()

	fmt.Println(cfg)

	//TODO: init logger

	//TODO: init logger

	//TODO: init router

	//TODO: run server

}
