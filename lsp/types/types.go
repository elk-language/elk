package types

type Type interface {
	lspType()
}

type Message struct {
	JSONRPC string `json:"jsonrpc"`
}

func (Message) lspType() {}

type MessageWithMethod struct {
	Message

	/**
	 * The method to be invoked.
	 */
	Method string `json:"method"`
}

type RequestMessage struct {
	MessageWithMethod

	/**
	 * The request id.
	 */
	ID string `json:"id"`

	/**
	 * The method to be invoked.
	 */
	Method string `json:"method"`
}

type ResponseMessage struct {
	Message

	/**
	 * The request id.
	 */
	ID *string `json:"id,omitempty"`
}

type NotificationMessage struct {
	MessageWithMethod
}
