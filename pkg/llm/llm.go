package llm

type LLM interface {
	SendRequest(prompt string) (string, error)
	CheckConnection() error
}
