package qvo

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

//WebpayRequest struct holds a webpay transaction request's params.
type WebpayRequest struct {
	Amount      int64  `json:"amount"`
	ReturnURL   string `json:"return_url"`
	CustomerID  string `json:"customer_id"`
	Description string `json:"description"`
}

//WebpayResponse struct holds the answer from qvo for a webpay transaction request.
type WebpayResponse struct {
	TransactionID  string    `json:"transaction_id"`
	RedirectURL    string    `json:"redirect_url"`
	ExpirationDate time.Time `json:"expiration_date"`
}

//WebpayTransaction begins a webpay transaction. If everything's ok, it'll return a transaction id (for later check), the redirect url to send the customer to, and the expiration date for this transaction.
func WebpayTransaction(c *Client, req WebpayRequest) (WebpayResponse, error) {

	form := url.Values{}
	form.Add("amount", strconv.Itoa(int(req.Amount)))
	form.Add("customer_id", req.CustomerID)
	form.Add("return_url", req.ReturnURL)
	form.Add("Description", req.Description)

	body, err := c.request("POST", "webpay_plus/charge", form)
	if err != nil {
		return WebpayResponse{}, err
	}

	var resp WebpayResponse
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return WebpayResponse{}, err
	}

	return resp, nil

}
