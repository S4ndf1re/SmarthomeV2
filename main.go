package main

import (
	"Smarthome/gui"
	"Smarthome/scripting"
	"log"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	guiHandler := gui.NewGui()
	guiHandler.AddContainer(passwordChangeContainer())
	container := scripting.NewScriptContainer("scripts", guiHandler)
	container.RunAll().Wait()

	if err := http.ListenAndServe(":8080", guiHandler); err != nil {
		panic(err)
	}
}
