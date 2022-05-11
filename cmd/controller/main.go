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

	deviceId := "d073d568f92a"
	plugins["com.lifx.bulbs"].Command(deviceId, "SET_COLOR", `{"color": "#00FF00"}`)
	// plugins["com.lifx.bulbs"].Command(deviceId, "SET_BRIGHTNESS", `{"brightness": 1.0}`)
	
	ctrl.OnRun()
}
