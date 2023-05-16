package chatresponse

type Message struct {
	ID         string   `json:"id"`
	Author     Author   `json:"author"`
	CreateTime float64  `json:"create_time"`
	UpdateTime float64  `json:"update_time"`
	Content    Content  `json:"content"`
	EndTurn    *bool    `json:"end_turn,omitempty"`
	Weight     float64  `json:"weight"`
	Metadata   Metadata `json:"metadata"`
	Recipient  string   `json:"recipient"`
}

type Author struct {
	Role     string            `json:"role"`
	Name     *string           `json:"name,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

type Content struct {
	ContentType string   `json:"content_type"`
	Parts       []string `json:"parts"`
}

type Metadata struct {
	MessageType   string        `json:"message_type"`
	ModelSlug     string        `json:"model_slug"`
	FinishDetails FinishDetails `json:"finish_details"`
}

type FinishDetails struct {
	Type string `json:"type"`
	Stop string `json:"stop"`
}

type Conversation struct {
	Message        Message `json:"message"`
	ConversationID string  `json:"conversation_id"`
	Error          *string `json:"error,omitempty"`
}
