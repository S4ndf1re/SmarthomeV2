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
	container := scripting.NewScriptContainer("scripts", guiHandler)
	container.RunAll().Wait()

	//tempCont := gui.NewContainer("test", func(name string) { log.Printf("%s\n", name) })

	//tempBtn := gui.NewButton("button1", "Click Me", func(name string) { log.Printf("button1 clicked by %s\n", name) })
	//tempCont.Add(tempBtn)

	//tempCheck := gui.NewCheckbox("check1", "Check Me")
	//tempCheck.SetChangeCallback(func(username string, state bool) {
	//	log.Printf("check1 clicked by %s. Status: %t\n", username, state)
	//})
	//tempCont.Add(tempCheck)

	//guiHandler.AddContainer(tempCont)

	if err := http.ListenAndServe(":8080", guiHandler); err != nil {
		panic(err)
	}
}
