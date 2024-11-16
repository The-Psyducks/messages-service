package users_connector

import (
	"github.com/stretchr/testify/assert"
	"messages/src/auth"
	"testing"
)

//This tests don't work if there is no user service running

func TestUsersConnector_GetUserNameAndImage(t *testing.T) {
	uc := NewUsersConnector()

	token, err := auth.GenerateToken("be89833a-836b-4f1d-ab8f-1df1b9789864", "Monke", false)
	if err != nil {
		t.Error(err)
	}

	_, _, err = uc.GetUserNameAndImage("be89833a-836b-4f1d-ab8f-1df1b9789864", "Bearer "+token)
	assert.Nil(t, err)
}

func TestUsersConnector_CheckUserExists(t *testing.T) {
	uc := NewUsersConnector()

	token, err := auth.GenerateToken("be89833a-836b-4f1d-ab8f-1df1b9789864", "Monke", false)
	if err != nil {
		t.Error(err)
	}

	userExists, err := uc.CheckUserExists("be89833a-836b-4f1d-ab8f-1df1b9789864", "Bearer "+token)
	assert.Nil(t, err)

	assert.True(t, userExists)

	//assert.True(t, userExists)

}
