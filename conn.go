package qvo

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var sandboxURI = "https://playground.qvo.cl"
var productionURI = "https://api.qvo.cl"

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

func NewConnection(token string, isSandbox bool) *Conn {
	return &Conn{
		Token:     token,
		IsSandbox: isSandbox,
	}
}

func (conn *Conn) getURI() string {
	if conn.IsSandbox {
		return sandboxURI
	}
	return productionURI
}

func (conn *Conn) getBearer() string {
	return fmt.Sprintf("Bearer: %s", conn.Token)
}

func (conn *Conn) post(endpoint string, values url.Values) (string, QVOError) {
	uri := fmt.Sprintf("%s/%s", conn.getURI(), endpoint)
	req, err := http.NewRequest("POST", uri, strings.NewReader(values.Encode()))
	req.Header.Add("Authorization", conn.getBearer())
}
