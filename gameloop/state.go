package gameloop

type GameHandler interface {
	// Integrate handles a logical step in the game, it must return the next state of the game
	// currentState is the current state of the game.
	// t is the time in seconds
	// deltaT is the time since the last update
	Integrate(currentState interface{}, t float64, deltaT float64) interface{}

	// Render should handle all the Rendering logic of the game.
	// _note:_ only display logic should go here
	// state is the current state of the game.
	// t is the time in seconds
	// alpha is the progression in between display frames. This allows for liner interpolation.
	Render(state interface{}, t float64, alpha float64)
}