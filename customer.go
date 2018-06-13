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

//Customer struct to represent a qvo customer object.
type Customer struct {
	ID                   string         `json:"id"`
	DefaultPaymentMethod Card           `json:"default_payment_method"`
	Name                 string         `json:"name"`
	Email                string         `json:"email"`
	Subscriptions        []Subscription `json:"subscriptions"`
	Cards                []Card         `json:"cards"`
	Transactions         []Transaction  `json:"transactions"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

//CreateCustomer creates a customer at the QVO account.
func CreateCustomer(c *Client, name, email string) (Customer, error) {

	form := url.Values{}
	form.Add("name", name)
	form.Add("email", email)

	body, err := c.request("POST", "customers", form)
	if err != nil {
		return Customer{}, err
	}

	var customer Customer
	err = json.Unmarshal(body, &customer)

	if err != nil {
		return Customer{}, err
	}

	return customer, nil

}

//GetCustomer retrieves a costumer given its id.
func GetCustomer(c *Client, id string) (Customer, error) {

	endpoint := fmt.Sprintf("customers/%s", id)

	form := url.Values{}
	form.Add("customer_id", id)

	body, err := c.request("GET", endpoint, form)
	if err != nil {
		return Customer{}, err
	}

	var customer Customer
	err = json.Unmarshal(body, &customer)

	if err != nil {
		return Customer{}, err
	}

	return customer, nil

}

//UpdateCustomer updates a customer given its id.
func UpdateCustomer(c *Client, id, name, email, defaultPaymentMethodID string) (Customer, error) {

	endpoint := fmt.Sprintf("customers/%s", id)

	form := url.Values{}
	form.Add("customer_id", id)
	form.Add("name", name)
	form.Add("email", email)
	form.Add("default_payment_method_id", defaultPaymentMethodID)

	body, err := c.request("PUT", endpoint, form)
	if err != nil {
		return Customer{}, err
	}

	var customer Customer
	err = json.Unmarshal(body, &customer)

	if err != nil {
		return Customer{}, err
	}

	return customer, nil

}

//DeleteCustomer deletes a customer given its id.
func DeleteCustomer(c *Client, id string) error {

	endpoint := fmt.Sprintf("customers/%s", id)

	form := url.Values{}
	form.Add("customer_id", id)

	_, err := c.request("DELETE", endpoint, form)
	if err != nil {
		return err
	}

	return nil

}

//ListCustomers retrieves a list of customers with given pages, filters and order.
func ListCustomers(c *Client, page, perPage int, where map[string]map[string]interface{}, orderBy string) ([]Customer, error) {

	var customers = make([]Customer, 0)

	form := url.Values{}
	if page > 0 && perPage > 0 {
		form.Add("page", strconv.Itoa(page))
		form.Add("per_page", strconv.Itoa(perPage))
	}

	if len(where) > 0 {
		jBytes, err := json.Marshal(where)
		if err != nil {
			log.Errorf("errored at where: %s", err)
			return customers, err
		}
		form.Add("where", string(jBytes))
	}

	if orderBy != "" {
		form.Add("order_by", orderBy)
	}

	body, err := c.request("GET", "customers", form)
	log.Debugf("\n\nbody: %s\n\n", body)
	if err != nil {
		log.Errorf("errored at body: %s", err)
		return customers, err
	}

	err = json.NewDecoder(bytes.NewReader(body)).Decode(&customers)
	//err = json.Unmarshal(body, &customers)

	log.Debugf("body: %s\n%v\n", string(body), body)

	if err != nil {
		log.Errorf("errored at unmarshal: %s", err)
		return customers, err
	}

	return customers, nil

}
