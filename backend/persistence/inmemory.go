package persistence

import "log"

type inmemory struct {
	users []string
}

func newInMemory() Persistence {
	return &inmemory{}
}

func (i *inmemory) UserExists(id string) (bool, error) {
	log.Printf("checking if user '%s' exists", id)
	for _, u := range i.users {
		if u == id {
			log.Printf("user '%s' exists", id)
			return true, nil
		}
	}
	log.Printf("user '%s' does not exist", id)
	return false, nil
}

func (i *inmemory) CreateUser(id string) error {
	log.Printf("creating user '%s'", id)
	i.users = append(i.users, id)
	return nil
}
