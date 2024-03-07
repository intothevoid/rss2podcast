package tts

import (
	"io/ioutil"
	"os/exec"
)

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertToAudio(text, filename string) error {
	err := ioutil.WriteFile("temp.txt", []byte(text), 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("say", "-v", "Alex", "-o", filename, "file", "temp.txt")
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
