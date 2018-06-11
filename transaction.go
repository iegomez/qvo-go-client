package qvo

import "time"

//GatewayResponse struct to deal with gateway response from transactions.
type GatewayResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

//Transaction struct to represent a qvo transaction object.
type Transaction struct {
	ID              string                 `json:"id"`
	Amount          int64                  `json:"amount"`
	Currency        string                 `json:"currency"` //CLP or USD.
	Description     string                 `json:"description"`
	Gateway         string                 `json:"gateway"` //One of: webpay_plus, webpay_oneclick, olpays.
	Credits         int64                  `json:"credits"`
	Status          string                 `json:"ststaus"` //One of: successful, rejected, unable_to_charge, refunded, waiting_for_response, response_timeout.
	Customer        Customer               `json:"customer"`
	Payment         *Payment               `json:"payment"` //Nullable, so it's a pointer.
	Refund          *Refund                `json:"refund"`  //Nullable, so it's a pointer.
	Transable       map[string]interface{} //API sends a "hash", so we are limited to an interfaces map.
	GatewayResponse GatewayResponse        `json:"gateway_response"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}
