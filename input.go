package noodle

//Key is a Keycode representation of a keyboard character.
// use the https://keycode.info/ tool to help identify missing key codes
type Key int

const (
	//Dash -
	KeyDash         = Key(189)
	KeyApostrophe   = Key(222)
	KeySemicolon    = Key(186)
	KeyEquals       = Key(187)
	KeyComma        = Key(188)
	KeyPeriod       = Key(190)
	KeySlash        = Key(191)
	KeyBackslash    = Key(220)
	KeyBackspace    = Key(8)
	KeyTab          = Key(9)
	KeyCapsLock     = Key(20)
	KeySpace        = Key(32)
	KeyEnter        = Key(13)
	KeyEscape       = Key(27)
	KeyInsert       = Key(45)
	KeyPrintScreen  = Key(42)
	KeyDelete       = Key(46)
	KeyPageUp       = Key(33)
	KeyPageDown     = Key(34)
	KeyHome         = Key(36)
	KeyEnd          = Key(35)
	KeyPause        = Key(19)
	KeyScrollLock   = Key(145)
	KeyArrowLeft    = Key(37)
	KeyArrowRight   = Key(39)
	KeyArrowDown    = Key(40)
	KeyArrowUp      = Key(38)
	KeyLeftBracket  = Key(219)
	KeyLeftShift    = Key(16)
	KeyLeftControl  = Key(17)
	KeyLeftSuper    = Key(73)
	KeyLeftAlt      = Key(18)
	KeyRightBracket = Key(221)
	KeyRightShift   = Key(16)
	KeyRightControl = Key(17)
	KeyRightSuper   = Key(73)
	KeyRightAlt     = Key(18)
	KeyZero         = Key(48)
	KeyOne          = Key(49)
	KeyTwo          = Key(50)
	KeyThree        = Key(51)
	KeyFour         = Key(52)
	KeyFive         = Key(53)
	KeySix          = Key(54)
	KeySeven        = Key(55)
	KeyEight        = Key(56)
	KeyNine         = Key(57)
	KeyF1           = Key(112)
	KeyF2           = Key(113)
	KeyF3           = Key(114)
	KeyF4           = Key(115)
	KeyF5           = Key(116)
	KeyF6           = Key(117)
	KeyF7           = Key(118)
	KeyF8           = Key(119)
	KeyF9           = Key(120)
	KeyF10          = Key(121)
	KeyF11          = Key(122)
	KeyF12          = Key(123)
	KeyA            = Key(65)
	KeyB            = Key(66)
	KeyC            = Key(67)
	KeyD            = Key(68)
	KeyE            = Key(69)
	KeyF            = Key(70)
	KeyG            = Key(71)
	KeyH            = Key(72)
	KeyI            = Key(73)
	KeyJ            = Key(74)
	KeyK            = Key(75)
	KeyL            = Key(76)
	KeyM            = Key(77)
	KeyN            = Key(78)
	KeyO            = Key(79)
	KeyP            = Key(80)
	KeyQ            = Key(81)
	KeyR            = Key(82)
	KeyS            = Key(83)
	KeyT            = Key(84)
	KeyU            = Key(85)
	KeyV            = Key(86)
	KeyW            = Key(87)
	KeyX            = Key(88)
	KeyY            = Key(89)
	KeyZ            = Key(90)
	KeyNumLock      = Key(144)
	KeyNumMultiply  = Key(106)
	KeyNumDivide    = Key(111)
	KeyNumAdd       = Key(107)
	KeyNumSubtract  = Key(109)
	KeyNumZero      = Key(96)
	KeyNumOne       = Key(97)
	KeyNumTwo       = Key(98)
	KeyNumThree     = Key(99)
	KeyNumFour      = Key(100)
	KeyNumFive      = Key(101)
	KeyNumSix       = Key(102)
	KeyNumSeven     = Key(103)
	KeyNumEight     = Key(104)
	KeyNumNine      = Key(105)
	KeyNumDecimal   = Key(110)
	KeyNumEnter     = Key(13)
)

type keyState int

const (
	keyStateUp keyState = iota
	keyStatePrePressed
	keyStatePressed
	keyStateDown
	keyStatePreRelease
	keyStateRelease
)

//InputHandler handles the different states of the input
type InputHandler struct {
	mouseDelta         Vector2 //mouseDelta is the movement since last frame
	mouseDeltaPrevious Vector2 //mouseDeltaPrevious is the delta before the frame
	mouseX             int     //mouseX position
	mouseY             int     //mouseY position

	mouseScrollDeltaPrevious float32 //mouseScrollDeltaPrevious is the scroll before the frame
	mouseScrollDelta         float32 //mouseScrollDelta is the scroll after the frame

	buttonStates [5]keyState
	keyStates    map[Key]keyState

	normalizeCoordinates bool //normalizeCoordinates will map mouse position within the bounding box. Values will still be pixels.
}

func newInput() *InputHandler {
	i := InputHandler{}
	i.mouseX = 0
	i.mouseY = 0
	i.keyStates = make(map[Key]keyState, 6)
	return &i
}

/// Stores the mouse position
func (i *InputHandler) setMousePosition(x, y int) {
	i.mouseDeltaPrevious = Vector2{float32(x - i.mouseX), float32(y - i.mouseY)}
	i.mouseX = x
	i.mouseY = y
}

// setMouseScroll sets the mouse delta scroll
func (i *InputHandler) setMouseScroll(delta float32) {
	i.mouseScrollDeltaPrevious = delta
}

//Sets the mouse Down State
func (i *InputHandler) setMouseDown(button int) {
	i.buttonStates[button] = keyStatePrePressed
}

//Sets the mouse up state
func (i *InputHandler) setMouseUp(button int) {
	i.buttonStates[button] = keyStatePreRelease
}

func (i *InputHandler) setKeyDown(key int) {
	i.keyStates[Key(key)] = keyStatePrePressed
}
func (i *InputHandler) setKeyUp(key int) {
	i.keyStates[Key(key)] = keyStatePreRelease
}

//Processes the input changes
func (i *InputHandler) update() {
	//Update the scroll
	i.mouseScrollDelta = i.mouseScrollDeltaPrevious
	i.mouseScrollDeltaPrevious = 0

	i.mouseDelta = i.mouseDeltaPrevious
	i.mouseDeltaPrevious = Vector2{0, 0}

	//Handle mouse inputs
	for index, state := range i.buttonStates {

		if state == keyStatePrePressed {
			i.buttonStates[index] = keyStatePressed
		} else if state == keyStatePressed {
			i.buttonStates[index] = keyStateDown
		}

		if state == keyStatePreRelease {
			i.buttonStates[index] = keyStateRelease
		} else if state == keyStateRelease {
			i.buttonStates[index] = keyStateUp
		}
	}

	//Handle key inputs. These are a bit tricker
	for index, state := range i.keyStates {

		//Pressed
		if state == keyStatePrePressed {
			i.keyStates[index] = keyStatePressed
		} else if state == keyStatePressed {
			i.keyStates[index] = keyStateDown
		}

		//Released
		if state == keyStatePreRelease {
			i.keyStates[index] = keyStateRelease
		} else if state == keyStateRelease {
			i.keyStates[index] = keyStateUp

			//Key is up, so lets remove the index so we dont loop it again
			delete(i.keyStates, index)
		}
	}
}

//GetMouseX gets the current mouse position on the screen in pixels
func (i *InputHandler) GetMouseX() int { return i.mouseX }

//GetMouseY gets the current mouse position on the screen in pixels
func (i *InputHandler) GetMouseY() int { return i.mouseY }

//GetMousePosition gets the current mouse position
func (i *InputHandler) GetMousePosition() Vector2 {
	return NewVector2(float32(i.mouseX), float32(i.mouseY))
}

//GetMouseDelta gets the mouse movement
func (i *InputHandler) GetMouseDelta() Vector2 {
	return i.mouseDelta
}

//GetMouseScroll gets the current mouse scroll delta
func (i *InputHandler) GetMouseScroll() float32 { return i.mouseScrollDelta }

// GetButton gets if the current button is pressed
func (i *InputHandler) GetButton(button int) bool {
	if i.buttonStates[button] >= keyStatePressed {
		return true
	}
	return false
}

// GetButtonDown gets if the current button has started being pressed
func (i *InputHandler) GetButtonDown(button int) bool {
	if i.buttonStates[button] == keyStatePressed {
		return true
	}
	return false
}

//GetButtonUp gets if the current button was released
func (i *InputHandler) GetButtonUp(button int) bool {
	if i.buttonStates[button] == keyStateRelease {
		return true
	}
	return false
}

// GetKey gets if the current key is pressed
func (i *InputHandler) GetKey(key Key) bool {
	if i.keyStates[key] >= keyStatePressed {
		return true
	}
	return false
}

// GetKeyDown gets if the current key has started being pressed
func (i *InputHandler) GetKeyDown(key Key) bool {
	if i.keyStates[key] == keyStatePressed {
		return true
	}
	return false
}

//GetKeyUp gets if the current key was released
func (i *InputHandler) GetKeyUp(key Key) bool {
	if i.keyStates[key] == keyStateRelease {
		return true
	}
	return false
}

//GetAxis gets a normalized 1D axis
func (i *InputHandler) GetAxis(negativeKey, positiveKey Key) float32 {
	sum := float32(0)
	if i.GetKey(negativeKey) {
		sum -= 1.0
	}
	if i.GetKey(positiveKey) {
		sum += 1.0
	}
	return sum
}

//GetAxis2D gets a normalized 2D axis based on the keys being pressed
func (i *InputHandler) GetAxis2D(negHorizontalKey, posHorizontalKey, negVerticalKey, posVerticalKey Key) Vector2 {
	axis := Vector2{i.GetAxis(negHorizontalKey, posHorizontalKey), i.GetAxis(negVerticalKey, posVerticalKey)}
	return axis.Normalize()
}

//Cursor updates the canvas cursor style. Set to 'none' if drawing custom cursors.
func (i *InputHandler) Cursor(cursor string) {
	canvas.Get("style").Set("cursor", cursor)
}
