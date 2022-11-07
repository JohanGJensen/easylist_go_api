package structs

type Message struct {
	Message string `json:"message" binding:"required"`
	Token   string `json:"token,omitempty"`
}
