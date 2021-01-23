package gamepad

import (
	"fmt"
	"os"
	"time"
)

func GetAnalogDirection(horizontal bool, val int16) uint8 {
	if horizontal {
		return GetAnalogXDirection(val)
	}
	return GetAnalogYDirection(val)
}
func GetAnalogXDirection(val int16) uint8 {
	if val <= -16383 {
		return DirLeft
	} else if val >= 16383 {
		return DirRight
	}
	return 0
}
func GetAnalogYDirection(val int16) uint8 {
	if val <= -16383 {
		return DirUp
	} else if val >= 16383 {
		return DirDown
	}
	return 0
}

func SetAnalogValue(ev *Event, currentVal int16, minVal, maxVal int) {
	if ev.Value <= -16383 {
		ev.InputValue = minVal
	} else if ev.Value >= 16383 {
		ev.InputValue = maxVal
	} else {
		if currentVal <= 0 {
			ev.InputValue = minVal
		} else if currentVal >= 0 {
			ev.InputValue = maxVal
		}
	}
}

func (gpad *GamePad) Poll(fn EventHandler) {
	gpad.handlers = append(gpad.handlers, fn)
}

func (gpad *GamePad) SendEvent(ev *Event) {
	if gpad.eventChannel != nil {
		gpad.eventChannel <- ev
	}
}

func (gpad *GamePad) StartLoop() {
	for {
		devicePath := "/dev/input/js0"
		file, err := os.Open(devicePath)
		if err != nil {
			fmt.Printf("No gamepad found (%v)\n", time.Now().Unix()%1000)
			time.Sleep(2 * time.Second)
			continue
		}
		gpad.FD = file.Fd()

		c := make(chan *Event, 1)
		gpad.eventChannel = c

		go func() {
			for {
				ev := <-c
				if ev == nil {
					break
				}
				for _, fn := range gpad.handlers {
					fn(ev)
				}
				fmt.Printf("type=%v, number=%v, value=%v\n", ev.Type, ev.Number, ev.Value)
			}
		}()

		for {
			var ev *Event
			if ev = gpad.Read(); ev == nil {
				break
			}
			emit := true
			ev.gpad = gpad
			ev.Pressed = ev.Value != 0

			if ev.Type == JsEventButton {
				if ev.Number >= 0 && ev.Number <= 9 {
					gpad.SetButtonState(ev.Number, ev.Value != 0)
					ev.SetInput(InputButton, int(ev.Number))
				}
			} else if ev.Type == JsEventAxis {
				if ev.Number >= 0 && ev.Number <= 4 {
					left := ev.Number <= 1
					horizontal := ev.Number == 0 || ev.Number == 3
					dir := GetAnalogDirection(horizontal, ev.Value)
					//dir := GetAnalogYDirection(ev.Value)
					prevDir := gpad.GetAnalogDirection(left, horizontal)
					if left {
						ev.InputType = InputAnalogLeft
					} else {
						ev.InputType = InputAnalogRight
					}
					ev.Pressed = dir != 0

					inputValue := dir
					if inputValue == 0 {
						inputValue = prevDir
					}
					ev.InputValue = int(inputValue)

					emit = dir != prevDir
					gpad.SetAnalogState(left, horizontal, ev.Value)
					//gpad.State.LeftStick.Y = ev.Value
					/*
						} else if ev.Number == 0 { // leftanalogX
							dir := GetAnalogXDirection(ev.Value)
							prevDir := GetAnalogXDirection(gpad.State.LeftStick.X)
							ev.InputType = InputAnalogLeft
							ev.Pressed = dir != 0
							ev.InputValue = int(dir)

							emit = dir != prevDir
							gpad.State.LeftStick.X = ev.Value
							// TODO: do not emit if same direction
							//if GetAnalogXDirection(ev.Value) == GetAnalogXDirection(gpad.state.LeftStick.X) && ev.Pressed == gpad.pre
							//SetAnalogValue(ev, gpad.State.LeftStick.X, DirLeft, DirRight)
						} else if ev.Number == 1 { // leftanalogY
							dir := GetAnalogYDirection(ev.Value)
							prevDir := GetAnalogYDirection(gpad.State.LeftStick.Y)
							ev.InputType = InputAnalogLeft
							ev.Pressed = dir != 0

							inputValue := dir
							if inputValue == 0 {
								inputValue = prevDir
							}
							ev.InputValue = int(inputValue)

							emit = dir != prevDir
							gpad.State.LeftStick.Y = ev.Value
							//fmt.Printf(">ev.Value=%v, emit=%v, dir=%v, gpad.IsLeftAnalog=%v, pressed: %v\n", ev.Value, emit, inputValue, gpad.IsLeftAnalog(dir), ev.Pressed)
							//SetAnalogValue(ev, gpad.State.LeftStick.Y, DirUp, DirDown)
							//ev.InputType = InputAnalogLeft
							//ev.Pressed = ev.Value != 0
							//gpad.State.LeftStick.Y = ev.Value
						} else if ev.Number == 3 { // rightanalogX
							//SetAnalogValue(ev, gpad.State.RightStick.X, DirLeft, DirRight)
							//ev.InputType = InputAnalogLeft
							//ev.Pressed = ev.Value != 0
							//gpad.State.RightStick.X = ev.Value
						} else if ev.Number == 4 { // rightanalogY
							//SetAnalogValue(ev, gpad.State.RightStick.Y, DirUp, DirDown)
							//ev.InputType = InputAnalogLeft
							//ev.Pressed = ev.Value != 0
							//gpad.State.RightStick.Y = ev.Value
					*/
				} else if ev.Number == 2 {
					gpad.SetShoulderState(ShoulderL, ev.Value == 32767)
					ev.SetInput(InputShoulder, ShoulderL)
					ev.Pressed = ev.Value == 32767
				} else if ev.Number == 5 {
					gpad.SetShoulderState(ShoulderR, ev.Value == 32767)
					ev.SetInput(InputShoulder, ShoulderR)
					ev.Pressed = ev.Value == 32767
				} else if ev.Number == 6 {
					if ev.Value == -32767 {
						gpad.SetDpadState(DirLeft, true)
						gpad.SetDpadState(DirRight, false)
						ev.SetInput(InputDpad, DirLeft)
					} else if ev.Value == 32767 {
						gpad.SetDpadState(DirRight, true)
						gpad.SetDpadState(DirLeft, false)
						ev.SetInput(InputDpad, DirRight)
					} else {
						if gpad.IsButtonDown(DirRight) {
							ev.SetInput(InputDpad, DirRight)
						} else if gpad.IsButtonDown(DirLeft) {
							ev.SetInput(InputDpad, DirLeft)
						}
						gpad.SetDpadState(DirRight, false)
						gpad.SetDpadState(DirLeft, false)
					}
				} else if ev.Number == 7 {
					if ev.Value == -32767 {
						gpad.SetDpadState(DirUp, true)
						gpad.SetDpadState(DirDown, false)
						ev.SetInput(InputDpad, DirUp)
					} else if ev.Value == 32767 {
						gpad.SetDpadState(DirDown, true)
						gpad.SetDpadState(DirUp, false)
						ev.SetInput(InputDpad, DirDown)
					} else {
						if gpad.IsButtonDown(DirUp) {
							ev.SetInput(InputDpad, DirUp)
						} else if gpad.IsButtonDown(DirDown) {
							ev.SetInput(InputDpad, DirDown)
						}
						gpad.SetDpadState(DirUp, false)
						gpad.SetDpadState(DirDown, false)
					}
				}
			}

			if emit && ev != nil {
				gpad.LastEvent = ev
				c <- ev
			}
		}

		close(c)
		file.Close()
		println("Gamepad disconnected!")
	}
}
