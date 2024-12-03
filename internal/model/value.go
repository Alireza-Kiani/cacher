package model

import (
	"reflect"
	"time"
)

type Value struct {
	TTL             time.Time
	EncodedRawValue any
}

type StoredValue struct {
	Raw  []byte
	Kind reflect.Kind
}

type PossibleValue struct {
	Value Value
	Ok    bool
	Error error
}
