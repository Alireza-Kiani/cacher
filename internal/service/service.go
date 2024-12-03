package service

import (
	"cache/internal/model"
	"log"
	"sync"
	"time"
)

type CacheService struct {
	store       map[model.Key]model.Value
	commandChan chan model.Command
	m           sync.Mutex
}

func New() *CacheService {
	service := &CacheService{
		store:       map[model.Key]model.Value{},
		commandChan: make(chan model.Command, 5000),
		m:           sync.Mutex{},
	}

	service.initReceiver()
	service.initRemover()

	return service
}

func (c *CacheService) initRemover() {
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for range ticker.C {
			c.m.Lock()
			for k, v := range c.store {
				if v.TTL.Before(time.Now()) {
					c.ReceiveCommand(model.NewDeleteCommand(k))
				}
			}
			c.m.Unlock()
		}
	}()
}

func (c *CacheService) initReceiver() {
	go func() {
		for cmd := range c.commandChan {
			switch cmd.Type {
			case model.CommandSet:
				c.set(cmd)
			case model.CommandDelete:
				c.delete(cmd)
			case model.CommandGet:
				c.get(cmd)
			default:
				log.Println("command type is not supported", "given type:", cmd.Type)
			}
		}
	}()
}

func (c *CacheService) ReceiveCommand(cmd model.Command) model.Result {
	cmd.InitiateResultChannel()
	c.commandChan <- cmd
	return cmd.ResultChan
}

func (c *CacheService) get(cmd model.Command) {
	c.m.Lock()
	v, ok := c.store[cmd.Key]
	c.m.Unlock()
	if !ok {
		cmd.ResultChan <- model.PossibleValue{
			Ok: ok,
		}
		return
	}

	cmd.ResultChan <- model.PossibleValue{
		Value: v,
		Ok:    ok,
	}
}

func (c *CacheService) set(cmd model.Command) {
	c.m.Lock()
	c.store[cmd.Key] = cmd.Value
	c.m.Unlock()

	cmd.ResultChan <- model.PossibleValue{
		Ok: true,
	}
}

func (c *CacheService) delete(cmd model.Command) {
	c.m.Lock()
	delete(c.store, cmd.Key)
	c.m.Unlock()
	cmd.ResultChan <- model.PossibleValue{
		Ok: true,
	}
}
