package server

type Message struct {
	User    string `json:"name"`
	Message string `json:"message"`
}
