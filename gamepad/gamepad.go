package gamepad

// #include <stdio.h>
// #include <linux/joystick.h>
// #include <xdo.h>
// #cgo LDFLAGS: -lxdo
import "C"
import (
	"unsafe"
)

type Vec struct {
	X int16
	Y int16
}

func (v *Vec) Normalize() Vec {
	w := Vec{}
	if w.X != 0 {
		if w.X <= -32767/2 {
			w.X = -1
		} else if w.X <= -32767/2 {
			w.X = 1
		}
	}
	if w.Y != 0 {
		if w.Y <= -32767/2 {
			w.Y = -1
		} else if w.Y <= -32767/2 {
			w.Y = 1
		}
	}
	return w
}

const (
	JsEventButton = 0x01
	JsEventAxis   = 0x02
	JsEventInit   = 0x80
)

const (
	InputDpad = iota
	InputButton
	InputShoulder
	InputAnalogLeft
	InputAnalogRight
)

const (
	ButtonB = iota
	ButtonA
	ButtonY
	ButtonX
	ButtonL
	ButtonR
	ButtonSelect
	ButtonStart
	ButtonLeftStick
	ButtonRightStick
)

const (
	DirLeft = iota
	DirRight
	DirUp
	DirDown
)

const (
	ShoulderL = iota
	ShoulderR
)

type EventHandler func(event *Event)

type GamePad struct {
	event C.struct_js_event

	FD       uintptr
	State    State
	handlers []EventHandler
}

type State struct {
	Capslock bool

	DpadFlags     uint8
	ButtonFlags   uint8
	ShoulderFlags uint8

	StartDown      bool
	SelectDown     bool
	LeftStickDown  bool
	RightStickDown bool

	LeftStick  Vec
	RightStick Vec
}

type Event struct {
	Type   uint8
	Number uint8
	Value  int16

	Pressed    bool
	InputType  int
	InputValue int
}

func (ev *Event) SetInput(inputType int, inputValue int) {
	ev.InputType = inputType
	ev.InputValue = inputValue
}

func (ev *Event) IsButton(button int) bool {
	return ev.InputType == InputButton && ev.InputValue == button
}

func (ev *Event) IsDpad(dpad int) bool {
	return ev.InputType == InputDpad && ev.InputValue == dpad
}

func New() *GamePad {
	return &GamePad{}
}

func (gpad *GamePad) Read() *Event {

	var bytes C.ssize_t
	bytes = C.read(C.int(gpad.FD), unsafe.Pointer(&gpad.event), C.sizeof_struct_js_event)
	if bytes < C.sizeof_struct_js_event {
		return nil
	}

	return &Event{
		Type:   uint8(gpad.event._type),
		Number: uint8(gpad.event.number),
		Value:  int16(gpad.event.value),
	}
}

func (gpad *GamePad) SetButtonState(button uint8, pressed bool) {
	if pressed {
		gpad.State.ButtonFlags |= 1 << button
	} else {
		gpad.State.ButtonFlags &= ^(1 << button)
	}
}

func (gpad *GamePad) SetShoulderState(shoulder uint8, pressed bool) {
	if pressed {
		gpad.State.ShoulderFlags |= 1 << shoulder
	} else {
		gpad.State.ShoulderFlags &= ^(1 << shoulder)
	}
}
func (gpad *GamePad) SetDpadState(dpad uint8, pressed bool) {
	if pressed {
		gpad.State.DpadFlags |= 1 << dpad
	} else {
		gpad.State.DpadFlags &= ^(1 << dpad)
	}
}

func (gpad *GamePad) IsLeftAnalog(dir uint8) bool {
	if dir == DirLeft {
		return gpad.State.LeftStick.X <= -16383
	} else if dir == DirRight {
		return gpad.State.LeftStick.X >= 16383
	}
	if dir == DirUp {
		return gpad.State.LeftStick.Y <= -16383
	} else if dir == DirDown {
		return gpad.State.LeftStick.Y >= 16383
	}
	return false
}

func (gpad *GamePad) IsRightAnalog(dir uint8) bool {
	if dir == DirLeft {
		return gpad.State.RightStick.X <= -16383
	} else if dir == DirRight {
		return gpad.State.RightStick.X >= 16383
	}
	if dir == DirUp {
		return gpad.State.RightStick.Y <= -16383
	} else if dir == DirDown {
		return gpad.State.RightStick.Y >= 16383
	}
	return false
}

func (gpad *GamePad) IsButtonDown(button uint8) bool {
	return gpad.State.ButtonFlags&(1<<button) != 0
}
func (gpad *GamePad) IsDpadDown(dpad uint8) bool {
	return gpad.State.DpadFlags&(1<<dpad) != 0
}
func (gpad *GamePad) IsShoulderDown(shoulder uint8) bool {
	return gpad.State.ShoulderFlags&(1<<shoulder) != 0
}
