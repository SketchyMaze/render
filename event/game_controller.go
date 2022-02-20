package event

// GameController holds event state for one or more (Xbox-style) controllers.
type GameController interface {
	ID() int      // Usually the controller index number
	Name() string // Friendly name of the controller

	// State setters, to be called by the engine.
	// Note: button names are implementation-specific, use Button*() methods in your code.
	SetButtonState(name string, pressed bool)
	GetButtonState(name string) bool
	SetAxisState(name string, value int) // value maybe -32768 to 32767

	// State getters.
	ButtonA() bool
	ButtonB() bool
	ButtonX() bool
	ButtonY() bool
	ButtonL1() bool // Left shoulder
	ButtonR1() bool // Right shoulder
	ButtonL2() bool // Left trigger (digital)
	ButtonR2() bool // Right trigger (digital)
	ButtonLStick() bool
	ButtonRStick() bool
	ButtonStart() bool
	ButtonSelect() bool // Back button
	ButtonHome() bool   // Guide button

	// D-Pad buttons.
	ButtonUp() bool
	ButtonLeft() bool
	ButtonRight() bool
	ButtonDown() bool

	// Axis getters. Returns Vectors ranging from -1.0 to 1.0 being a
	// percentage of the axis between neutral and maxed out.
	LeftStick() Vector
	RightStick() Vector
	LeftTrigger() float64
	RightTrigger() float64
}
