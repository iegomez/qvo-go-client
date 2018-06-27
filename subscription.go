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

//Subscription struct to represent a qvo subscription object.
type Subscription struct {
	ID                 string        `json:"id"`
	Status             string        `json:"status"` //One of: active, canceled, trialing, retrying, inactive, unpaid.
	Debt               int64         `json:"debt"`
	Start              time.Time     `json:"start"`
	End                time.Time     `json:"end"`
	CycleCount         int32         `json:"cycle_count"`
	CurrentPeriodStart time.Time     `json:"current_period_start"`
	CurrentPeriodEnd   time.Time     `json:"current_period_end"`
	Customer           Customer      `json:"customer"`
	Plan               Plan          `json:"plan"`
	Transactions       []Transaction `json:"transactions"`
	TaxName            string        `json:"tax_name"`
	TaxPercent         float64       `json:"tax_percent"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
}

//CreateSubscription creates a subscription for a customer and plan. Returns a copy of the subscription if successful, and an error if not.
//customerID and planID are required.
//cycleCount <= 0 will be omitted.
//If taxName is "" or taxPercent isn´t in [0.0, 100.0], they'll be omitted.
func CreateSubscription(c *Client, customerID, planID, start, taxName string, taxPercent float64, cycleCount int64) (Subscription, error) {

	var subscription Subscription

	//Validate required fields.
	if customerID == "" {
		return Subscription{}, errors.New("can't create a subscription without customer id")
	}

	if planID == "" {
		return Subscription{}, errors.New("can't create a subscription without a plan id")
	}

	form := url.Values{}
	form.Add("customer_id", customerID)
	form.Add("plan_id", planID)
	if start != "" {
		form.Add("start", start)
	}
	if cycleCount > 0 {
		form.Add("cycle_count", strconv.FormatInt(cycleCount, 10))
	}
	if taxName != "" && taxPercent >= 0 && taxPercent <= 100 {
		form.Add("tax_name", taxName)
		form.Add("tax_percent", strconv.FormatFloat(taxPercent, 'f', -1, 64))
	}

	body, err := c.request("POST", "subscriptions", form)
	if err != nil {
		return Subscription{}, err
	}

	err = json.Unmarshal(body, &subscription)

	if err != nil {
		return Subscription{}, err
	}

	return subscription, nil
}

//GetSubscription returns the subscription or an error.
func GetSubscription(c *Client, subscriptionID string) (Subscription, error) {
	endpoint := fmt.Sprintf("subscriptions/%s", subscriptionID)

	form := url.Values{}
	form.Add("subscription_id", subscriptionID)

	body, err := c.request("GET", endpoint, form)
	if err != nil {
		return Subscription{}, err
	}

	var subscription Subscription
	err = json.Unmarshal(body, &subscription)

	if err != nil {
		return Subscription{}, err
	}

	return subscription, nil
}

//UpdateSubscription updates a subscription's plan given its id.
func UpdateSubscription(c *Client, subscriptionID, planID string) (Subscription, error) {

	endpoint := fmt.Sprintf("subscriptions/%s", subscriptionID)

	form := url.Values{}
	form.Add("subscription_id", subscriptionID)
	form.Add("plan_id", planID)

	body, err := c.request("PUT", endpoint, form)
	if err != nil {
		return Subscription{}, err
	}

	var subscription Subscription
	err = json.Unmarshal(body, &subscription)

	if err != nil {
		return Subscription{}, err
	}

	return subscription, nil

}

//CancelSubscription cancels a subscription.
//Depending on cancelAtPeriodEnd, it'll be canceled when the current period end is reached (if true), or immediately (if false).
//If subscription was ianctive, it'll be canceled immediately anyway.
func CancelSubscription(c *Client, subscriptionID string, cancelAtePeriodEnd bool) error {

	endpoint := fmt.Sprintf("subscriptions/%s", subscriptionID)

	form := url.Values{}
	form.Add("subscription_id", subscriptionID)
	form.Add("cancel_at_period_end", strconv.FormatBool(cancelAtePeriodEnd))

	_, err := c.request("DELETE", endpoint, form)
	if err != nil {
		return err
	}

	return nil

}

//ListSubscriptions retrieves a list of subscriptions with given pages, filters and order.
func ListSubscriptions(c *Client, page, perPage int, where map[string]map[string]interface{}, orderBy string) ([]Subscription, error) {

	var subscriptions = make([]Subscription, 0)

	form := url.Values{}
	if page > 0 && perPage > 0 {
		form.Add("page", strconv.Itoa(page))
		form.Add("per_page", strconv.Itoa(perPage))
	}

	if len(where) > 0 {
		jBytes, err := json.Marshal(where)
		if err != nil {
			log.Errorf("errored at where: %s", err)
			return subscriptions, err
		}
		form.Add("where", string(jBytes))
	}

	if orderBy != "" {
		form.Add("order_by", orderBy)
	}

	body, err := c.request("GET", "subscriptions", form)
	//log.Debugf("\n\nbody: %s\n\n", body)
	if err != nil {
		log.Errorf("errored at body: %s", err)
		return subscriptions, err
	}

	err = json.NewDecoder(bytes.NewReader(body)).Decode(&subscriptions)

	if err != nil {
		log.Errorf("errored at unmarshal: %s", err)
		return subscriptions, err
	}

	return subscriptions, nil

}
