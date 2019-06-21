package bot

import (
	"encoding/json"
	"regexp"
	"fmt"
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

// SearchMessages searches through the group messages for a specified string (does not take capitalization into consideration)
func (res *Response) SearchMessages(str string) []Message {
	var matchedMessages []Message
	reg := fmt.Sprintf("(?i)%s", str)

	for _, v := range res.Data.Messages {
		match, err := regexp.MatchString(reg, v.Message)
		if err != nil {
			Handle(err)
		} else if match {
			matchedMessages = append(matchedMessages, v)
		}
	}

	return matchedMessages
}