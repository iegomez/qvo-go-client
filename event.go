package qvo

import "time"

//Event struct to represent a qvo event object.
type Event struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} //API sends a "hash", so we are limited to an interfaces map.
	previous  map[string]interface{} //API sends a "hash", so we are limited to an interfaces map.
	CreatedAt time.Time              `json:"created_at"`
}
