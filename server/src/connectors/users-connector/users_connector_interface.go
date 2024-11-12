package users_connector

type Interface interface {
	CheckUserExists(id string, header string) (bool, error)
	GetUserNameAndImage(id string, header string) (string, string, error)
}
