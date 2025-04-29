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
	const podcastSummaryPrompt = `Your task is to transform the following title:%s, description:%s, content:%s into a concise summary that a human podcast host can read verbatim.
		IMPORTANT: Generate ONLY the plain summary text with NO additional elements whatsoever:

		NO introductions (like "Here's a summary...")
		NO headings (like "Summary:")
		NO sound effect instructions (like "Sound of musical transition")
		NO meta-commentary about the article
		NO phrases like "the article states/details/mentions"
		NO formatting markers or special characters
		NO conclusions or sign-offs

		The summary should:

		Be 200-350 words in length
		Use natural, conversational language appropriate for verbal delivery
		Present only factual information from the original article
		Be written in the first person as if the podcaster is directly sharing the news
		Flow as a coherent, stand-alone piece that needs no introduction or conclusion

		Begin the summary immediately without any preamble and end it without any closing remarks.`

	retval, err := pod.ollama.SendRequest(fmt.Sprintf(podcastSummaryPrompt, title, description, content))

	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}

	return retval, nil
}

func (pod *podcast) GenerateIntroduction(subject string, podcaster string) (string, error) {
	retval, err := pod.ollama.SendRequest(fmt.Sprintf("Generate a very short introduction for a podcast. Subject: %s. Podcaster: %s.", subject, podcaster))
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}

	return retval, nil
}
