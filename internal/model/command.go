package model

import "time"

type CommandType int8

const (
	CommandSet CommandType = iota + 1
	CommandDelete
	CommandGet
)

type Command struct {
	Type  CommandType
	Key   Key
	Value Value
	// Core writes and controller reads
	ResultChan Result
}

func NewSetCommand(k Key, rawValue any, ttl int) Command {
	ttlTime := time.Time{}
	if ttl > 0 {
		ttlTime = time.Now().Add(time.Millisecond * time.Duration(ttl))
	}

	return Command{
		Type: CommandSet,
		Key:  k,
		Value: Value{
			TTL:             ttlTime,
			EncodedRawValue: rawValue,
		},
	}
}

func NewGetCommand(k Key) Command {
	return Command{
		Type: CommandGet,
		Key:  k,
	}
}

func NewDeleteCommand(k Key) Command {
	return Command{
		Type: CommandDelete,
		Key:  k,
	}
}

func (c *Command) InitiateResultChannel() {
	c.ResultChan = make(Result, 1)
}
