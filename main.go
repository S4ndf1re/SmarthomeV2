package main

import (
	"time"
)

func main() {

	container := NewScriptContainer("scripts")
	container.RunAll().Wait()

	for {
		time.Sleep(1000 * time.Millisecond)
	}
}
