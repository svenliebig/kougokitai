package persistence

type Persistence interface {
	UserExists(id string) (bool, error)
	CreateUser(id string) error
}
