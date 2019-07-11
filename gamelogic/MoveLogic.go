package gamelogic

import "math"

type MoveLogic struct {
}

func (m *MoveLogic) GetPossibleMoves(board *Board, player *Player) []*Move {
	var moves []*Move
	var fields = m.GetPiranhas(board, player)

	for _, field := range fields {
		for _, dir := range Directions() {
			dist := m.CalculateMoveDistance(board, field, dir)
			if dist > 0 {
				move := NewMove(field.X, field.Y, dir)
				if m.IsValidMove(board, player, move, dist) {
					moves = append(moves, move)
				}
			}
		}
	}

	return moves
}

func (m *MoveLogic) GetMovesToSwarm(board *Board, player *Player) []*Move {
	moves := m.GetPossibleMoves(board, player)
	swarm := m.GetSwarm(board, player)

	var res []*Move

	for _, move := range moves {
		minDistance := 1000.0
		minTargetDistance := 1000.0

		targetField := m.GetFieldInDirection(board, move, m.CalculateMoveDistance(board, board.GetField(move.X, move.Y), move.Direction))

		for _, s := range swarm {
			dist := math.Sqrt(math.Pow(float64(s.X-move.X), 2) + math.Pow(float64(s.Y-move.Y), 2))
			targetDist := math.Sqrt(math.Pow(float64(s.X-targetField.X), 2) + math.Pow(float64(s.Y-targetField.Y), 2))
			if dist < minDistance {
				minDistance = dist
			}
			if targetDist < minTargetDistance {
				minTargetDistance = targetDist
			}
		}

		if minTargetDistance < minDistance && minDistance > 0.0 {
			res = append(res, move)
		}
	}

	return res
}

func (m *MoveLogic) IsValidMove(board *Board, player *Player, move *Move, distance int) bool {

	nextField := m.GetFieldInDirection(board, move, distance)

	if nextField == nil {
		return false
	}

	fieldsInDirection := m.getFieldsInDirection(board, move, distance)
	for _, field := range fieldsInDirection {
		if field != nextField {
			if field.IsPiranha() && !field.IsPiranhaOfPlayer(player) {
				return false
			}
		}
	}

	if nextField.IsPiranhaOfPlayer(player) {
		return false
	}

	if nextField.IsObstructed() {
		return false
	}

	return true
}

func (m *MoveLogic) GetFieldInDirection(board *Board, move *Move, distance int) *Field {
	targetX := move.X
	targetY := move.Y

	switch move.Direction {
	case DirectionLeft:
		targetX -= distance
		break
	case DirectionRight:
		targetX += distance
		break
	case DirectionUp:
		targetY += distance
		break
	case DirectionDown:
		targetY -= distance
		break
	case DirectionUpRight:
		targetY += distance
		targetX += distance
		break
	case DirectionDownLeft:
		targetY -= distance
		targetX -= distance
		break
	case DirectionDownRight:
		targetY -= distance
		targetX += distance
		break
	case DirectionUpLeft:
		targetY += distance
		targetX -= distance
		break
	}

	if targetX < 0 || targetX >= board.width || targetY < 0 || targetY >= board.height {
		return nil
	}

	return board.GetField(targetX, targetY)
}

func (m *MoveLogic) getFieldsInDirection(board *Board, move *Move, distance int) []*Field {
	var fields []*Field
	for i := 0; i < distance; i++ {
		fields = append(fields, m.GetFieldInDirection(board, move, i))
	}
	return fields
}

func (m *MoveLogic) CalculateMoveDistance(board *Board, field *Field, direction Direction) int {
	switch direction {
	case DirectionLeft:
		fallthrough
	case DirectionRight:
		return m.moveDistanceHorizontal(board, field.Y)
	case DirectionUp:
		fallthrough
	case DirectionDown:
		return m.moveDistanceVertical(board, field.X)
	case DirectionUpRight:
		fallthrough
	case DirectionDownLeft:
		return m.moveDistanceDiagonalRising(board, field.X, field.Y)
	case DirectionDownRight:
		fallthrough
	case DirectionUpLeft:
		return m.moveDistanceDiagonalFalling(board, field.X, field.Y)
	}
	return -1
}

func (m *MoveLogic) moveDistanceHorizontal(board *Board, y int) int {
	count := 0
	for x := 0; x < board.width; x++ {
		if board.GetField(x, y).IsPiranha() {
			count++
		}
	}

	return count
}

func (m *MoveLogic) moveDistanceVertical(board *Board, x int) int {
	count := 0
	for y := 0; y < board.height; y++ {
		if board.GetField(x, y).IsPiranha() {
			count++
		}
	}
	return count
}

func (m *MoveLogic) moveDistanceDiagonalRising(board *Board, x int, y int) int {
	count := 0
	cX := x
	cY := y

	for cX >= 0 && cY >= 0 {
		if board.GetField(cX, cY).IsPiranha() {
			count++
		}
		cY--
		cX--
	}

	cX = x + 1
	cY = y + 1

	for cX < board.width && cY < board.height {
		if board.GetField(cX, cY).IsPiranha() {
			count++
		}
		cY++
		cX++
	}

	return count
}

func (m *MoveLogic) moveDistanceDiagonalFalling(board *Board, x int, y int) int {
	count := 0
	cX := x
	cY := y

	for cX < board.width && cY >= 0 {
		if board.GetField(cX, cY).IsPiranha() {
			count++
		}
		cY--
		cX++
	}

	cX = x - 1
	cY = y + 1

	for cX >= 0 && cY < board.height {
		if board.GetField(cX, cY).IsPiranha() {
			count++
		}
		cY++
		cX--
	}

	return count
}

func (m *MoveLogic) CalculateSwarmDistance(board *Board, player *Player) float64 {
	ownFields := m.GetPiranhas(board, player)
	sumX := 0
	sumY := 0
	for _, field := range ownFields {
		sumX += field.X
		sumY += field.Y
	}
	avgX := float64(sumX) / float64(len(ownFields))
	avgY := float64(sumY) / float64(len(ownFields))

	distance := 0.0
	for _, field := range ownFields {
		distance += math.Sqrt(math.Pow(float64(field.X)-avgX, 2) + math.Pow(float64(field.Y)-avgY, 2))
	}

	return distance
}

func (m *MoveLogic) CalculateDistanceToSwarm(board *Board, player *Player) float64 {
	swarm := m.GetSwarm(board, player)
	piranhas := m.GetPiranhas(board, player)
	distance := 0.0

	for _, p := range piranhas {
		minDistance := math.MaxFloat64
		for _, s := range swarm {
			d := math.Sqrt(math.Pow(float64(p.X)-float64(s.X), 2) + math.Pow(float64(p.Y)-float64(s.Y), 2))
			if d < minDistance {
				minDistance = d
			}
		}
		distance += minDistance
	}
	return distance
}

func (m *MoveLogic) ApplyMove(board *Board, move *Move) *Board {
	newBoard := board.Clone()

	sourceField := board.GetField(move.X, move.Y)

	targetField := m.GetFieldInDirection(board, move, m.CalculateMoveDistance(board, sourceField, move.Direction))

	newBoard.SetField(NewField(sourceField.X, sourceField.Y, FieldTypeEmpty))
	newBoard.SetField(NewField(targetField.X, targetField.Y, sourceField.T))

	return newBoard
}

func (m *MoveLogic) getSwarmHelper(board *Board, player *Player, field *Field, swarm *map[*Field]int) {
	if _, ok := (*swarm)[field]; !ok {
		(*swarm)[field] = 1
		x := field.X
		y := field.Y
		if x > 0 && board.GetField(x-1, y).IsPiranhaOfPlayer(player) {
			m.getSwarmHelper(board, player, board.GetField(x-1, y), swarm)
		}
		if x < 9 && board.GetField(x+1, y).IsPiranhaOfPlayer(player) {
			m.getSwarmHelper(board, player, board.GetField(x+1, y), swarm)
		}
		if y > 0 && board.GetField(x, y-1).IsPiranhaOfPlayer(player) {
			m.getSwarmHelper(board, player, board.GetField(x, y-1), swarm)
		}
		if y < 9 && board.GetField(x, y+1).IsPiranhaOfPlayer(player) {
			m.getSwarmHelper(board, player, board.GetField(x, y+1), swarm)
		}
		if x > 0 && y > 0 && board.GetField(x-1, y-1).IsPiranhaOfPlayer(player) {
			m.getSwarmHelper(board, player, board.GetField(x-1, y-1), swarm)
		}
		if x < 9 && y > 0 && board.GetField(x+1, y-1).IsPiranhaOfPlayer(player) {
			m.getSwarmHelper(board, player, board.GetField(x+1, y-1), swarm)
		}
		if x > 0 && y < 9 && board.GetField(x-1, y+1).IsPiranhaOfPlayer(player) {
			m.getSwarmHelper(board, player, board.GetField(x-1, y+1), swarm)
		}
		if x < 9 && y < 9 && board.GetField(x+1, y+1).IsPiranhaOfPlayer(player) {
			m.getSwarmHelper(board, player, board.GetField(x+1, y+1), swarm)
		}
	}
}

func (m *MoveLogic) GetSwarm(board *Board, player *Player) []*Field {
	fields, ok := board.swarm[player]
	if !ok {
		dict := make(map[*Field]int)
		piranhas := m.GetPiranhas(board, player)

		x := 0
		y := 0
		count := 0
		for count != len(piranhas) {
			f := board.GetField(x, y)
			_, ok := dict[f]
			if f.IsPiranhaOfPlayer(player) && !ok {
				swarm := make(map[*Field]int)
				m.getSwarmHelper(board, player, f, &swarm)
				idx := count
				for k := range swarm {
					dict[k] = idx
					count++
				}
			}
			x++
			if x == 10 {
				x = 0
				y++
			}
		}

		d := make(map[int][]*Field)
		for k, v := range dict {
			d[v] = append(d[v], k)
		}

		maxSize := 0
		for _, swarm := range d {
			if len(swarm) > maxSize {
				maxSize = len(swarm)
				fields = swarm
			}
		}
		board.swarm[player] = fields
	}
	return fields
}

func (m *MoveLogic) IsInSwarm(board *Board, player *Player, field *Field) bool {
	swarm := m.GetSwarm(board, player)
	for _, f := range swarm {
		if f.X == field.X && f.Y == field.Y {
			return true
		}
	}
	return false
}

func (m *MoveLogic) CalculateSwarmSize(board *Board, player *Player) int {
	return len(m.GetSwarm(board, player))
}

func (m *MoveLogic) GetPiranhas(board *Board, player *Player) []*Field {
	piranhas, ok := board.piranhas[player]
	if !ok {
		piranhas = nil
		for y := 0; y < board.height; y++ {
			for x := 0; x < board.width; x++ {
				if board.GetField(x, y).IsPiranhaOfPlayer(player) {
					piranhas = append(piranhas, board.GetField(x, y))
				}
			}
		}
		board.piranhas[player] = piranhas
	}

	return piranhas
}

func (m *MoveLogic) GetPiranhaCount(board *Board, player *Player) int {
	return len(m.GetPiranhas(board, player))
}

func (m *MoveLogic) HasPlayerWon(board *Board, player *Player) bool {
	return m.GetPiranhaCount(board, player) == m.CalculateSwarmSize(board, player)
}
