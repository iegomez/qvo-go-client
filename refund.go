package qvo

import "time"

//Refund struct to represent a qvo refund object.
type Refund struct {
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
