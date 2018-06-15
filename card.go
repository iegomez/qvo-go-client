package qvo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"
)

//Status constants
const (
	Succeeded string = "succeeded"
	Failed    string = "failed"
)

//CardInscriptionRequest holds a customer id and the return url.
type CardInscriptionRequest struct {
	CustomerID string `json:"customer_id"`
	ReturnURL  string `json:"return_url"`
}

//CardInscriptionResponse struct holds the answer from qvo for a card inscription response.
type CardInscriptionResponse struct {
	InscriptionUID string    `json:"inscription_uid"`
	RedirectURL    string    `json:"redirect_url"`
	ExpirationDate time.Time `json:"expiration_date"`
}

//CardInscriptionState holds the state of a card inscription request.
type CardInscriptionState struct {
	UID       string        `json:"uid"`
	Status    string        `json:"status"`
	Card      *Card         `json:"card"`
	Error     *errorWrapper `json:"error"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

//Card struct to represent a qvo card object.
type Card struct {
	ID           string    `json:"id"`
	Lats4Digits  string    `json:"last_4_digits"`
	CardType     string    `json:"card_type"`    //VISA or MASTERCARD
	PaymentType  string    `json:"payment_type"` //CD (credit) or DB (debit)
	FailureCount int32     `json:"failure_count"`
	CreatedAt    time.Time `json:"created_at"`
}

//CreateCardInscription begins a card inscription request. If everything's ok, it'll return an inscription uid, the redirect url to send the customer to, and the expiration date for this transaction.
func CreateCardInscription(c *Client, customerID string, req CardInscriptionRequest) (CardInscriptionResponse, error) {

	endpoint := fmt.Sprintf("customers/%s/cards/inscriptions", customerID)

	form := url.Values{}
	form.Add("customer_id", req.CustomerID)
	form.Add("return_url", req.ReturnURL)

	body, err := c.request("POST", endpoint, form)
	if err != nil {
		return CardInscriptionResponse{}, err
	}

	var resp CardInscriptionResponse
	err = json.Unmarshal(body, &resp)

	if err != nil {
		return CardInscriptionResponse{}, err
	}

	return resp, nil

}

//GetCardInscription returns the inscription's state and a card (if successful).
func GetCardInscription(c *Client, customerID, inscriptionUID string) (CardInscriptionState, error) {
	endpoint := fmt.Sprintf("customers/%s/cards/inscriptions/%s", customerID, inscriptionUID)

	form := url.Values{}
	form.Add("customer_id", customerID)
	form.Add("inscription_uid", inscriptionUID)

	body, err := c.request("GET", endpoint, form)
	if err != nil {
		return CardInscriptionState{}, err
	}

	var cardInscriptionState CardInscriptionState
	err = json.Unmarshal(body, &cardInscriptionState)

	if err != nil {
		return CardInscriptionState{}, err
	}

	return cardInscriptionState, nil
}

//GetCard returns a card given a customer id and a card id.
func GetCard(c *Client, customerID, cardID string) (Card, error) {
	endpoint := fmt.Sprintf("customers/%s/cards/%s", customerID, cardID)

	form := url.Values{}
	form.Add("customer_id", customerID)
	form.Add("card_id", cardID)

	body, err := c.request("GET", endpoint, form)
	if err != nil {
		return Card{}, err
	}

	var card Card
	err = json.Unmarshal(body, &card)

	if err != nil {
		return Card{}, err
	}

	return card, nil

}
