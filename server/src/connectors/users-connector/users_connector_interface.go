package users_connector

type Interface interface {
	CheckUserExists(id string, header string) (bool, error)
}
