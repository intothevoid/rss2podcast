package audio

import (
	"log"
	"os/exec"
)

func ConvertWavToMp3(wavFile string, mp3File string) error {
	cmd := exec.Command("ffmpeg", "-i", wavFile, mp3File)

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Println("Conversion complete!")
	return nil
}
