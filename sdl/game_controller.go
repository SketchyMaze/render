package sdl

import (
	"git.kirsle.net/go/render/event"
	"github.com/veandco/go-sdl2/sdl"
)

// User tuneable properties.
var (
	// For controllers having a digital (non-analog) Left/Right Trigger, the press percentage
	// for which to consider it a boolean press.
	TriggerAxisBooleanThreshold float64 = 0.5
)

// GameController holds an abstraction around SDL2 GameControllers.
type GameController struct {
	id     int
	name   string
	active bool

	// Underlying SDL2 GameController.
	ctrl *sdl.GameController

	// Button states.
	buttons map[string]bool
	axes    map[string]int
}

// NewGameController creates a GameController from an SDL2 controller.
func NewGameController(index int, name string, ctrl *sdl.GameController) *GameController {
	return &GameController{
		id:      index,
		name:    name,
		ctrl:    ctrl,
		buttons: map[string]bool{},
		axes:    map[string]int{},
	}
}

// ID returns the controller index as SDL2 knows it.
func (gc *GameController) ID() int {
	return gc.id
}

// Name returns the controller name.
func (gc *GameController) Name() string {
	return gc.name
}

// SetButtonState sets the state using the SDL2 button names.
func (gc *GameController) SetButtonState(name string, pressed bool) {
	gc.buttons[name] = pressed
}

// GetButtonState returns the button state by SDL2 button name.
func (gc *GameController) GetButtonState(name string) bool {
	if v, ok := gc.buttons[name]; ok {
		return v
	}
	return false
}

// SetAxisState sets the axis state.
func (gc *GameController) SetAxisState(name string, value int) {
	gc.axes[name] = value
}

// GetAxisState returns the underlying SDL2 axis state.
func (gc *GameController) GetAxisState(name string) int {
	if v, ok := gc.axes[name]; ok {
		return v
	}
	return 0
}

// ButtonA returns whether the logical Xbox button is pressed.
func (gc *GameController) ButtonA() bool {
	return gc.GetButtonState("a")
}

// ButtonB returns whether the logical Xbox button is pressed.
func (gc *GameController) ButtonB() bool {
	return gc.GetButtonState("b")
}

// ButtonX returns whether the logical Xbox button is pressed.
func (gc *GameController) ButtonX() bool {
	return gc.GetButtonState("x")
}

// ButtonY returns whether the logical Xbox button is pressed.
func (gc *GameController) ButtonY() bool {
	return gc.GetButtonState("y")
}

// ButtonL1 returns whether the Left Shoulder button is pressed.
func (gc *GameController) ButtonL1() bool {
	return gc.GetButtonState("leftshoulder")
}

// ButtonR1 returns whether the Right Shoulder button is pressed.
func (gc *GameController) ButtonR1() bool {
	return gc.GetButtonState("rightshoulder")
}

// ButtonL2 returns whether the Left Trigger (digital) button is pressed.
// Returns true if the LeftTrigger is 50% pressed or TriggerAxisBooleanThreshold.
func (gc *GameController) ButtonL2() bool {
	return gc.axisToFloat("lefttrigger") > TriggerAxisBooleanThreshold
}

// ButtonR2 returns whether the Left Trigger (digital) button is pressed.
// Returns true if the LeftTrigger is 50% pressed or TriggerAxisBooleanThreshold.
func (gc *GameController) ButtonR2() bool {
	return gc.axisToFloat("righttrigger") > TriggerAxisBooleanThreshold
}

// ButtonLStick returns whether the Left Stick button is pressed.
func (gc *GameController) ButtonLStick() bool {
	return gc.GetButtonState("leftstick")
}

// ButtonRStick returns whether the Right Stick button is pressed.
func (gc *GameController) ButtonRStick() bool {
	return gc.GetButtonState("rightstick")
}

// ButtonStart returns whether the logical Xbox button is pressed.
func (gc *GameController) ButtonStart() bool {
	return gc.GetButtonState("start")
}

// ButtonSelect returns whether the Xbox "back" button is pressed.
func (gc *GameController) ButtonSelect() bool {
	return gc.GetButtonState("back")
}

// ButtonUp returns whether the Xbox D-Pad button is pressed.
func (gc *GameController) ButtonUp() bool {
	return gc.GetButtonState("dpup")
}

// ButtonDown returns whether the Xbox D-Pad button is pressed.
func (gc *GameController) ButtonDown() bool {
	return gc.GetButtonState("dpdown")
}

// ButtonLeft returns whether the Xbox D-Pad button is pressed.
func (gc *GameController) ButtonLeft() bool {
	return gc.GetButtonState("dpleft")
}

// ButtonRight returns whether the Xbox D-Pad button is pressed.
func (gc *GameController) ButtonRight() bool {
	return gc.GetButtonState("dpright")
}

// ButtonHome returns whether the Xbox "guide" button is pressed.
func (gc *GameController) ButtonHome() bool {
	return gc.GetButtonState("guide")
}

// LeftStick returns the vector of X and Y of the left analog stick.
func (gc *GameController) LeftStick() event.Vector {
	return event.Vector{
		X: gc.axisToFloat("leftx"),
		Y: gc.axisToFloat("lefty"),
	}
}

// RightStick returns the vector of X and Y of the right analog stick.
func (gc *GameController) RightStick() event.Vector {
	return event.Vector{
		X: gc.axisToFloat("rightx"),
		Y: gc.axisToFloat("righty"),
	}
}

// LeftTrigger returns the vector of the left analog trigger.
func (gc *GameController) LeftTrigger() float64 {
	return gc.axisToFloat("lefttrigger")
}

// RightTrigger returns the vector of the left analog trigger.
func (gc *GameController) RightTrigger() float64 {
	return gc.axisToFloat("righttrigger")
}

// axisToFloat converts an SDL2 Axis value to a float between -1.0..1.0
func (gc *GameController) axisToFloat(name string) float64 {
	// SDL2 Axis is an int16 ranging from -32768 to 32767,
	// convert this into a percentage +- 0 to 1.
	axis := gc.GetAxisState(name)
	if axis < 0 {
		return float64(axis) / 32768
	} else {
		return float64(axis) / 32767
	}
}
