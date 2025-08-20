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

	apis, err := portaudio.HostApis()
	if err != nil {
		panic(errors.New("Could not query host APIs"))
	}

	LogAPIs(apis)

}

func LogAPIs(apis []*portaudio.HostApiInfo) {
	fmt.Println("APIs -------------------------------------")
	for i, a := range apis {
		fmt.Printf("API %d: %s\n", i, a.Name)
		fmt.Printf("Type: %s\n", a.Type)
		fmt.Println()
		LogIO(a.Devices)
	}
	fmt.Println()
}

func QueryIO() ([]*portaudio.DeviceInfo, error) {
	devices, err := portaudio.Devices()

	if err != nil {
		return []*portaudio.DeviceInfo{}, errors.New("Error enumerating audio devices")
	}

	return devices, nil
}

func LogIO(devices []*portaudio.DeviceInfo) {
	fmt.Println("---------------------------------- Devices")
	for _, d := range devices {
		var isInput bool
		if d.MaxInputChannels > 0 {
			isInput = true
		} else {
			isInput = false
		}

		if isInput {
			fmt.Printf("Input Device %d: %s\n", d.Index, d.Name)
			fmt.Printf("Max Channels: %d\n", d.MaxInputChannels)
		} else {
			fmt.Printf("Output Device %d: %s\n", d.Index, d.Name)
			fmt.Printf("Max Channels: %d\n", d.MaxOutputChannels)
		}
		fmt.Println("------------------------------------------")
	}

}
