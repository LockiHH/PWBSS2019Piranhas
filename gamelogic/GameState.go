package gamelogic

type Player struct {
	color Color
}

func NewPlayer(color Color) *Player {
	return &Player{color: color}
}

type Color int

const (
	ColorBlue Color = 0
	ColorRed Color = 1
)

func (c Color) OppositeColor() Color {
	if c == ColorBlue {
		return ColorRed
	} else {
		return ColorBlue
	}
}

type GameState struct {
	board* Board
}

func NewGameState(board *Board) *GameState {
	return &GameState{board: board}
}

type Board struct {
	fields [][]*Field
	width int
	height int
	swarm map[*Player][]*Field
	piranhas map[*Player][]*Field
}

func NewBoard(fields [][]*Field, width int, height int) *Board {
	return &Board{fields: fields, width: width, height: height, swarm:make(map[*Player][]*Field), piranhas:make(map[*Player][]*Field)}
}

func (b *Board) GetField(x int, y int) *Field {
	return b.fields[y][x]
}

func (b *Board) SetField(field *Field) {
	b.fields[field.Y][field.X] = field
}

func (b *Board) Clone() *Board {
	newFields := make([][]*Field, len(b.fields))
	for i := range b.fields {
		newFields[i] = make([]*Field, len(b.fields[i]))
		copy(newFields[i], b.fields[i])
	}
	return &Board{fields: newFields, width:b.width, height:b.height, swarm:make(map[*Player][]*Field), piranhas:make(map[*Player][]*Field)}
}

type Field struct {
	X int
	Y int
	T FieldType
}

func NewField(x int, y int, t FieldType) *Field {
	return &Field{X: x, Y: y, T: t}
}


func(f *Field) IsPiranha() bool {
	return f.T == FieldTypeRed || f.T == FieldTypeBlue
}

func(f *Field) IsPiranhaOfPlayer(player *Player) bool {
	return (f.T == FieldTypeRed && player.color == ColorRed) || (f.T == FieldTypeBlue && player.color == ColorBlue)
}

func(f *Field) IsObstructed() bool {
	return f.T == FieldTypeObstructed
}

type FieldType int

const (
	FieldTypeEmpty FieldType = 0
	FieldTypeObstructed FieldType = 1
	FieldTypeBlue FieldType = 2
	FieldTypeRed FieldType = 3
)

type Direction int

func (d Direction) String() string {
	switch d {
		case DirectionUp:
			return "UP"
		case DirectionUpRight:
			return "UP_RIGHT"
		case DirectionRight:
			return "RIGHT"
		case DirectionDownRight:
			return "DOWN_RIGHT"
		case DirectionDown:
			return "DOWN"
		case DirectionDownLeft:
			return "DOWN_LEFT"
		case DirectionLeft:
			return "LEFT"
		case DirectionUpLeft:
			return "UP_LEFT"
	}
	return ""
}

const (
	DirectionUp Direction = 0
	DirectionUpRight Direction = 1
	DirectionRight Direction = 2
	DirectionDownRight Direction = 3
	DirectionDown Direction = 4
	DirectionDownLeft Direction = 5
	DirectionLeft Direction = 6
	DirectionUpLeft Direction = 7
)

func Directions() []Direction {
	return []Direction {DirectionUp,
						DirectionUpRight,
						DirectionRight,
						DirectionDownRight,
						DirectionDown,
						DirectionDownLeft,
						DirectionLeft,
						DirectionUpLeft}
}

type Move struct {
	X int
	Y int
	Direction Direction
}

func NewMove(x int, y int, direction Direction) *Move {
	return &Move{X: x, Y: y, Direction: direction}
}


