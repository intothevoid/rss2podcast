package tts

import (
	"fmt"
)

type ConverterConfig struct {
	URL    string
	Voice  string
	Speed  float64
	Format string
}

type Converter struct {
	engine string
	coqui  *CoquiConverter
	kokoro *KokoroConverter
}

func NewConverter(engine string, coquiConfig, kokoroConfig *ConverterConfig) *Converter {
	return &Converter{
		engine: engine,
		coqui:  NewCoquiConverter(coquiConfig),
		kokoro: NewKokoroConverter(kokoroConfig),
	}
}

// ConvertToAudio converts text to audio using the configured TTS engine
func (c *Converter) ConvertToAudio(content string, fileName string) error {
	switch c.engine {
	case "coqui":
		return c.coqui.ConvertToAudio(content, fileName)
	case "kokoro":
		return c.kokoro.ConvertToAudio(content, fileName)
	default:
		return fmt.Errorf("unsupported TTS engine: %s", c.engine)
	}
}
