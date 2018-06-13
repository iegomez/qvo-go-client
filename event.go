package qvo

import "time"

//Event struct to represent a qvo event object.
type Event struct {
	ID        string                  `json:"id"`
	Type      string                  `json:"type"`
	Data      map[string]interface{}  `json:"data"`               //API sends a "hash", so we are limited to an interfaces map.
	Previous  *map[string]interface{} `json:"previous,omitempty"` //API sends a "hash", so we are limited to an interfaces map. Also, it's nullable, so it's a pointer.
	CreatedAt time.Time               `json:"created_at"`
}
