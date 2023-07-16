package game

type Game interface {
	Round()
	IsOver() bool
}
