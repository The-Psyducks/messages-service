package usersConnector

type Interface interface {
	CheckUserExists(id string, header string) (bool, error)
}
