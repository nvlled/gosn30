package xdo

// #include <stdlib.h>
// #include <xdo.h>
// #cgo LDFLAGS: -lxdo
import "C"
import "unsafe"

const CURRENTWINDOW = 0

const (
	MbLeft = iota + 1
	MbMid
	MbRight
	MbWheelUp
	MbWheelDown
)

const (
	KeyReturn = "Return"
	KeySpace  = "space"
	KeyDelete = "Delete"
)

type Window int

type Xdo struct {
	xdo      *C.xdo_t
	Window   int
	KeyDelay int
}

func New() *Xdo {
	x := new(Xdo)
	x.xdo = C.xdo_new(nil)
	x.KeyDelay = 12000
	x.Window = CURRENTWINDOW
	return x
}

func (t *Xdo) MouseMove(x, y int) {
	C.xdo_move_mouse_relative(t.xdo, C.int(x), C.int(y))
}

func (t *Xdo) MouseDown(mouseButton int) {
	C.xdo_mouse_down(t.xdo, C.Window(t.Window), C.int(mouseButton))
}
func (t *Xdo) MouseUp(mouseButton int) {
	C.xdo_mouse_up(t.xdo, C.Window(t.Window), C.int(mouseButton))
}
func (t *Xdo) MousePress(mouseButton int, pressed bool) {
	if pressed {
		t.MouseDown(mouseButton)
	} else {
		t.MouseUp(mouseButton)
	}
}
func (t *Xdo) MouseClick(mouseButton int) {
	t.MouseDown(mouseButton)
	t.MouseUp(mouseButton)
}

func (t *Xdo) KeyPress(keyseq string) {
	str := C.CString(keyseq)
	defer C.free(unsafe.Pointer(str))
	C.xdo_send_keysequence_window(t.xdo, C.Window(t.Window), str, C.useconds_t(t.KeyDelay))
}
