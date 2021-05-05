package govim

type SimpleCursor struct {
	x, y int
}

func NewSimpleCursor() *SimpleCursor {
	return &SimpleCursor{}
}

func NewSimpleCursorAt(x, y int) *SimpleCursor {
	return &SimpleCursor{x, y}
}

func (c *SimpleCursor) XPos() int {
	return c.x
}

func (c *SimpleCursor) YPos() int {
	return c.y
}

func (c *SimpleCursor) MoveTo(x, y int) {
	c.x = x
	c.y = y
}

func (c *SimpleCursor) MoveLeft() {
	c.x--
}

func (c *SimpleCursor) MoveRight() {
	c.x++
}

func (c *SimpleCursor) MoveUp() {
	c.y--
}

func (c *SimpleCursor) MoveDown() {
	c.y++
}
