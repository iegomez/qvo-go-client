package qvo

import "time"

//Client struct to represent a qvo client object.
type Client struct {
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
