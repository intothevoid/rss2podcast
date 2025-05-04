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
	mlx    *MLXAudioConverter
}

func NewConverter(engine string, coquiConfig, kokoroConfig, mlxConfig *ConverterConfig) *Converter {
	return &Converter{
		engine: engine,
		coqui:  NewCoquiConverter(coquiConfig),
		kokoro: NewKokoroConverter(kokoroConfig),
		mlx:    NewMLXAudioConverter(mlxConfig),
	}
}

// ConvertToAudio converts text to audio using the configured TTS engine
func (c *Converter) ConvertToAudio(content string, fileName string) error {
	switch c.engine {
	case "coqui":
		return c.coqui.ConvertToAudio(content, fileName)
	case "kokoro":
		return c.kokoro.ConvertToAudio(content, fileName)
	case "mlx":
		return c.mlx.ConvertToAudio(content, fileName)
	default:
		return fmt.Errorf("unsupported TTS engine: %s", c.engine)
	}
}

// PlayAudio plays the audio file using the configured TTS engine
func (c *Converter) PlayAudio(fileName string) error {
	if c.engine != "mlx" {
		return fmt.Errorf("play audio is only supported for MLX Audio engine")
	}
	return c.mlx.PlayAudio(fileName)
}

// StopAudio stops any currently playing audio using the configured TTS engine
func (c *Converter) StopAudio() error {
	if c.engine != "mlx" {
		return fmt.Errorf("stop audio is only supported for MLX Audio engine")
	}
	return c.mlx.StopAudio()
}

// OpenOutputFolder opens the output folder in the system's file explorer
func (c *Converter) OpenOutputFolder() error {
	if c.engine != "mlx" {
		return fmt.Errorf("open output folder is only supported for MLX Audio engine")
	}
	return c.mlx.OpenOutputFolder()
}
