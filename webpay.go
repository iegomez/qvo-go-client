package qvo

import (
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

//WebpayResponse struct holds the answer from qvo for a webpay transaction request.
type WebpayResponse struct {
	TransactionID  string    `json:"transaction_id"`
	RedirectURL    string    `json:"redirect_url"`
	ExpirationDate time.Time `json:"expiration_date"`
}

//WebpayTransaction begins a webpay transaction. If everything's ok, it'll return a transaction id (for later check), the redirect url to send the customer to, and the expiration date for this transaction.
func WebpayTransaction(c *Client, customerID, returnURL, description string, amount int64) (WebpayResponse, error) {

	form := url.Values{}
	form.Add("amount", strconv.FormatInt(amount, 10))
	form.Add("customer_id", customerID)
	form.Add("return_url", returnURL)
	form.Add("Description", description)

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
