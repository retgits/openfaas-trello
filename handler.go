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
	"net/http"
	"os"

	"github.com/adlio/trello"
	"github.com/hashicorp/vault/api"
	handler "github.com/openfaas-incubator/go-function-sdk"
	"github.com/tidwall/gjson"
)

const (
	trelloSecretPath = "secret/trello"
	trelloKeyName    = "appkey"
	trelloTokenName  = "apptoken"
	trelloMemberName = "membername"
)

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	vaultAddress := os.Getenv("VAULT_ADDRESS")
	vaultToken := os.Getenv("VAULT_TOKEN")

	if len(vaultAddress) == 0 || len(vaultToken) == 0 {
		return errorResponse("vault address or token or trello member name not set", http.StatusBadRequest)
	}

	jsonstring := string(req.Body)
	title := gjson.Get(jsonstring, "title").String()
	description := gjson.Get(jsonstring, "description").String()
	board := gjson.Get(jsonstring, "board").String()
	list := gjson.Get(jsonstring, "list").String()

	if len(title) == 0 || len(description) == 0 || len(board) == 0 || len(list) == 0 {
		return errorResponse("invalid request, all fields are mandatory", http.StatusBadRequest)
	}

	secrets, err := getSecretFromVault(vaultAddress, vaultToken, trelloSecretPath)
	if err != nil {
		return errorResponse("could not locate Trello secrets", http.StatusBadRequest)
	}

	trelloClient := trello.NewClient(secrets[trelloKeyName].(string), secrets[trelloTokenName].(string))

	member, err := trelloClient.GetMember(secrets[trelloMemberName].(string), trello.Defaults())
	if err != nil {
		return errorResponse(err.Error(), http.StatusBadRequest)
	}

	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {
		return errorResponse(err.Error(), http.StatusBadRequest)
	}

	var boardid string
	for idx := range boards {
		if boards[idx].Name == board {
			boardid = boards[idx].ID
			break
		}
	}

	currBoard, err := trelloClient.GetBoard(boardid, trello.Defaults())
	if err != nil {
		return errorResponse(err.Error(), http.StatusBadRequest)
	}

	lists, err := currBoard.GetLists(trello.Defaults())
	if err != nil {
		return errorResponse(err.Error(), http.StatusBadRequest)
	}

	var listid string
	for idx := range lists {
		if lists[idx].Name == list {
			listid = lists[idx].ID
			break
		}
	}

	card := &trello.Card{
		Name:    title,
		Desc:    description,
		IDList:  listid,
		IDBoard: boardid,
	}

	err = trelloClient.CreateCard(card, trello.Defaults())
	if err != nil {
		return errorResponse(err.Error(), http.StatusBadRequest)
	}

	return handler.Response{
		Body:       []byte("card successfully created"),
		StatusCode: http.StatusOK,
	}, nil
}

func errorResponse(body string, statusCode int) (handler.Response, error) {
	return handler.Response{
		Body:       []byte(body),
		StatusCode: statusCode,
	}, fmt.Errorf("invalid request")
}

func getSecretFromVault(address string, token string, path string) (map[string]interface{}, error) {
	conf := &api.Config{
		Address: address,
	}

	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	client.SetToken(token)

	c := client.Logical()

	secret, err := c.Read(path)
	if err != nil {
		return nil, err
	}

	if secret.Data == nil {
		return nil, fmt.Errorf("element does not exist")
	}

	return secret.Data, nil
}
