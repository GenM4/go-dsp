package main

import (
	"errors"
	"fmt"
	"math"
	"time"

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

	p := buildStreamParams(nil, devices[1], 2)

	sine, err := newStereoSine(256, 320, p)
	if err != nil {
		panic(err)
	}
	defer sine.Close()

	err = sine.Start()
	if err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second)

	err = sine.Stop()
	if err != nil {
		panic(err)
	}
}

type stereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
}

func (g *stereoSine) processAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * g.phaseL))
		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
		out[1][i] = float32(math.Sin(2 * math.Pi * g.phaseR))
		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
	}
}

func newStereoSine(freqL, freqR float64, p *portaudio.StreamParameters) (*stereoSine, error) {
	s := &stereoSine{nil, freqL / p.SampleRate, 0, freqR / p.SampleRate, 0}

	var err error
	s.Stream, err = portaudio.OpenStream(*p, s.processAudio)
	if err != nil {
		return nil, errors.New("Could not open stream")
	}

	return s, nil
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
		Input:      *buildStreamDeviceParams(in, channels),
		Output:     *buildStreamDeviceParams(out, channels),
		SampleRate: sr,
		Flags:      portaudio.NoFlag,
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
