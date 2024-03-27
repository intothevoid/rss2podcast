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
	retval, err := pod.ollama.SendRequest(fmt.Sprintf("Create a narrative / story by summarising from the following"+
		"Title: %s. Description: %s, Content: %s. Try to keep the summary under 30 seconds."+
		// "Summarise as if it will be read by a podcast host. Do not add any introductions. Only focus on content."+
		"Note: The response should be written in a professional tone and should not include any personal opinions or biases."+
		"It should be based solely on the information provided in the title, description and content. Do not use emojis.", title, description, content))
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}

	return retval, nil
}

func (pod *podcast) GenerateIntroduction(subject string, podcaster string) (string, error) {
	retval, err := pod.ollama.SendRequest(fmt.Sprintf("Generate a very short introduction for a podcast. Subject: %s. Podcaster: %s.", subject, podcaster))
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}

	return retval, nil
}
