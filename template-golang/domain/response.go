package domain

import "time"

type Response struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}
