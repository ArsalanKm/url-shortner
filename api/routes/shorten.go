package routes

import (
	"time"
)

type request struct {
	URL         		string        `json:"url"`
	CustomShort 		string        `json:"short"`
	Expiry      		time.Duration `json:"expiry"`
}

type response struct {
	URL            		string `json:"url"`
	CustomShort    		string `json:"custom_short"`
	Expiry         		string `json:"expiry"`
	XrateRemaining 		string `json:"rate_remaining"`
	XrateLimitRest 		string `json:"rate_limit_rest"`
}
