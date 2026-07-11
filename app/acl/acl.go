package acl

type Acl struct {
	storage *Storage
}

func New() *Acl {
	storage := NewStorage()
	return &Acl{
		storage: storage,
	}
}

func (a *Acl) GetUser(name string) *User {
	u, ok := a.storage.data[name]
	if !ok {
		return nil
	}

	return u
}