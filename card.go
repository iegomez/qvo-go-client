package qvo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

//Status constants
const (
	Succeeded string = "succeeded"
	Failed    string = "failed"
)

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
func CreateCardInscription(c *Client, customerID, returnURL string) (CardInscriptionResponse, error) {

	endpoint := fmt.Sprintf("customers/%s/cards/inscriptions", customerID)

	form := url.Values{}
	form.Add("customer_id", customerID)
	form.Add("return_url", returnURL)

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

//ChargeCard creates a charge for given customer and card.
func ChargeCard(c *Client, customerID, cardID, description string, amount int64) (Transaction, error) {
	endpoint := fmt.Sprintf("customers/%s/cards/%s/charge", customerID, cardID)

	form := url.Values{}
	form.Add("customer_id", customerID)
	form.Add("card_id", cardID)
	form.Add("amount", strconv.FormatInt(amount, 10))
	form.Add("description", description)

	body, err := c.request("POST", endpoint, form)
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

//DeleteCard deletes a card for a given customer.
func DeleteCard(c *Client, customerID, cardID string) error {

	endpoint := fmt.Sprintf("customers/%s/cards/%s", customerID, cardID)

	form := url.Values{}
	form.Add("customer_id", customerID)
	form.Add("card_id", cardID)

	_, err := c.request("DELETE", endpoint, form)
	if err != nil {
		return err
	}

	return nil

}

//ListCards retrieves cards for a given customer.
func ListCards(c *Client, customerID string) ([]Card, error) {

	var cards = make([]Card, 0)

	form := url.Values{}
	form.Add("customer_id", customerID)

	body, err := c.request("GET", "customers", form)
	//log.Debugf("\n\nbody: %s\n\n", body)
	if err != nil {
		log.Errorf("errored at body: %s", err)
		return cards, err
	}

	err = json.NewDecoder(bytes.NewReader(body)).Decode(&cards)

	if err != nil {
		log.Errorf("errored at unmarshal: %s", err)
		return cards, err
	}

	return cards, nil

}
