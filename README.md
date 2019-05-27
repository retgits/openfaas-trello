# openfaas-trello

An [OpenFaaS](https://openfaas.com) function to create [Trello](https://trello.com) cards.

## Deploy

### Secrets

To deploy this app to OpenFaaS you'll need to have a few secrets ready:

| Name            | Description                               |
|-----------------|-------------------------------------------|
| trello-apikey   | The API key to connect to Trello          |
| trello-apptoken | The app token to connect to Trello        |
| trello-member   | Your Trello name                          |
| trello-board    | The default board name to create cards on |
| trello-list     | The default list name to create cards on  |

### Template

This app makes use of a custom template: `faas-cli template pull https://github.com/retgits/of-templates`

## Sample message

To create a card on the default board and list the input to the function should be

```json
{
  "card": {
    "title": "Hello World",
    "description": "This is pretty awesome"
  }
}
```

To override the default board and list add a `config` element to the JSON payload

```json
{
  "card": {
    "title": "Hello World",
    "description": "This is pretty awesome"
  },
  "config": {
      "board": "name of the board",
      "list": "name of the list"
  }
}
```