package qvo

import "time"

//WebpayReq struct holds a webpay transaction request's params.
type WebpayReq struct {
	Amount      int64  `json:"amount"`
	ReturnURL   string `json:"return_url"`
	CustomerID  string `json:"customer_id"`
	Description string `json:"description"`
}

//WebpayResp struct holds the answer from qvo for a webpay transaction request.
type WebpayResp struct {
	TransactionID  string    `json:"transaction_id"`
	RedirectURL    string    `json:"redirect_url"`
	ExpirationDate time.Time `json:"expiration_date"`
}

//WebpayTransaction begins a webpay transaction.
func WebpayTransaction(c *Client, req WebpayReq) error {

	return nil
}
