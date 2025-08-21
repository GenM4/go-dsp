package main

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/GenM4/go-dsp/internal/gen"
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

	devices, err := QueryIO()
	if err != nil {
		panic(err)
	}

	p := buildStreamParams(nil, devices[2], 2)

	sine, err := gen.NewStereoSine(1000, 1000, p)
	if err != nil {
		panic(err)
	}
	defer sine.Close()

	err = sine.Start()
	fmt.Println("sine started")
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	err = sine.Stop()
	fmt.Println("sine stopped")
	if err != nil {
		panic(err)
	}
}

func buildStreamParams(in, out *portaudio.DeviceInfo, channels int) *portaudio.StreamParameters {
	var sp *portaudio.StreamParameters
	var sr float64
	if in != nil {
		sr = in.DefaultSampleRate
	} else {
		sr = out.DefaultSampleRate
	}

	sp = &portaudio.StreamParameters{
		Input:           *buildStreamDeviceParams(in, channels),
		Output:          *buildStreamDeviceParams(out, channels),
		SampleRate:      sr,
		FramesPerBuffer: 512,
		Flags:           portaudio.NoFlag,
	}

	return sp
}

func buildStreamDeviceParams(device *portaudio.DeviceInfo, channels int) *portaudio.StreamDeviceParameters {
	var sdp *portaudio.StreamDeviceParameters

	var lat time.Duration
	if device != nil && device.MaxInputChannels > 0 {
		lat = device.DefaultHighInputLatency
	} else if device != nil && device.MaxInputChannels == 0 {
		lat = device.DefaultHighOutputLatency
	} else {
		lat = 0
	}

	sdp = &portaudio.StreamDeviceParameters{
		Device:   device,
		Channels: channels,
		Latency:  lat,
	}

	return sdp
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
	fmt.Println("Devices ----------------------------------")
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
		fmt.Printf("Sample Rate: %0.0f\n", d.DefaultSampleRate)
		fmt.Println("------------------------------------------")
	}

}
