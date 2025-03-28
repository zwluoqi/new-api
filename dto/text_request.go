package dto

import "encoding/json"

type GeneralOpenAIRequest struct {
	Model               string         `json:"model,omitempty"`
	ReasoningEffort     string         `json:"reasoning_effort,omitempty"`
	Messages            []Message      `json:"messages,omitempty"`
	Prompt              any            `json:"prompt,omitempty"`
	BestOf              int            `json:"best_of,omitempty"`
	Echo                bool           `json:"echo,omitempty"`
	Stream              bool           `json:"stream,omitempty"`
	StreamOptions       *StreamOptions `json:"stream_options,omitempty"`
	Suffix              string         `json:"suffix,omitempty"`
	MaxTokens           uint           `json:"max_tokens,omitempty"`
	MaxCompletionTokens uint           `json:"max_completion_tokens,omitempty"`
	Temperature         float64        `json:"temperature,omitempty"`
	TopP                float64        `json:"top_p,omitempty"`
	TopK                int            `json:"top_k,omitempty"`
	Stop                any            `json:"stop,omitempty"`
	N                   int            `json:"n,omitempty"`
	Input               any            `json:"input,omitempty"`
	Instruction         string         `json:"instruction,omitempty"`
	Size                string         `json:"size,omitempty"`
	Functions           any            `json:"functions,omitempty"`
	FrequencyPenalty    float64        `json:"frequency_penalty,omitempty"`
	PresencePenalty     float64        `json:"presence_penalty,omitempty"`
	ResponseFormat      any            `json:"response_format,omitempty"`
	Seed                float64        `json:"seed,omitempty"`
	Tools               []ToolCall     `json:"tools,omitempty"`
	ToolChoice          any            `json:"tool_choice,omitempty"`
	User                string         `json:"user,omitempty"`
	LogitBias           any            `json:"logit_bias,omitempty"`
	LogProbs            any            `json:"logprobs,omitempty"`
	TopLogProbs         int            `json:"top_logprobs,omitempty"`
	Dimensions          int            `json:"dimensions,omitempty"`
	ParallelToolCalls   bool           `json:"parallel_tool_calls,omitempty"`
	EncodingFormat      any            `json:"encoding_format,omitempty"`

	Thinking *Thinking `json:"thinking,omitempty"`
}

type Thinking struct {
	Type         string `json:"type"`
	BudgetTokens int    `json:"budget_tokens"`
}

type OpenAITools struct {
	Type     string         `json:"type"`
	Function OpenAIFunction `json:"function"`
}

type OpenAIFunction struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Parameters  any    `json:"parameters,omitempty"`
}

type StreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"`
}

func (r GeneralOpenAIRequest) GetMaxTokens() int {
	return int(r.MaxTokens)
}

func (r GeneralOpenAIRequest) ParseInput() []string {
	if r.Input == nil {
		return nil
	}
	var input []string
	switch r.Input.(type) {
	case string:
		input = []string{r.Input.(string)}
	case []any:
		input = make([]string, 0, len(r.Input.([]any)))
		for _, item := range r.Input.([]any) {
			if str, ok := item.(string); ok {
				input = append(input, str)
			}
		}
	}
	return input
}

type Message struct {
	Role             string          `json:"role"`
	Content          json.RawMessage `json:"content"`
	ReasoningContent *string         `json:"reasoning_content,omitempty"`
	Name             *string         `json:"name,omitempty"`
	ToolCalls        any             `json:"tool_calls,omitempty"`
	ToolCallId       string          `json:"tool_call_id,omitempty"`
}

type MediaMessage struct {
	Type     string `json:"type"`
	Text     string `json:"text"`
	ImageUrl any    `json:"image_url,omitempty"`
}

type MessageImageUrl struct {
	Url    string `json:"url"`
	Detail string `json:"detail"`
}

const (
	ContentTypeText     = "text"
	ContentTypeImageURL = "image_url"
)

func (m Message) StringContent() string {
	var stringContent string
	if err := json.Unmarshal(m.Content, &stringContent); err == nil {
		return stringContent
	}
	return string(m.Content)
}

func (m *Message) SetStringContent(content string) {
	jsonContent, _ := json.Marshal(content)
	m.Content = jsonContent
}

func (m Message) IsStringContent() bool {
	var stringContent string
	if err := json.Unmarshal(m.Content, &stringContent); err == nil {
		return true
	}
	return false
}

func (m Message) ParseContent() []MediaMessage {
	var contentList []MediaMessage
	var stringContent string
	if err := json.Unmarshal(m.Content, &stringContent); err == nil {
		contentList = append(contentList, MediaMessage{
			Type: ContentTypeText,
			Text: stringContent,
		})
		return contentList
	}
	var arrayContent []json.RawMessage
	if err := json.Unmarshal(m.Content, &arrayContent); err == nil {
		for _, contentItem := range arrayContent {
			var contentMap map[string]any
			if err := json.Unmarshal(contentItem, &contentMap); err != nil {
				continue
			}
			switch contentMap["type"] {
			case ContentTypeText:
				if subStr, ok := contentMap["text"].(string); ok {
					contentList = append(contentList, MediaMessage{
						Type: ContentTypeText,
						Text: subStr,
					})
				}
			case ContentTypeImageURL:
				if subObj, ok := contentMap["image_url"].(map[string]any); ok {
					detail, ok := subObj["detail"]
					if ok {
						subObj["detail"] = detail.(string)
					} else {
						subObj["detail"] = "high"
					}
					contentList = append(contentList, MediaMessage{
						Type: ContentTypeImageURL,
						ImageUrl: MessageImageUrl{
							Url:    subObj["url"].(string),
							Detail: subObj["detail"].(string),
						},
					})
				} else if url, ok := contentMap["image_url"].(string); ok {
					contentList = append(contentList, MediaMessage{
						Type: ContentTypeImageURL,
						ImageUrl: MessageImageUrl{
							Url:    url,
							Detail: "high",
						},
					})
				}

			}
		}
		return contentList
	}

	return nil
}
