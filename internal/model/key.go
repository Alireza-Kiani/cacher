package model

import "time"

type Key string

type KeyTTL struct {
	Key Key
	TTL time.Time
}
