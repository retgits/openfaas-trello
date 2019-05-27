// Package function handles the invocations of the function to create Trello cards.
package function

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/adlio/trello"
	handler "github.com/openfaas-incubator/go-function-sdk"
)

// TrelloEvent contains all the details of the card to be created
type TrelloEvent struct {
	Card   Card    `json:"card"`
	Config *Config `json:"config,omitempty"`
}

// Card contains the details of the card itself
type Card struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Config contains the details of where the card should be created
type Config struct {
	Board *string `json:"board,omitempty"`
	List  *string `json:"list,omitempty"`
}

// getSecret reads a secret mounted to the container in OpenFaaS. The input is the name
// of the secret and the function returns either a byte array containing the secret or
// an error.
func getSecret(name string) (secretBytes []byte, err error) {
	// read from the openfaas secrets folder
	secretBytes, err = ioutil.ReadFile(fmt.Sprintf("/var/openfaas/secrets/%s", name))
	if err != nil {
		// read from the original location for backwards compatibility with openfaas <= 0.8.2
		secretBytes, err = ioutil.ReadFile(fmt.Sprintf("/run/secrets/%s", name))
	}

	return secretBytes, err
}

// unmarshalTrelloEvent takes a byte array and unmarshals that into a TrelloEvent object.
// In case any errors occur, the error is sent back.
func unmarshalTrelloEvent(data []byte) (TrelloEvent, error) {
	var r TrelloEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

// Handle a function invocation
func Handle(req handler.Request) (handler.Response, error) {
	apiKey, err := getSecret("trello-apikey")
	if err != nil {
		return handler.Response{
			Body:       []byte(fmt.Sprintf("error reading Trello API key: %s", err.Error())),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	appToken, err := getSecret("trello-apptoken")
	if err != nil {
		return handler.Response{
			Body:       []byte(fmt.Sprintf("error reading Trello app token: %s", err.Error())),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	memberName, err := getSecret("trello-member")
	if err != nil {
		return handler.Response{
			Body:       []byte(fmt.Sprintf("error reading Trello member: %s", err.Error())),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	event, err := unmarshalTrelloEvent(req.Body)
	if err != nil {
		return handler.Response{
			Body:       []byte(fmt.Sprintf("error unmarshalling Trello event: %s", err.Error())),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	trelloClient := trello.NewClient(string(apiKey), string(appToken))

	member, err := trelloClient.GetMember(string(memberName), trello.Defaults())
	if err != nil {
		return handler.Response{
			Body:       []byte(fmt.Sprintf("error creating connection to Trello: %s", err.Error())),
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	var board string
	if event.Config == nil || event.Config.Board == nil {
		trelloBoard, err := getSecret("trello-board")
		if err != nil {
			return handler.Response{
				Body:       []byte(fmt.Sprintf("no Trello board set in request or secrets: %s", err.Error())),
				StatusCode: http.StatusBadRequest,
			}, err
		}
		board = string(trelloBoard)
	} else {
		board = *event.Config.Board
	}

	var list string
	if event.Config == nil || event.Config.List == nil {
		trelloList, err := getSecret("trello-list")
		if err != nil {
			return handler.Response{
				Body:       []byte(fmt.Sprintf("no Trello list set in request or secrets: %s", err.Error())),
				StatusCode: http.StatusBadRequest,
			}, err
		}
		list = string(trelloList)
	} else {
		list = *event.Config.List
	}

	boards, err := member.GetBoards(trello.Defaults())
	if err != nil {
		return handler.Response{
			Body:       []byte(fmt.Sprintf("error getting Trello boards: %s", err.Error())),
			StatusCode: http.StatusBadRequest,
		}, err
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
		return handler.Response{
			Body:       []byte(fmt.Sprintf("error getting Trello board details: %s", err.Error())),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	lists, err := currBoard.GetLists(trello.Defaults())
	if err != nil {
		return handler.Response{
			Body:       []byte(fmt.Sprintf("error getting Trello lists: %s", err.Error())),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	var listid string
	for idx := range lists {
		if lists[idx].Name == list {
			listid = lists[idx].ID
			break
		}
	}

	card := &trello.Card{
		Name:    event.Card.Title,
		Desc:    event.Card.Description,
		IDList:  listid,
		IDBoard: boardid,
	}

	err = trelloClient.CreateCard(card, trello.Defaults())
	if err != nil {
		return handler.Response{
			Body:       []byte(fmt.Sprintf("error creating Trello card: %s", err.Error())),
			StatusCode: http.StatusBadRequest,
		}, err
	}

	return handler.Response{
		Body:       []byte("card successfully created"),
		StatusCode: http.StatusOK,
	}, nil
}
