package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jjzcru/djinn/pkg/controller"
)

func main() {
	fmt.Println("Loading plugins:")
	ctrl := &controller.Controller{}
	err := ctrl.OnLoadPlugins()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	
	plugins := ctrl.Plugins

	for k := range plugins {
		fmt.Printf("Plugin '%s' loaded successfully\n", k)
	}
	
	ctrl.OnRun()
}
