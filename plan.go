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

//Plan struct to represent a qvo plan object.
type Plan struct {
	ID                string         `json:"id"`
	Name              string         `json:"name"`
	Price             float64        `json:"price"`
	Currency          string         `json:"currency"` //CLP or UF.
	Interval          string         `json:"interval"` //One of: day, week, month, year.
	IntervalCount     int32          `json:"interval_count"`
	TrialPeriodDays   int32          `json:"trial_period_days"`
	DefaultCycleCount int32          `json:"default_cycle_count"`
	Status            string         `json:"status"` //active or inactive.
	Subscription      []Subscription `json:"subscriptions"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

//CreatePlan creates a plan at QVOs end. Returns a copy of the plan if successful, and an error if not.
func CreatePlan(c *Client, plan Plan) (Plan, error) {
	//Validate required fields.
	if plan.ID == "" || plan.Name == "" {
		return Plan{}, errors.New("can't create a plan without id or name")
	}

	if plan.Currency == "" {
		return Plan{}, errors.New("can't create a plan without currency")
	}

	form := url.Values{}
	form.Add("id", plan.ID)
	form.Add("name", plan.Name)
	form.Add("price", strconv.FormatFloat(plan.Price, 'f', -1, 64))
	form.Add("currency", plan.Currency)
	form.Add("interval", plan.Interval)
	form.Add("interval_count", strconv.FormatInt(int64(plan.IntervalCount), 10))
	form.Add("trial_period_days", strconv.FormatInt(int64(plan.TrialPeriodDays), 10))
	form.Add("default_cycle_count", strconv.FormatInt(int64(plan.DefaultCycleCount), 10))

	body, err := c.request("POST", "plans", form)
	if err != nil {
		return Plan{}, err
	}

	err = json.Unmarshal(body, &plan)

	if err != nil {
		return Plan{}, err
	}

	return plan, nil
}

//GetPlan retrieves a plan by id.
func GetPlan(c *Client, id string) (Plan, error) {

	endpoint := fmt.Sprintf("plans/%s", id)

	form := url.Values{}
	form.Add("plan_id", id)

	body, err := c.request("GET", endpoint, form)
	if err != nil {
		return Plan{}, err
	}

	var plan Plan
	err = json.Unmarshal(body, &plan)

	if err != nil {
		return Plan{}, err
	}

	return plan, nil

}

//UpdatePlan updates a plan given its id.
func UpdatePlan(c *Client, plan Plan) (Plan, error) {

	endpoint := fmt.Sprintf("plans/%s", plan.ID)

	form := url.Values{}
	form.Add("plan_id", plan.ID)
	form.Add("name", plan.Name)
	form.Add("price", strconv.FormatFloat(plan.Price, 'f', -1, 64))
	form.Add("currency", plan.Currency)
	form.Add("interval", plan.Interval)
	form.Add("interval_count", strconv.FormatInt(int64(plan.IntervalCount), 10))
	form.Add("trial_period_days", strconv.FormatInt(int64(plan.TrialPeriodDays), 10))
	form.Add("default_cycle_count", strconv.FormatInt(int64(plan.DefaultCycleCount), 10))

	body, err := c.request("PUT", endpoint, form)
	if err != nil {
		return Plan{}, err
	}

	err = json.Unmarshal(body, &plan)

	if err != nil {
		return Plan{}, err
	}

	return plan, nil

}

//DeletePlan deletes a plan given its id.
func DeletePlan(c *Client, id string) error {

	endpoint := fmt.Sprintf("plans/%s", id)

	form := url.Values{}
	form.Add("plan_id", id)

	_, err := c.request("DELETE", endpoint, form)
	if err != nil {
		return err
	}

	return nil

}

//ListPlans retrieves a list of plans with given pages, filters and order.
func ListPlans(c *Client, page, perPage int, where map[string]map[string]interface{}, orderBy string) ([]Plan, error) {

	var plans = make([]Plan, 0)

	form := url.Values{}
	if page > 0 && perPage > 0 {
		form.Add("page", strconv.Itoa(page))
		form.Add("per_page", strconv.Itoa(perPage))
	}

	if len(where) > 0 {
		jBytes, err := json.Marshal(where)
		if err != nil {
			log.Errorf("errored at where: %s", err)
			return plans, err
		}
		form.Add("where", string(jBytes))
	}

	if orderBy != "" {
		form.Add("order_by", orderBy)
	}

	body, err := c.request("GET", "plans", form)
	//log.Debugf("\n\nbody: %s\n\n", body)
	if err != nil {
		log.Errorf("errored at body: %s", err)
		return plans, err
	}

	err = json.NewDecoder(bytes.NewReader(body)).Decode(&plans)

	if err != nil {
		log.Errorf("errored at unmarshal: %s", err)
		return plans, err
	}

	return plans, nil

}
