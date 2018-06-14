package qvo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

//GatewayResponse struct to deal with gateway response from transactions.
type GatewayResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

//Transaction struct to represent a qvo transaction object.
type Transaction struct {
	ID              string                  `json:"id"`
	Amount          int64                   `json:"amount"`
	Currency        string                  `json:"currency"` //CLP or USD.
	Description     string                  `json:"description"`
	Gateway         string                  `json:"gateway"` //One of: webpay_plus, webpay_oneclick, olpays.
	Credits         int64                   `json:"credits"`
	Status          string                  `json:"ststaus"` //One of: successful, rejected, unable_to_charge, refunded, waiting_for_response, response_timeout.
	Customer        Customer                `json:"customer"`
	Payment         *Payment                `json:"payment"` //Nullable, so it's a pointer.
	Refund          *Refund                 `json:"refund"`  //Nullable, so it's a pointer.
	Transable       *map[string]interface{} //API sends a "hash", so we are limited to an interfaces map. Also, it's nullable, so it's a pointer.
	GatewayResponse GatewayResponse         `json:"gateway_response"`
	CreatedAt       time.Time               `json:"created_at"`
	UpdatedAt       time.Time               `json:"updated_at"`
}

//GetTransaction retrieves a transaction by id.
func GetTransaction(c *Client, id string) (Transaction, error) {

	endpoint := fmt.Sprintf("transactions/%s", id)

	form := url.Values{}
	form.Add("transaction_id", id)

	body, err := c.request("GET", endpoint, form)
	if err != nil {
		return Transaction{}, err
	}

	var transaction Transaction
	err = json.Unmarshal(body, &transaction)

	if err != nil {
		return Transaction{}, err
	}

	return transaction, nil

}

//RefundTransaction makes a refund request for a given transaction id.
func RefundTransaction(c *Client, id string) (Refund, error) {
	endpoint := fmt.Sprintf("transactions/%s/refund", id)

	form := url.Values{}
	form.Add("transaction_id", id)

	body, err := c.request("POST", endpoint, form)
	if err != nil {
		return Refund{}, err
	}

	var refund Refund
	err = json.Unmarshal(body, &refund)

	if err != nil {
		return Refund{}, err
	}

	return refund, nil
}
