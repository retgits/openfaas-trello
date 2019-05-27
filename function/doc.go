// Package function handles the invocations of the function to create Trello cards.
//
// The input to the function is a JSON object
//
// {
// 	"card": {
// 	  "title": "Hello World",
// 	  "description": "This is pretty awesome"
// 	},
// 	"config": {
// 	  "board": "x",
// 	  "list": "y"
// 	}
// }
//
// The card object contains the details of the card to be created and the config object
// contains details of where (like board and list) the card should be created. The
// function takes defaults from secrets created through the OpenFaaS CLI in case the
// config object doesn't exist.
//
// To connect to trello the function will get the appkey and apptoken from secrets as
// well. Details on how to get the Trello tokens can be found in the
// Trello API documentation https://trello.readme.io/docs/get-started.
package function
