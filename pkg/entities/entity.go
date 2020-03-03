package entities

type Entity interface {
	GetUserPasswd() UserPasswd
}

type DefaultEntity struct {
	User *UserPasswd
}

func (e *DefaultEntity) GetUserPasswd() UserPasswd {
	return *e.User
}
