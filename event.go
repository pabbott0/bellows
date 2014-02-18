package bellows

import (
	"log"
)

type Event struct {
}

type EventEngine struct {
	EventChan chan *Event
}

func NewEventEngine(conf *Config) *EventEngine {
	ec := make(chan *Event, conf.Channels.EventQueueDepth)
	eve := &EventEngine{ec}
	eve.Start()
	return eve
}

func (eve *EventEngine) Start() {
	log.Println("starting event engine")
	go eve.handleEvents()
}

func (eve *EventEngine) handleEvents() {
	for {
		e := <-eve.EventChan
		log.Printf("got event: %v", e)
	}
}
