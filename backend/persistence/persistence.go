package persistence

type Persistence interface {
	UserExists(id string) (bool, error)
	CreateUser(id string) error
}

var registered Persistence

func Register(p Persistence) {
	if registered == nil {
		registered = p
		return
	}

	panic("persistence: already registered, you can only register one persistence implementation.")
}
