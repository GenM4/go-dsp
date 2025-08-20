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
		if d.MaxInputChannels != 0 {
			fmt.Printf("Input Device %d: %s\n", d.Index, d.Name)
			fmt.Printf("Max Channels: %d\n", d.MaxInputChannels)
		} else {
			fmt.Printf("Output Device %d: %s\n", d.Index, d.Name)
			fmt.Printf("Max Channels: %d\n", d.MaxOutputChannels)
		}
	}
}
