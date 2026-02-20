package models

type Res[T any] struct {
	Status  uint16 `json:"status"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}
