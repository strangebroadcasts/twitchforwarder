package discord

// WebhookExecution holds the parameters for an execution of a Discord webhook
type WebhookExecution struct {
	Content  string `json:"content"`
	Username string `json:"username,omitempty"`
}
