package core

type Actor struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name" validate:"required"`
	Sex  rune   `json:"sex" validate:"required"`
	Bd   string `json:"bd" validate:"required"`
}
