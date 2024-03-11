package podcast

import (
	"fmt"

	"github.com/intothevoid/rss2podcast/pkg/llm"
)

type Podcast interface {
	GenerateSummary(title string, description string, content string) (string, error)
	GenerateIntroduction(subject string, podcaster string) (string, error)
}

type podcast struct {
	ollama llm.LLM
}

func NewPodcast(ol llm.LLM) Podcast {
	return &podcast{
		ollama: ol,
	}
}

func (pod *podcast) GenerateSummary(title string, description string, content string) (string, error) {
	retval, err := pod.ollama.SendRequest(fmt.Sprintf("Process this article as a transcript to be read on a news podcast. Only talk about the title, description and the content. Do not add any introduction. Title: %s. Description: %s, Content: %s.", title, description, content))
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}

	return retval, nil
}

func (pod *podcast) GenerateIntroduction(subject string, podcaster string) (string, error) {
	retval, err := pod.ollama.SendRequest(fmt.Sprintf("Generate an introduction for a podcast. Subject: %s. Podcaster: %s.", subject, podcaster))
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}

	return retval, nil
}
