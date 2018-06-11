package qvo

import "time"

//Card struct to represent a qvo card object.
type Card struct {
	ID           string    `json:"id"`
	Lats4Digits  string    `json:"last_4_digits"`
	CardType     string    `json:"card_type"`    //VISA or MASTERCARD
	PaymentType  string    `json:"payment_type"` //CD (credit) or DB (debit)
	FailureCount int32     `json:"failure_count"`
	CreatedAt    time.Time `json:"created_at"`
}
