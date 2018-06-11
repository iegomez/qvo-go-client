package qvo

var sandboxURI = "https://playground.qvo.cl"

type Conn struct {
	Token     string
	IsSandbox bool
}

type QVOError struct {
	Type    string  `json:"type"`
	Message *string `json:"message"`
	Param   *string `json:"param"`
}

type Filter struct {
	Attribute string
	Operator  string
	Value     interface{}
}
