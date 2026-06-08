package error_responses

import (

	// commnuity pacakges

	"fmt"
)

type ErrorResponse struct {
	MessageID string
	Err       error
}

// Error implements error.
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("MessageID: %s, Error: %v", e.MessageID, e.Err)
}

func (e *ErrorResponse) ErrorString() string {
	return fmt.Sprintf("MessageID: %s, Error:%v", e.MessageID, e.Err)
}

func (e *ErrorResponse) NewErrorResponse(messageID string, err error) *ErrorResponse {
	return &ErrorResponse{
		MessageID: messageID,
		Err:       err,
	}
}
