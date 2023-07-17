package game

type Game interface {
	Round()
	Start()
	IsOver() bool
}
