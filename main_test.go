package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/models"

	"github.com/stretchr/testify/assert"
)



func TestNewUserHandler(t *testing.T) {
	server := httptest.NewServer(SetupServer())
	defer server.Close()

	user := models.User{
		FirstName: "Selcuk",
		LastName: "Sozuer",
		Email: "selcuk.sozuer@gmail.com",
		Password: "1234",		
	}

	raw, _ := json.Marshal(user)
	res, err := http.Post(fmt.Sprintf("%s/users", server.URL), "application/json", bytes.NewBuffer(raw))
	assert.Nil(t, err)

	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	data, _ := ioutil.ReadAll(res.Body)
	var newUser models.User
	json.Unmarshal(data, &newUser)
	assert.Equal(t, newUser.FirstName, "Selcuk")
	assert.Equal(t, newUser.Email, "selcuk.sozuer@gmail.com")
}

func TestListUsersHandler(t *testing.T) {
	server := httptest.NewServer(SetupServer())
	defer server.Close()
	
	res, err := http.Get(fmt.Sprintf("%s/users", server.URL))
	assert.Nil(t, err)

	defer res.Body.Close()
	assert.Equal(t, http.StatusOK, res.StatusCode)

	data, _ := ioutil.ReadAll(res.Body)
	var users []models.User
	json.Unmarshal(data, &users)
	assert.True(t, len(users) > 0)
}