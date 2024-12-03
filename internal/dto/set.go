package dto

import "errors"

type SetRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
	TTL   int    `json:"ttl"` //in milliseconds
}

func (r SetRequest) Validate() error {
	if r.TTL < 0 {
		return errors.New("ttl cannot be negative, for permanent storing, provide zero")
	}

	if len(r.Key) == 0 {
		return errors.New("key cannot be empty")
	}

	if r.Value == nil {
		return errors.New("value cannot be null")
	}

	return nil
}
