package gamelogic

import (
	"fmt"
	"math"
	"sort"
)

type Controller struct {
	state         *GameState
	roomID        string
	ownPlayer     *Player
	foreignPlayer *Player
	moveLogic     *MoveLogic
	calcPerRound  int
}

func (c *Controller) RoomID() string {
	return c.roomID
}

var singleton *Controller = nil

func GetController() *Controller {
	if singleton == nil {
		singleton = &Controller{moveLogic: &MoveLogic{}}
	}
	return singleton
}

func (c *Controller) UpdateState(newstate *GameState) {
	c.state = newstate
}

func (c *Controller) JoinRoom(roomID string) {
	c.roomID = roomID
	//TODO reset state
}

func (c *Controller) SetPlayer(ownColor Color) {
	c.ownPlayer = NewPlayer(ownColor)
	c.foreignPlayer = NewPlayer(ownColor.OppositeColor())
}

func (c *Controller) readyToPlay() bool {
	return c.state != nil && c.ownPlayer != nil && c.foreignPlayer != nil
}

func (c *Controller) CalculateStaticHeuristic(board *Board, oldBoard *Board, move *Move) float64 {
	c.calcPerRound += 1
	heuristic := 0.0

	targetField := c.moveLogic.GetFieldInDirection(oldBoard, move, c.moveLogic.CalculateMoveDistance(oldBoard, oldBoard.GetField(move.X, move.Y), move.Direction))

	//heuristic += c.moveLogic.CalculateSwarmDistance(oldBoard, c.ownPlayer) - c.moveLogic.CalculateSwarmDistance(board, c.ownPlayer)
	heuristic += c.moveLogic.CalculateDistanceToSwarm(oldBoard, c.ownPlayer) - c.moveLogic.CalculateDistanceToSwarm(board, c.ownPlayer)
	heuristic += float64(c.moveLogic.CalculateSwarmSize(board, c.ownPlayer)) - float64(c.moveLogic.CalculateSwarmSize(oldBoard, c.ownPlayer))
	heuristic += math.Min(math.Abs(float64(move.X)-4.5), math.Abs(float64(move.Y)-4.5))
	heuristic -= math.Max(math.Abs(float64(targetField.X)-4.5), math.Abs(float64(targetField.Y)-4.5))

	if c.moveLogic.HasPlayerWon(board, c.ownPlayer) {
		heuristic += 1000000.0
	}

	if c.moveLogic.HasPlayerWon(board, c.foreignPlayer) {
		heuristic = -100000.0
	}
	if targetField.IsPiranhaOfPlayer(c.foreignPlayer) {
		heuristic += 5 - math.Max(math.Abs(float64(move.X)-4.5), math.Abs(float64(move.Y)-4.5))
	}

	heuristic += (float64(c.moveLogic.CalculateSwarmSize(oldBoard, c.foreignPlayer)) - float64(c.moveLogic.CalculateSwarmSize(board, c.foreignPlayer))) / 2
	heuristic += (c.moveLogic.CalculateDistanceToSwarm(board, c.foreignPlayer) - c.moveLogic.CalculateDistanceToSwarm(oldBoard, c.foreignPlayer)) / 2
	heuristic += float64(len(c.moveLogic.GetMovesToSwarm(oldBoard, c.foreignPlayer))-len(c.moveLogic.GetMovesToSwarm(board, c.foreignPlayer))) / 2

	return heuristic
}

func (c *Controller) CalculateDynamicHeuristic(board *Board, depth int) float64 {
	moveHeuristic := 0.0

	if depth != 0 && c.calcPerRound < 50000 {
		possibleMoves := c.moveLogic.GetPossibleMoves(board, c.ownPlayer)
		type ratedMove struct {
			heuristic float64
			board     *Board
			move      *Move
		}
		var moves []ratedMove

		for _, m := range possibleMoves {
			targetBoard := c.moveLogic.ApplyMove(board, m)
			staticHeuristic := c.CalculateStaticHeuristic(targetBoard, board, m)
			moves = append(moves, ratedMove{staticHeuristic, targetBoard, m})
		}

		sort.Slice(moves, func(l, r int) bool {
			return moves[l].heuristic > moves[r].heuristic
		})
		maxBreath := 60
		if len(moves) < maxBreath {
			maxBreath = len(moves)
		}

		for _, m := range moves[:maxBreath] {
			nextDepth := depth - 1
			heuristic := c.CalculateDynamicHeuristic(m.board, nextDepth)/4 + m.heuristic
			if heuristic > moveHeuristic {
				moveHeuristic = heuristic
			}
		}
	}

	return moveHeuristic
}

func (c *Controller) NextTurn() (*Move, error) {
	c.calcPerRound = 0
	if !c.readyToPlay() {
		return nil, fmt.Errorf("controller is not ready to play")
	}

	possibleMoves := c.moveLogic.GetPossibleMoves(c.state.board, c.ownPlayer)

	bestMove := possibleMoves[0]
	bestHeuristic := -1000000000.0
	for _, move := range possibleMoves {
		targetBoard := c.moveLogic.ApplyMove(c.state.board, move)

		moveHeuristic := c.CalculateStaticHeuristic(targetBoard, c.state.board, move)
		moveHeuristic += c.CalculateDynamicHeuristic(targetBoard, 1) / 4

		if moveHeuristic > bestHeuristic {
			bestHeuristic = moveHeuristic
			bestMove = move
		}
	}

	fmt.Printf("%+v\n", bestMove)
	println(bestHeuristic)
	println(c.calcPerRound)

	return bestMove, nil
}
