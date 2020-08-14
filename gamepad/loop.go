package gamepad

import (
	"fmt"
	"os"
	"time"
)

func (gpad *GamePad) Poll(fn EventHandler) {
	gpad.handlers = append(gpad.handlers, fn)
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

		go func() {
			for {
				ev := <-c
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
			ev.Pressed = ev.Value != 0

			if ev.Type == JsEventButton {
				if ev.Number >= 0 && ev.Number <= 9 {
					gpad.SetButtonState(ev.Number, ev.Value != 0)
					ev.SetInput(InputButton, int(ev.Number))
				}
			} else if ev.Type == JsEventAxis {
				if ev.Number == 0 {
					gpad.State.LeftStick.X = ev.Value
					ev.InputType = InputAnalogLeft
					if ev.Value <= -16383 {
						ev.InputValue = DirLeft
					} else if ev.Value >= 16383 {
						ev.InputValue = DirRight
					}
				} else if ev.Number == 1 {
					gpad.State.LeftStick.Y = ev.Value
					ev.InputType = InputAnalogLeft
					if ev.Value <= -16383 {
						ev.InputValue = DirUp
					} else if ev.Value >= 16383 {
						ev.InputValue = DirDown
					}
				} else if ev.Number == 3 {
					gpad.State.RightStick.X = ev.Value
					ev.InputType = InputAnalogRight
					if ev.Value <= -16383 {
						ev.InputValue = DirLeft
					} else if ev.Value >= 16383 {
						ev.InputValue = DirRight
					}
				} else if ev.Number == 4 {
					gpad.State.RightStick.Y = ev.Value
					ev.InputType = InputAnalogRight
					if ev.Value <= -16383 {
						ev.InputValue = DirUp
					} else if ev.Value >= 16383 {
						ev.InputValue = DirDown
					}
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
							gpad.SetDpadState(DirRight, false)
						} else if gpad.IsButtonDown(DirLeft) {
							ev.SetInput(InputDpad, DirLeft)
							gpad.SetDpadState(DirLeft, false)
						}
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
							gpad.SetDpadState(DirUp, false)
						} else if gpad.IsButtonDown(DirDown) {
							ev.SetInput(InputDpad, DirDown)
							gpad.SetDpadState(DirDown, false)
						}
					}
				}
			}

			c <- ev
		}

		close(c)
		file.Close()
		println("Gamepad disconnected!")
	}
}
