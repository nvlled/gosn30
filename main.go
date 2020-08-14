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
			if !event.Pressed {
				return
			}
			if gpad.IsButtonDown(gamepad.ButtonL) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("d")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("c")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("s")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("r")
				}
			} else if gpad.IsButtonDown(gamepad.ButtonR) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("h")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("l")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("i")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("n")
				} else if event.IsDpad(gamepad.DirLeft) {
					xd.KeyPress("BackSpace")
				} else if event.IsDpad(gamepad.DirRight) {
					xd.KeyPress("space")
				} else if event.IsDpad(gamepad.DirUp) {
					// capslock/shift
					// xd.ToggleUpperCase()
				} else if event.IsDpad(gamepad.DirDown) {
					xd.KeyPress("Return")
				}
			} else if gpad.IsShoulderDown(gamepad.ShoulderL) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("g")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("b")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("w")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("f")
				}
			} else if gpad.IsShoulderDown(gamepad.ShoulderR) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("y")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("p")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("u")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("m")
				}
			} else if gpad.IsLeftAnalog(gamepad.DirUp) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("1")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("2")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("3")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("4")
				}
			} else if gpad.IsLeftAnalog(gamepad.DirDown) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("5")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("6")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("7")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("8")
				}
			} else if gpad.IsLeftAnalog(gamepad.DirLeft) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("q")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("z")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("v")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("x")
				}
			} else if gpad.IsLeftAnalog(gamepad.DirRight) {
				if event.IsButton(gamepad.ButtonY) {
					xd.KeyPress("9")
				} else if event.IsButton(gamepad.ButtonX) {
					xd.KeyPress("0")
				} else if event.IsButton(gamepad.ButtonB) {
					xd.KeyPress("k")
				} else if event.IsButton(gamepad.ButtonA) {
					xd.KeyPress("j")
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
					xd.KeyPress("Left")
				} else if event.IsDpad(gamepad.DirRight) {
					xd.KeyPress("Right")
				} else if event.IsDpad(gamepad.DirUp) {
					xd.KeyPress("Up")
				} else if event.IsDpad(gamepad.DirDown) {
					xd.KeyPress("Down")
				} else if event.IsButton(gamepad.ButtonSelect) {
					mode = ModeMouse
					beeep.Notify("mouse", "", "")
				} else if event.IsButton(gamepad.ButtonStart) {
					xd.KeyPress("Escape")
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
				}
			}
		}

		gpad.Poll(func(event *gamepad.Event) {
			if mode == ModeMouse {
				processMouseInput(event)
			} else {
				processKeyInput(event)
			}
		})

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
