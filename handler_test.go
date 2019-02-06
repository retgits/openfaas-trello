// Package function handles the invocations of the function to create Trello cards.
// The input to the function is a JSON object containing a title, description,
// boardname and list name. Based on the board and list name, the function will find
// the correct identifiers and create a new Trello card. To access Trello, the function
// will get the appkey and apptoken from Vault.
// Details on how to get the Trello tokens can be found in the
// [Trello API documentation](https://trello.readme.io/docs/get-started).
package function

import (
	"fmt"
	"os"
	"testing"

	handler "github.com/openfaas-incubator/go-function-sdk"
	"github.com/stretchr/testify/assert"
)

const (
	testVaultAddress    = "."
	testVaultToken      = "."
	testTrelloMember    = "."
	testCardTitle       = "Hello"
	testCardDescription = "Hello World is the best description ever!"
	testCardBoard       = "Main"
	testCardList        = "Done"
)

func TestEnvironmentVariableSettings(t *testing.T) {
	assert := assert.New(t)

	os.Unsetenv("VAULT_ADDRESS")
	os.Unsetenv("VAULT_TOKEN")

	message := handler.Request{
		Body: []byte(`{"title":"","description":"","board":"","list":"",}`),
	}

	resp, err := Handle(message)
	assert.Error(err)
	assert.EqualValues(resp.StatusCode, 400)
	assert.EqualValues(string(resp.Body), "vault address or token or trello member name not set")

	os.Setenv("VAULT_ADDRESS", testVaultAddress)
	resp, err = Handle(message)
	assert.Error(err)
	assert.EqualValues(resp.StatusCode, 400)
	assert.EqualValues(string(resp.Body), "vault address or token or trello member name not set")

	os.Setenv("VAULT_ADDRESS", testVaultAddress)
	os.Setenv("VAULT_TOKEN", testVaultToken)
	resp, err = Handle(message)
	assert.Error(err)
	assert.EqualValues(resp.StatusCode, 400)
	assert.EqualValues(string(resp.Body), "invalid request, all fields are mandatory")
}

func TestMessageSettings(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("VAULT_ADDRESS", testVaultAddress)
	os.Setenv("VAULT_TOKEN", "DEMOTOKEN")

	message := handler.Request{
		Body: []byte(`{"title":"","description":"","board":"","list":"",}`),
	}

	resp, err := Handle(message)
	assert.Error(err)
	assert.EqualValues(resp.StatusCode, 400)
	assert.EqualValues(string(resp.Body), "invalid request, all fields are mandatory")

	message = handler.Request{
		Body: []byte(fmt.Sprintf(`{"title":"%s","description":"","board":"","list":"",}`, testCardTitle)),
	}

	resp, err = Handle(message)
	assert.Error(err)
	assert.EqualValues(resp.StatusCode, 400)
	assert.EqualValues(string(resp.Body), "invalid request, all fields are mandatory")

	message = handler.Request{
		Body: []byte(fmt.Sprintf(`{"title":"%s","description":"%s","board":"","list":"",}`, testCardTitle, testCardDescription)),
	}

	resp, err = Handle(message)
	assert.Error(err)
	assert.EqualValues(resp.StatusCode, 400)
	assert.EqualValues(string(resp.Body), "invalid request, all fields are mandatory")

	message = handler.Request{
		Body: []byte(fmt.Sprintf(`{"title":"%s","description":"%s","board":"%s","list":"",}`, testCardTitle, testCardDescription, testCardBoard)),
	}

	resp, err = Handle(message)
	assert.Error(err)
	assert.EqualValues(resp.StatusCode, 400)
	assert.EqualValues(string(resp.Body), "invalid request, all fields are mandatory")

	message = handler.Request{
		Body: []byte(fmt.Sprintf(`{"title":"%s","description":"%s","board":"%s","list":"%s",}`, testCardTitle, testCardDescription, testCardBoard, testCardList)),
	}

	resp, err = Handle(message)
	assert.Error(err)
	assert.EqualValues(resp.StatusCode, 400)
	assert.EqualValues(string(resp.Body), "could not locate Trello secrets")
}

func TestVault(t *testing.T) {
	assert := assert.New(t)

	os.Setenv("VAULT_ADDRESS", testVaultAddress)
	os.Setenv("VAULT_TOKEN", "DEMOTOKEN")

	message := handler.Request{
		Body: []byte(fmt.Sprintf(`{"title":"%s","description":"%s","board":"%s","list":"%s",}`, testCardTitle, testCardDescription, testCardBoard, testCardList)),
	}

	resp, err := Handle(message)
	assert.Error(err)
	assert.EqualValues(resp.StatusCode, 400)
	assert.EqualValues(string(resp.Body), "could not locate Trello secrets")

	os.Setenv("VAULT_ADDRESS", testVaultAddress)
	os.Setenv("VAULT_TOKEN", testVaultToken)

	message = handler.Request{
		Body: []byte(fmt.Sprintf(`{"title":"%s","description":"%s","board":"%s","list":"%s",}`, testCardTitle, testCardDescription, testCardBoard, testCardList)),
	}

	resp, err = Handle(message)
	assert.NoError(err)
	assert.EqualValues(resp.StatusCode, 200)
	assert.EqualValues(string(resp.Body), "card successfully created")
}
