package bot

import (
	"encoding/json"
)

// Response contains Messages
type Response struct {
	Data Messages `json:"response"`
}

// Messages a slice of type Message
type Messages struct {
	Messages []Message `json:"messages"`
}

// Message contains an ID, Name of who sent the message, and the text it contains
type Message struct {
	MsgID string `json:"id"`
	From string `json:"name"`
	Message string `json:"text"`
}

// ParseMessages Unmarshals a slice of bytes into a Response of Messages
func ParseMessages(data []byte) *Response {
	var res Response
	json.Unmarshal(data, &res)
	
	return &res
}