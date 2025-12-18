package tui

func cursorUp(c, top int) int {
	if c > 0 {
		c--
	} else {
		c = top
	}

	return c
}

func cursorDown(c, top int) int {
	if c < top {
		c++
	} else {
		c = 0
	}
	return c
}
