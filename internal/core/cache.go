package core

import (
	"cache/internal/model"
	"log"
)

type Cache struct {
	store       map[model.Key]model.Value
	CommandChan chan *model.Command
}

func New() *Cache {
	service := &Cache{
		store:       map[model.Key]model.Value{},
		CommandChan: make(chan *model.Command, 5000),
	}

	service.InitReceiver()

	return service
}

func (c *Cache) InitReceiver() {
	go func() {
		for cmd := range c.CommandChan {
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

func (c *Cache) get(cmd *model.Command) {

}

func (c *Cache) set(cmd *model.Command) {

}

func (c *Cache) delete(cmd *model.Command) {

}
