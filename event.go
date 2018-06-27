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

//Event types
const (
	CustomerCreated             string = "customer.created"
	CustomerUpdated             string = "customer.updated"
	CustomerDeleted             string = "customer.deleted"
	PlanCreated                 string = "plan.created"
	PlanUpdated                 string = "plan.updated"
	PlanDeleted                 string = "plan.deleted"
	CustomerCardCreated         string = "customer.card.created"
	CustomerCardDeleted         string = "customer.card.deleted"
	CustomerSubscriptionCreated string = "customer.subscription.created"
	CustomerSubscriptionUpdated string = "customer.subscription.updated"
	CustomerSubscriptionDeleted string = "customer.subscription.deleted"
	TransactionPaymentSucceeded string = "transaction.payment_succeeded"
	TransactionPaymentFailed    string = "transaction.payment_failed"
	TransactionPRefunded        string = "transaction.refunded"
	TransactionResponseTimeout  string = "transaction.response_timeout"
)

//Event struct to represent a qvo event object.
type Event struct {
	ID        string                  `json:"id"`
	Type      string                  `json:"type"`
	Data      map[string]interface{}  `json:"data"`               //API sends a "hash", so we are limited to an interfaces map.
	Previous  *map[string]interface{} `json:"previous,omitempty"` //API sends a "hash", so we are limited to an interfaces map. Also, it's nullable, so it's a pointer.
	CreatedAt time.Time               `json:"created_at"`
}

//GetEvent retrieves a event given its id.
func GetEvent(c *Client, id string) (Event, error) {

	endpoint := fmt.Sprintf("events/%s", id)

	form := url.Values{}
	form.Add("event_id", id)

	body, err := c.request("GET", endpoint, form)
	if err != nil {
		return Event{}, err
	}

	var event Event
	err = json.Unmarshal(body, &event)

	if err != nil {
		return Event{}, err
	}

	return event, nil

}

//ListEvents retrieves a list of events with given pages, filters and order.
func ListEvents(c *Client, page, perPage int, where map[string]map[string]interface{}, orderBy string) ([]Event, error) {

	var events = make([]Event, 0)

	form := url.Values{}
	if page > 0 && perPage > 0 {
		form.Add("page", strconv.Itoa(page))
		form.Add("per_page", strconv.Itoa(perPage))
	}

	if len(where) > 0 {
		jBytes, err := json.Marshal(where)
		if err != nil {
			log.Errorf("errored at where: %s", err)
			return events, err
		}
		form.Add("where", string(jBytes))
	}

	if orderBy != "" {
		form.Add("order_by", orderBy)
	}

	body, err := c.request("GET", "events", form)
	//log.Debugf("\n\nbody: %s\n\n", body)
	if err != nil {
		log.Errorf("errored at body: %s", err)
		return events, err
	}

	err = json.NewDecoder(bytes.NewReader(body)).Decode(&events)

	if err != nil {
		log.Errorf("errored at unmarshal: %s", err)
		return events, err
	}

	return events, nil

}
