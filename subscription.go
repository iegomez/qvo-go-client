package qvo

import "time"

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
