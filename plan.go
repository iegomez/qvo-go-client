package qvo

import "time"

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
