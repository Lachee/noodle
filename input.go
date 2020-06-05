package noodle

type Key int

var (
	Dash         = Key(189)
	Apostrophe   = Key(222)
	Semicolon    = Key(186)
	Equals       = Key(187)
	Comma        = Key(188)
	Period       = Key(190)
	Slash        = Key(191)
	Backslash    = Key(220)
	Backspace    = Key(8)
	Tab          = Key(9)
	CapsLock     = Key(20)
	Space        = Key(32)
	Enter        = Key(13)
	Escape       = Key(27)
	Insert       = Key(45)
	PrintScreen  = Key(42)
	Delete       = Key(46)
	PageUp       = Key(33)
	PageDown     = Key(34)
	Home         = Key(36)
	End          = Key(35)
	Pause        = Key(19)
	ScrollLock   = Key(145)
	ArrowLeft    = Key(37)
	ArrowRight   = Key(39)
	ArrowDown    = Key(40)
	ArrowUp      = Key(38)
	LeftBracket  = Key(219)
	LeftShift    = Key(16)
	LeftControl  = Key(17)
	LeftSuper    = Key(73)
	LeftAlt      = Key(18)
	RightBracket = Key(221)
	RightShift   = Key(16)
	RightControl = Key(17)
	RightSuper   = Key(73)
	RightAlt     = Key(18)
	Zero         = Key(48)
	One          = Key(49)
	Two          = Key(50)
	Three        = Key(51)
	Four         = Key(52)
	Five         = Key(53)
	Six          = Key(54)
	Seven        = Key(55)
	Eight        = Key(56)
	Nine         = Key(57)
	F1           = Key(112)
	F2           = Key(113)
	F3           = Key(114)
	F4           = Key(115)
	F5           = Key(116)
	F6           = Key(117)
	F7           = Key(118)
	F8           = Key(119)
	F9           = Key(120)
	F10          = Key(121)
	F11          = Key(122)
	F12          = Key(123)
	A            = Key(65)
	B            = Key(66)
	C            = Key(67)
	D            = Key(68)
	E            = Key(69)
	F            = Key(70)
	G            = Key(71)
	H            = Key(72)
	I            = Key(73)
	J            = Key(74)
	K            = Key(75)
	L            = Key(76)
	M            = Key(77)
	N            = Key(78)
	O            = Key(79)
	P            = Key(80)
	Q            = Key(81)
	R            = Key(82)
	S            = Key(83)
	T            = Key(84)
	U            = Key(85)
	V            = Key(86)
	W            = Key(87)
	X            = Key(88)
	Y            = Key(89)
	Z            = Key(90)
	NumLock      = Key(144)
	NumMultiply  = Key(106)
	NumDivide    = Key(111)
	NumAdd       = Key(107)
	NumSubtract  = Key(109)
	NumZero      = Key(96)
	NumOne       = Key(97)
	NumTwo       = Key(98)
	NumThree     = Key(99)
	NumFour      = Key(100)
	NumFive      = Key(101)
	NumSix       = Key(102)
	NumSeven     = Key(103)
	NumEight     = Key(104)
	NumNine      = Key(105)
	NumDecimal   = Key(110)
	NumEnter     = Key(13)
)

type KeyState int

const (
	KeyStateUp KeyState = iota
	KeyStatePrePressed
	KeyStatePressed
	KeyStateDown
	KeyStatePreRelease
	KeyStateRelease
)

type Input struct {
	mouseX       int
	mouseY       int
	buttonStates [5]KeyState
}

func newInput() *Input {
	i := Input{}
	i.mouseX = 0
	i.mouseY = 0
	return &i
}

/// Stores the mouse position
func (i *Input) setMousePosition(x, y int) {
	i.mouseX = x
	i.mouseY = y
}

//Sets the mouse Down State
func (i *Input) setMouseDown(button int) {
	//println("Mouse Down", button)
	i.buttonStates[button] = KeyStatePrePressed
}

//Sets the mouse up state
func (i *Input) setMouseUp(button int) {
	//println("Mouse Up", button)
	i.buttonStates[button] = KeyStatePreRelease
}

//Processes the input changes
func (i *Input) update() {
	for index, state := range i.buttonStates {

		if state == KeyStatePrePressed {
			i.buttonStates[index] = KeyStatePressed
		} else if state == KeyStatePressed {
			i.buttonStates[index] = KeyStateDown
		}

		if state == KeyStatePreRelease {
			i.buttonStates[index] = KeyStateRelease
		} else if state == KeyStateRelease {
			i.buttonStates[index] = KeyStateUp
		}
	}
}

//GetMouseX gets the current mouse position on the screen in pixels
func (i *Input) GetMouseX() int { return i.mouseX }

//GetMouseY gets the current mouse position on the screen in pixels
func (i *Input) GetMouseY() int { return i.mouseY }

//GetMousePosition gets the current mouse position
func (i *Input) GetMousePosition() Vector2 { return NewVector2(float64(i.mouseX), float64(i.mouseY)) }

// GetButton gets if the current button is pressed
func (i *Input) GetButton(button int) bool {
	if i.buttonStates[button] >= KeyStatePressed {
		return true
	}
	return false
}

// GetButtonDown gets if the current button has started being pressed
func (i *Input) GetButtonDown(button int) bool {
	if i.buttonStates[button] == KeyStatePressed {
		return true
	}
	return false
}

//GetButtonUp gets if the current button was released
func (i *Input) GetButtonUp(button int) bool {
	if i.buttonStates[button] == KeyStateRelease {
		return true
	}
	return false
}
