package core

type Movie struct {
	Id      int    `json:"id,omitempty"`
	Title   string `json:"title" validate:"required"`
	Descr   string `json:"descr" validate:"required"`
	Release string `json:"release" validate:"required"`
	Rating  int    `json:"rating" validate:"required"`
	Actors  []int  `json:"actors" validate:"required"`
}
