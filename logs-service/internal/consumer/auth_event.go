package consumer

type AuthEvent struct {
	EventType string `json:"event_type"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
}
