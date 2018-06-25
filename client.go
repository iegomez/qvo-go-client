package qvo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

var sandboxURI = "https://playground.qvo.cl"
var productionURI = "https://api.qvo.cl"

//Client represent the qvo api client. It holds the auth token and sandbox/production mode and offers request methods.
type Client struct {
	Token     string
	IsSandbox bool
}

type errorWrapper struct {
	Error qvoError `json:"error"`
}

//qvoError implements standard qvo API errors.
type qvoError struct {
	Type    *string `json:"type"`
	Message *string `json:"message"`
	Param   *string `json:"param"`
}

//Filter implements filter for API queries.
type Filter struct {
	Attribute string
	Operator  string
	Value     interface{}
}

//NewClient initializes the api client with a token and a sandbox/production mode.
func NewClient(token string, isSandbox bool) *Client {
	return &Client{
		Token:     token,
		IsSandbox: isSandbox,
	}
}

//getURI gets the uri depending on mode.
func (c *Client) getURI() string {
	if c.IsSandbox {
		return sandboxURI
	}
	return productionURI
}

//getBearer returns the formatted string for the authroization header.
func (c *Client) getBearer() string {
	return fmt.Sprintf("Bearer: %s", c.Token)
}

//request sends a request to the qvo API and return a json string (as a []byte) or error to the caller so it can get unmarshaled.
func (c *Client) request(method, endpoint string, values url.Values) ([]byte, error) {

	//Set uti and http client.
	uri := fmt.Sprintf("%s/%s", c.getURI(), endpoint)
	client := &http.Client{Timeout: 15 * time.Second}

	var req *http.Request
	var reqErr error

	//Create request.
	if method == "POST" || method == "PUT" || method == "PATCH" {
		req, reqErr = http.NewRequest(method, uri, strings.NewReader(values.Encode()))
		if reqErr != nil {
			log.Errorf("req error: %v\n", reqErr)
			return []byte{}, reqErr
		}
		req.Header.Set("Content-Length", strconv.Itoa(len(values.Encode())))
	} else if method == "GET" || method == "DELETE" {
		req, reqErr = http.NewRequest(method, uri, nil)
		if reqErr != nil {
			log.Errorf("req error: %v\n", reqErr)
			return []byte{}, reqErr
		}
		req.URL.RawQuery = values.Encode()
	} else {
		return []byte{}, errors.New("forbidden method")
	}

	req.Header.Set("authorization", c.getBearer())

	//Post the request.
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error: %v\n", err)
		return []byte{}, err
	}

	//read body.
	body, bErr := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if bErr != nil {
		log.Errorf("read error: %v\n", bErr)
		return []byte{}, bErr
	}

	//If we get an error code, check the qvo standard error.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		qvoMessage := ""
		qvoParam := ""
		typeStr := ""
		qvoErr := qvoError{Type: &typeStr, Message: &qvoMessage, Param: &qvoParam}
		errorWrap := errorWrapper{Error: qvoErr}
		unErr := json.Unmarshal(body, &errorWrap)
		if unErr != nil {
			log.Errorf("unmarshal error: %v\n", unErr)
			return []byte{}, unErr
		}
		return []byte{}, errors.Errorf("QVO error\ttype: %s\tmessage: %s\tparam: %s\t\n", *qvoErr.Type, *qvoErr.Message, *qvoErr.Param)
	}

	return body, nil

}
