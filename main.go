package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/nvlled/gosn30/gamepad"
	"github.com/nvlled/gosn30/xdo"
)

const (
	ModeKeyb = iota
	ModeMouse
)

func abs16(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

func mapValueRange(x int16, srcMin, srcMax, tgtMin, tgtMax float32) int {
	n := (float32(x) - srcMin) / (srcMax - srcMin)
	return int((tgtMax-tgtMin)*n + tgtMin)
}

func handleLockFile() {
	lockFilename := ".gosn30-lock"
	lockPath := os.Getenv("HOME") + "/" + lockFilename

	if _, err := os.Stat(lockPath); !os.IsNotExist(err) && os.Getenv("NOPE") != "1" {
		println("Another instance of this program is already running. If this is not the case,")
		println("add an ENV variable NOPE=1 to run anyway.")
		os.Exit(1)
	}

	if _, err := os.Create(lockPath); err != nil {
		panic(err)
	}

	killSignal := make(chan os.Signal, 1)
	signal.Notify(killSignal, os.Interrupt)
	<-killSignal
	println("removing lockfile")
	os.Remove(lockPath)
	os.Exit(0)
}

func main() {
	for {
		mode := ModeKeyb
		gpad := gamepad.New()
		xd := xdo.New()

		go handleLockFile()
		go gpad.StartLoop()

		processKeyInput := func(event *gamepad.Event) {
			// TODO:
			//fmt.Printf("pressed: %v, type=%v, dir=%v | %v, %v\n", event.Pressed, event.InputType, event.InputValue, gamepad.InputAnalogLeft, gamepad.DirUp)
			if event.IsLeftAnalog(gamepad.DirUp) {
				xd.SetShift(event.Pressed)
			} else if event.IsRightAnalog(gamepad.DirUp) {
				xd.SetShift(event.Pressed)
			}

			if event.IsLeftAnalog(gamepad.DirRight) {
				fmt.Printf(">ctrl: %v\n", event.Pressed)
				xd.SetCtrl(event.Pressed)
			}

			if !event.Pressed {
				return
			}
			println("X")

			if gpad.IsButtonDown(gamepad.ButtonL) && gpad.IsButtonDown(gamepad.ButtonR) {
				if event.IsDpad(gamepad.DirLeft) {
					xd.KeyPress("Left")
				} else if event.IsDpad(gamepad.DirUp) {
					xd.KeyPress("Up")
				} else if event.IsDpad(gamepad.DirDown) {
					xd.KeyPress("Down")
				} else if event.IsDpad(gamepad.DirRight) {
					xd.KeyPress("Right")
				}
			} else if gpad.IsButtonDown(gamepad.ButtonL) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("h")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("l")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("i")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("n")
				} else if event.IsDpad(gamepad.DirLeft) {
					xd.KeyPress("g")
				} else if event.IsDpad(gamepad.DirUp) {
					xd.KeyPress("b")
				} else if event.IsDpad(gamepad.DirDown) {
					xd.KeyPress("w")
				} else if event.IsDpad(gamepad.DirRight) {
					xd.KeyPress("f")
				}
			} else if gpad.IsButtonDown(gamepad.ButtonR) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("y")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("p")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("u")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("m")
				} else if event.IsDpad(gamepad.DirLeft) {
					xd.KeyPress("q")
				} else if event.IsDpad(gamepad.DirUp) {
					xd.KeyPress("z")
				} else if event.IsDpad(gamepad.DirDown) {
					xd.KeyPress("v")
				} else if event.IsDpad(gamepad.DirRight) {
					xd.KeyPress("x")
				}
			} else if gpad.IsShoulderDown(gamepad.ShoulderL) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("5")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("6")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("7")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("8")
				} else if event.IsDpad(gamepad.DirLeft) {
					xd.KeyPress("1")
				} else if event.IsDpad(gamepad.DirUp) {
					xd.KeyPress("2")
				} else if event.IsDpad(gamepad.DirDown) {
					xd.KeyPress("3")
				} else if event.IsDpad(gamepad.DirRight) {
					xd.KeyPress("4")
				}
			} else if gpad.IsShoulderDown(gamepad.ShoulderR) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("0")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("9")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("k")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("j")
				} else if event.IsDpad(gamepad.DirLeft) {
					xd.KeyPress("BackSpace")
				} else if event.IsDpad(gamepad.DirUp) {
					xd.KeyPress("Delete")
				} else if event.IsDpad(gamepad.DirDown) {
					xd.KeyPress("Return")
				} else if event.IsDpad(gamepad.DirRight) {
					xd.KeyPress("space")
				}
			} else if gpad.IsLeftAnalog(gamepad.DirLeft) {
				if event.IsButton(gamepad.ButtonY) {
					xd.EnterText("ä")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.EnterText("ö")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.EnterText("å")
				} else if event.IsButton(gamepad.ButtonA) {
				}
			} else if gpad.IsRightAnalog(gamepad.DirLeft) {
				if event.IsDpad(gamepad.DirLeft) {
					xd.KeyPress("period")
				} else if event.IsDpad(gamepad.DirUp) {
					xd.KeyPress("comma")
				} else if event.IsDpad(gamepad.DirDown) {
					xd.KeyPress("colon")
				} else if event.IsDpad(gamepad.DirRight) {
				}
			} else {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("a")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("o")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("e")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("t")
				} else if event.IsDpad(gamepad.DirLeft) {
					xd.KeyPress("d")
				} else if event.IsDpad(gamepad.DirUp) {
					xd.KeyPress("c")
				} else if event.IsDpad(gamepad.DirDown) {
					xd.KeyPress("s")
				} else if event.IsDpad(gamepad.DirRight) {
					xd.KeyPress("r")
				} else if event.IsButton(gamepad.ButtonSelect) {
					mode = ModeMouse
					beeep.Notify("mouse", "", "")
				} else if event.IsButton(gamepad.ButtonStart) {
					xd.ToggleCapsLock()
					if xd.IsCapsLock() {
						beeep.Notify("uppercase", "", "")
					} else {
						beeep.Notify("lowercase", "", "")
					}
				} else if event.IsButton(gamepad.ButtonLeftStick) {
					xd.ToggleCtrl()
				} else if event.IsButton(gamepad.ButtonRightStick) {
					xd.ToggleAlt()
				}
			}
		}

		processMouseInput := func(event *gamepad.Event) {
			if event.IsButton(gamepad.ButtonA) {
				xd.MousePress(xdo.MbLeft, event.Pressed)
			} else if event.IsButton(gamepad.ButtonB) {
				xd.MousePress(xdo.MbRight, event.Pressed)
			} else if event.Pressed {
				if event.IsButton(gamepad.ButtonSelect) {
					mode = ModeKeyb
					beeep.Notify("keyboard", "", "")
				} else if gpad.IsShoulderDown(gamepad.ShoulderL) && gpad.IsRightAnalog(gamepad.DirLeft) {
					xd.KeyPress("Alt_L+Left")
				} else if gpad.IsShoulderDown(gamepad.ShoulderL) && gpad.IsRightAnalog(gamepad.DirRight) {
					xd.KeyPress("Alt_L+Right")
				} else if event.IsDpad(gamepad.DirLeft) {
					xd.KeyPress("Left")
				} else if event.IsDpad(gamepad.DirUp) {
					xd.KeyPress("Up")
				} else if event.IsDpad(gamepad.DirDown) {
					xd.KeyPress("Down")
				} else if event.IsDpad(gamepad.DirRight) {
					xd.KeyPress("Right")
				}
			}
		}

		keyChan := make(chan *gamepad.Event)
		gpad.Poll(func(event *gamepad.Event) {
			if mode == ModeMouse {
				processMouseInput(event)
			} else {
				keyChan <- event
			}
		})
		go func() {
			for {
				processKeyInput(<-keyChan)
			}
		}()
		/*
			go func() {
				buttonDown := false
				downTime := time.Now().UnixNano()
				micros := int64(1000 * 1000)
				for {
					time.Sleep(100 * time.Millisecond)
					lastEvent := gpad.LastEvent

					if mode != ModeKeyb || lastEvent == nil {
						continue
					}

					eventVal := lastEvent.InputValue
					if eventVal == gamepad.ButtonL ||
						eventVal == gamepad.ButtonR {
						continue
					}
					if lastEvent.InputType != gamepad.InputDpad && lastEvent.InputType != gamepad.InputButton {
						continue
					}

					if !lastEvent.Pressed {
						buttonDown = false
						continue
					}

					if !buttonDown && lastEvent.Pressed {
						downTime = time.Now().UnixNano()
						buttonDown = true
					} else {
						buttonDown = lastEvent.Pressed
						now := time.Now().UnixNano()
						if buttonDown {
							if (now-downTime)/micros > 1200 {
								//fmt.Printf("button down elapsed: %v\n", (now-downTime)/micros)
								keyChan <- lastEvent
							}
						}
					}

				}
			}()
		*/

		lastScroll := time.Now().UnixNano()
		for {

			if mode == ModeMouse {
				var maxSpeed float32 = 20.0
				if gpad.IsButtonDown(gamepad.ButtonR) {
					maxSpeed = 7
				}

				var dx, dy int
				if gpad.State.LeftStick.X != 0 {
					dx = mapValueRange(gpad.State.LeftStick.X, -32767, 32767, -maxSpeed, maxSpeed)
				}
				if gpad.State.LeftStick.Y != 0 {
					dy = mapValueRange(gpad.State.LeftStick.Y, -32767, 32767, -maxSpeed, maxSpeed)
				}
				if dx != 0 || dy != 0 {
					xd.MouseMove(dx, dy)
				}

				if gpad.State.RightStick.Y != 0 {
					delay := mapValueRange(abs16(gpad.State.RightStick.Y), 0, 32767, 128, 0)
					now := time.Now().UnixNano()
					millis := (now - lastScroll) / 1000000
					fmt.Printf("delay: %v, diff: %v\n", delay, millis)
					if millis >= int64(delay) {
						if gpad.State.RightStick.Y < 0 {
							xd.MouseClick(xdo.MbWheelUp)
						} else {
							xd.MouseClick(xdo.MbWheelDown)
						}
						lastScroll = now
					}
				}
			}
			time.Sleep(20 * time.Millisecond)
		}

	}
}
