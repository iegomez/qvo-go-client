package qvo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

//Withdrawal struct to represent a qvo withdraw object.
type Withdrawal struct {
	ID        string    `json:"id"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"` //One of: processing, rejected, transfered.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//CreateWithdrawal creates a withdrawal of the given amount. Return a Withdrawal object or an error.
func CreateWithdrawal(c *Client, amount int64) (Withdrawal, error) {

	var withdrawal Withdrawal

	//Validate required fields.
	if amount <= 0 {
		return Withdrawal{}, errors.New("can't create a withdrawal with negative or 0 amount")
	}

	form := url.Values{}
	form.Add("amount", strconv.FormatInt(amount, 10))

	body, err := c.request("POST", "withdrawals", form)
	if err != nil {
		return Withdrawal{}, err
	}

	err = json.Unmarshal(body, &withdrawal)

	if err != nil {
		return Withdrawal{}, err
	}

	return withdrawal, nil
}

//GetWithdrawal retrieves a withdrawal given its id.
func GetWithdrawal(c *Client, id string) (Withdrawal, error) {

	endpoint := fmt.Sprintf("withdrawals/%s", id)

	form := url.Values{}
	form.Add("withdrawal_id", id)

	body, err := c.request("GET", endpoint, form)
	if err != nil {
		return Withdrawal{}, err
	}

	var withdrawal Withdrawal
	err = json.Unmarshal(body, &withdrawal)

	if err != nil {
		return Withdrawal{}, err
	}

	return withdrawal, nil

}

//ListWithdrawals retrieves a list of withdrawals with given pages, filters and order.
func ListWithdrawals(c *Client, page, perPage int, where map[string]map[string]interface{}, orderBy string) ([]Withdrawal, error) {

	var withdrawals = make([]Withdrawal, 0)

	form := url.Values{}
	if page > 0 && perPage > 0 {
		form.Add("page", strconv.Itoa(page))
		form.Add("per_page", strconv.Itoa(perPage))
	}

	if len(where) > 0 {
		jBytes, err := json.Marshal(where)
		if err != nil {
			log.Errorf("errored at where: %s", err)
			return withdrawals, err
		}
		form.Add("where", string(jBytes))
	}

	if orderBy != "" {
		form.Add("order_by", orderBy)
	}

	body, err := c.request("GET", "withdrawals", form)
	//log.Debugf("\n\nbody: %s\n\n", body)
	if err != nil {
		log.Errorf("errored at body: %s", err)
		return withdrawals, err
	}

	err = json.NewDecoder(bytes.NewReader(body)).Decode(&withdrawals)

	if err != nil {
		log.Errorf("errored at unmarshal: %s", err)
		return withdrawals, err
	}

	return withdrawals, nil

}
