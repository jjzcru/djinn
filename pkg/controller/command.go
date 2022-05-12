package controller

type Command struct {
	Id string `json:"id"`
	Plugin string `json:"plugin"`
	Command string `json:"command"`
	Payload string `json:"payload"`
}