package bellows

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type StorageEngine struct {
	MutationChan chan *Mutation
	db           *gorm.DB
	eve          *EventEngine
}

type Storable interface {
	Storable()
}

type Mutation struct {
	Data    Storable
	ErrChan chan error
}

func NewStorageEngine(conf *Config, eve *EventEngine) *StorageEngine {
	db, err := gorm.Open(conf.Storage.Driver, conf.Storage.Dsn)
	if err != nil {
		log.Panicf("can't connect to db: %v", err)
	} else {
		log.Println("connected to db")
	}

	mc := make(chan *Mutation, conf.Channels.MutationQueueDepth)
	ste := &StorageEngine{mc, &db, eve}
	ste.Start()

	return ste
}

func (ste *StorageEngine) Start() {
	log.Println("starting storage engine")
	go ste.runMutator()
}

func (ste *StorageEngine) runMutator() {
	for {
		mut := <-ste.MutationChan
		log.Printf("got mutation: %v", mut.Data)
		//s.db.Begin()
		//defer s.db.Commit()
		ste.db.Save(mut.Data)
		mut.ErrChan <- ste.db.Error //race condition - needs to be deferred or something

	}
}

func (ste *StorageEngine) Store(dat Storable) error {
	errChan := make(chan error)
	mut := &Mutation{dat, errChan}
	ste.MutationChan <- mut
	err := <-errChan
	return err
}

func (ste *StorageEngine) Get(dat Storable) error {
	q := ste.db.First(dat)
	if q.Error != nil {
		return q.Error
	}

	return nil
}
