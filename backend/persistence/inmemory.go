package persistence

type inmemory struct {
	users []string
}

func newInMemory() Persistence {
	return &inmemory{}
}

func (i *inmemory) UserExists(email string) (bool, error) {
	for _, u := range i.users {
		if u == email {
			return true, nil
		}
	}
	return false, nil
}

func (i *inmemory) CreateUser(email string) error {
	i.users = append(i.users, email)
	return nil
}
