package gen

import (
	"errors"
	"github.com/gordonklaus/portaudio"
	"math"
)

type StereoSine struct {
	*portaudio.Stream
	stepL, phaseL float64
	stepR, phaseR float64
}

func (g *StereoSine) ProcessAudio(out [][]float32) {
	for i := range out[0] {
		out[0][i] = float32(math.Sin(2 * math.Pi * g.phaseL))
		_, g.phaseL = math.Modf(g.phaseL + g.stepL)
		out[1][i] = float32(math.Sin(2 * math.Pi * g.phaseR))
		_, g.phaseR = math.Modf(g.phaseR + g.stepR)
	}
}

func NewStereoSine(freqL, freqR float64, p *portaudio.StreamParameters) (*StereoSine, error) {
	s := &StereoSine{nil, freqL / p.SampleRate, 0, freqR / p.SampleRate, 0}

	var err error
	s.Stream, err = portaudio.OpenStream(*p, s.ProcessAudio)
	if err != nil {
		return nil, errors.New("Could not open stream")
	}

	return s, nil
}
