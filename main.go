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

	devices, err := QueryIO()
	if err != nil {
		panic(err)
	}

	LogIO(devices)

	fmt.Println()
	fmt.Println()

	apis, err := portaudio.HostApis()
	fmt.Println("------------------------------------------")
	for _, a := range apis {
		fmt.Println(a.Name)
		fmt.Println(a.Type)
		LogIO(a.Devices)
		fmt.Println("------------------------------------------")
	}

}

func QueryIO() ([]*portaudio.DeviceInfo, error) {
	devices, err := portaudio.Devices()

	if err != nil {
		return []*portaudio.DeviceInfo{}, errors.New("Error enumerating audio devices")
	}

	return devices, nil
}

func LogIO(devices []*portaudio.DeviceInfo) {
	fmt.Println("------------------------------------------")
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
