package audio

import (
	"log"
	"os/exec"
)

// ConvertWavToMp3 converts a WAV audio file to an MP3 audio file.
// It uses the ffmpeg command-line tool to perform the conversion.
// The input WAV file path is specified by the `wavFile` parameter,
// and the output MP3 file path is specified by the `mp3File` parameter.
// If the conversion is successful, it returns `nil`. Otherwise, it
// returns an error indicating the cause of the failure.
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
