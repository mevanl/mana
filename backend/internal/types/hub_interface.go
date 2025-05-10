package types

type HubInterface interface {
	BroadcastMessage(event Event)
	UnregisterClient(client *Client)
	RegisterClient(client *Client)
}
