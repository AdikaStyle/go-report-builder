package models

type PrintOptions struct {
	PageHeight  string
	PageWidth   string
	Orientation PrintOrientation
}

type PrintOrientation int

const (
	Orientation_Landscape = 1
	Orientation_Portrait  = 2
)
