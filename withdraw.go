package qvo

import "time"

//Withdraw struct to represent a qvo withdraw object.
type Withdraw struct {
	ID        string    `json:"id"`
	Amount    int64     `json:"amount"`
	Status    string    `json:"status"` //One of: processing, rejected, transfered.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
