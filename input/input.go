// Package input provides functions that can be used to handle input devices on your game
package input

import (
	"log"

	"github.com/ungerik/go3d/vec2"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	// Joystick
	joysticks            map[sdl.JoystickID]*sdl.Joystick = make(map[sdl.JoystickID]*sdl.Joystick)
	joystickAxises       map[sdl.JoystickID][]int16       = make(map[sdl.JoystickID][]int16)
	joystickButtons      map[sdl.JoystickID][]bool        = make(map[sdl.JoystickID][]bool)
	joystickHats         map[sdl.JoystickID][]uint8       = make(map[sdl.JoystickID][]uint8)
	joysticksInitialised bool

	// Mouse
	mouseLocation vec2.T = vec2.Zero
	mouseState    []bool = make([]bool, 3)
)

const (
	// Mouse Button Constants

	// MOUSE_LEFT is the value for the left mouse button
	MOUSE_LEFT = 0

	// MOUSE_MIDDLE is the value for the middle mouse button ( aka the scrollwheel )
	MOUSE_MIDDLE = 1

	// MOUSE_RIGHT is the value for the right mouse button
	MOUSE_RIGHT = 2

	// Joystick Hat Position Constants

	// JOYSTICK_HAT_N is the value for the north joystick hat position
	JOYSTICK_HAT_N = 1

	// JOYSTICK_HAT_NE is the value for the north east joystick hat position
	JOYSTICK_HAT_NE = 3

	// JOYSTICK_HAT_E is the value for the east joystick hat position
	JOYSTICK_HAT_E = 2

	// JOYSTICK_HAT_SE is the value for the south east joystick hat position
	JOYSTICK_HAT_SE = 6

	// JOYSTICK_HAT_S is the value for the south joystick hat position
	JOYSTICK_HAT_S = 4

	// JOYSTICK_HAT_SW is the value for the south east joystick hat position
	JOYSTICK_HAT_SW = 9

	// JOYSTICK_HAT_W is the value for the west joystick hat position
	JOYSTICK_HAT_W = 8

	// JOYSTICK_HAT_NW is the value for the north west joystick hat position
	JOYSTICK_HAT_NW = 12
)

// InitJoystick initializes the Joystick Subsystem and add available joysticks
func InitJoystick() {
	if sdl.WasInit(sdl.INIT_JOYSTICK) == 0 {
		sdl.InitSubSystem(sdl.INIT_JOYSTICK)
	}

	if sdl.NumJoysticks() > 0 {
		for i := 0; i < sdl.NumJoysticks(); i++ {
			id := sdl.JoystickID(i)

			addJoystick(id)
		}

		sdl.JoystickEventState(sdl.ENABLE)

		joysticksInitialised = true
	}
}

// HandleEvents handles input device specific events
// like keyboard input, mouse input, and joystick input
func HandleEvents(e sdl.Event) {
	switch t := e.(type) {
	case *sdl.JoyDeviceEvent:
		handleJoyDeviceEvent(t)
		break
	case *sdl.JoyAxisEvent:
		handleJoyAxisEvent(t)
		break
	case *sdl.JoyButtonEvent:
		handleJoyButtonEvent(t)
		break
	case *sdl.JoyHatEvent:
		handleJoyHatEvent(t)
		break
	case *sdl.MouseMotionEvent:
		handleMouseMotionEvent(t)
		break
	case *sdl.MouseButtonEvent:
		handleMouseButtonEvent(t)
		break
	}
}

// Joystick Event Handlers

func handleJoyDeviceEvent(t *sdl.JoyDeviceEvent) {
	if t.Type == sdl.JOYDEVICEADDED {
		addJoystick(t.Which)
	} else if t.Type == sdl.JOYDEVICEREMOVED {
		remJoystick(t.Which)
	}

}

func handleJoyAxisEvent(t *sdl.JoyAxisEvent) {
	joystickAxises[t.Which][t.Axis] = t.Value
}

func handleJoyButtonEvent(t *sdl.JoyButtonEvent) {
	if t.State == 1 {
		joystickButtons[t.Which][t.Button] = true
	} else {
		joystickButtons[t.Which][t.Button] = false
	}
}

func handleJoyHatEvent(t *sdl.JoyHatEvent) {
	joystickHats[t.Which][t.Hat] = t.Value
}

// Mouse Event Handlers

func handleMouseMotionEvent(t *sdl.MouseMotionEvent) {
	mouseLocation[0] = float32(t.X)
	mouseLocation[1] = float32(t.Y)
}

func handleMouseButtonEvent(t *sdl.MouseButtonEvent) {
	if t.Type == sdl.MOUSEBUTTONDOWN {
		if t.Button == sdl.BUTTON_LEFT {
			mouseState[MOUSE_LEFT] = true
		}

		if t.Button == sdl.BUTTON_MIDDLE {
			mouseState[MOUSE_MIDDLE] = true
		}

		if t.Button == sdl.BUTTON_RIGHT {
			mouseState[MOUSE_RIGHT] = true
		}
	} else {
		if t.Button == sdl.BUTTON_LEFT {
			mouseState[MOUSE_LEFT] = false
		}

		if t.Button == sdl.BUTTON_MIDDLE {
			mouseState[MOUSE_MIDDLE] = false
		}

		if t.Button == sdl.BUTTON_RIGHT {
			mouseState[MOUSE_RIGHT] = false
		}
	}
}

// Axis returns the Axis value for a certain Joystick's Axis.
// The returned value is a int16 number usually to show how 'full' the analog trigger are pressed or how far the stick has gone
func Axis(id sdl.JoystickID, axis uint) int16 {
	return joystickAxises[id][axis]
}

// Axisf returns the Axis value for a certain Joystick's Axis in float32.
// The value returns a float32 number that goes from -1 to 1 usually to show how 'full' the analog trigger are pressed or how far the stick has gone
func Axisf(id sdl.JoystickID, axis uint) float32 {
	return float32(Axis(id, axis)) / 65536
}

// Button returns the Joystick's button state
func Button(id sdl.JoystickID, button uint) bool {
	return joystickButtons[id][button]
}

// Hat returns the Joystick's hat position.
// Use the JOYSTICK_HAT_* constants to know what position it's on
func Hat(id sdl.JoystickID, hat uint) uint8 {
	return joystickHats[id][hat]
}

// MouseLocation returns the mouse location relative to the window in a 2D Vector
func MouseLocation() vec2.T {
	return mouseLocation
}

// Mouse returns the state of a mouse button.
// Use the MOUSE_* constants to know what button it is
func Mouse(button uint8) bool {
	return mouseState[button]
}

// Key returns the state of a keyboard button
func Key(key sdl.Scancode) bool {
	keyState := sdl.GetKeyboardState()

	if keyState[key] == 1 {
		return true
	}

	return false
}

// Clean removes every used joystick
func Clean() {
	for k := range joysticks {
		remJoystick(k)
	}
}

func addJoystick(id sdl.JoystickID) {
	if joy := sdl.JoystickOpen(id); joy != nil {
		id = joy.InstanceID()

		joysticks[id] = joy
		joystickAxises[id] = make([]int16, joy.NumAxes())
		joystickButtons[id] = make([]bool, joy.NumButtons())
		joystickHats[id] = make([]uint8, joy.NumHats())

		log.Printf("Input // Added %s as Joystick %d\n", joy.Name(), id)
	}
}

func remJoystick(id sdl.JoystickID) {
	if joy := joysticks[id]; joy != nil {
		joy.Close()

		delete(joysticks, id)
		delete(joystickAxises, id)
		delete(joystickButtons, id)

		log.Printf("Input // Removed Joystick %d\n", id)
	}
}
