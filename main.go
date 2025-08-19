package main

import (
	"errors"
	"fmt"

	"github.com/gordonklaus/portaudio"
)

func main() {

	if portaudio.Initialize() != nil {
		panic(errors.New("Portaudio couldn't initialize"))
	}

	defer portaudio.Terminate()
	devices, err := portaudio.Devices()

	if err != nil {
		panic(errors.New("Error enumerating audio devices"))
	}

	for _, d := range devices {
		fmt.Println(d.Name)
	}
}
