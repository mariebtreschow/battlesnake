package main

import (
	"log"
	"math"
	"sort"
)

type errorString struct{ s string }

func (e *errorString) Error() string { return e.s }

func New(text string) error {
	return &errorString{text}
}

func info() BattlesnakeInfoResponse {
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Marie Treschow",
		Color:      "#FF0000",
		Head:       "default",
		Tail:       "default",
	}
}

func start(state GameState) {
	log.Printf("%s START:\n", state.Game.ID)
}

func end(state GameState) {
	log.Printf("%s END:\n\n", state.Game.ID)
}

func collideWithWalls(nextMove Coord, myHead Coord, boardWidth int, boardHeight int) bool {
	if nextMove.Y >= boardHeight || nextMove.Y < 0 || nextMove.X >= boardWidth || nextMove.X < 0 {
		return true
	}
	return false
}

func collideWithMyself(move string, nextMove Coord, mybody []Coord, amountOfMoves []int) bool {
	if len(amountOfMoves) > 1 {
		for _, moveCount := range amountOfMoves {
			moveAfterNext := nextCoordinate(move, nextMove, moveCount)
			for _, coordinate := range mybody {
				if (nextMove.Y == coordinate.Y && nextMove.X == coordinate.X) || (moveAfterNext.Y == coordinate.Y && moveAfterNext.X == coordinate.X) {
					return true
				}
			}
		}
	}
	for _, coordinate := range mybody {
		if nextMove.Y == coordinate.Y && nextMove.X == coordinate.X {
			return true
		}
	}
	return false
}

func nextCoordinate(move string, myHead Coord, lengthToNextMove int) Coord {
	switch move {
	case "up":
		myHead.Y = myHead.Y + lengthToNextMove
	case "down":
		myHead.Y = myHead.Y - lengthToNextMove
	case "right":
		myHead.X = myHead.X + lengthToNextMove
	case "left":
		myHead.X = myHead.X - lengthToNextMove
	}
	return myHead
}

func indexOfShortestDistance(arr map[int]float64) int {
	log.Println("Finding the shortest distance:", arr)
	min := arr[0]
	index := 0
	for i, value := range arr {
		if value < min {
			index = i
		}
	}
	return index
}

type Optimal struct {
	Move  string
	Value float64
}

func getStortedMap(optimalMoves map[string]float64) []Optimal {
	var sorted []Optimal
	for key, value := range optimalMoves {
		sorted = append(sorted, Optimal{key, value})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value
	})
	return sorted
}

func getLongestDistance(optimalMoves map[string]float64) string {
	log.Println("Comparing optimal moves:", optimalMoves)
	sortedMoves := getStortedMap(optimalMoves)
	return sortedMoves[0].Move
}

func getShortestDistance(optimalMoves map[string]float64) string {
	log.Println("Comparing optimal moves:", optimalMoves)
	sortedMoves := getStortedMap(optimalMoves)
	return sortedMoves[len(sortedMoves)-1].Move
}

func calulateDistanceBetweenCoord(firstPointX float64, firstPointY float64, secondPointX float64, secondPointY float64) float64 {
	return math.Sqrt(math.Pow((secondPointX-firstPointX), 2) + math.Pow((secondPointY-firstPointY), 2))
}

func distanceToFoodPerMove(nextMove Coord, food []Coord) float64 {
	log.Println("Calculate distance to each food coordinate for each move:", nextMove)
	distanceToEachFoodPerIndex := map[int]float64{}

	for i, coordinate := range food {
		firstPointX := float64(nextMove.X)
		firstPointY := float64(nextMove.Y)
		secondPointX := float64(coordinate.X)
		secondPointY := float64(coordinate.Y)
		distance := calulateDistanceBetweenCoord(firstPointX, firstPointY, secondPointX, secondPointY)
		distanceToEachFoodPerIndex[i] = distance
	}
	index := indexOfShortestDistance(distanceToEachFoodPerIndex)
	return distanceToEachFoodPerIndex[index]
}

func stringInSlice(find string, list []string) bool {
	log.Println("Finding item in list:", find, list)
	for _, element := range list {
		if element == find {
			return true
		}
	}
	return false
}

func shortestMoveToFood(myHead Coord, food []Coord, safeMoves []string) string {
	var firstMove string = safeMoves[0]
	optimalMoves := make(map[string]float64)
	for i := 0; i < len(safeMoves); i++ {
		optimalMoves[safeMoves[i]] = 0.0
	}

	if len(optimalMoves) > 1 {
		for _, evaluateMove := range safeMoves {
			lengthToNextMove := 1
			coordEvaluateMove := nextCoordinate(evaluateMove, myHead, lengthToNextMove)
			distance := distanceToFoodPerMove(coordEvaluateMove, food)
			optimalMoves[evaluateMove] = distance
		}
		return getShortestDistance(optimalMoves)
	} else {
		return firstMove
	}
}

func makeRange(max int) []int {
	min := 0
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func movesToCollideWithMyself(move string, myHead Coord, lengthToNextMove int, mybody []Coord) bool {
	coordinate := nextCoordinate(move, myHead, lengthToNextMove)
	moves := makeRange(lengthToNextMove)
	return collideWithMyself(move, coordinate, mybody, moves)
}

func moveWithLongestDistance(myHead Coord, safeMoves []string, mybody []Coord, boardWidth int, boardHeight int) string {
	firstMove := safeMoves[0]
	checkMovesDistanceToWall := make(map[string]float64)

	for _, evaluateMove := range safeMoves {
		nextMove := nextCoordinate(evaluateMove, myHead, 1)

		if evaluateMove == "right" {
			checkMovesDistanceToWall[evaluateMove] = float64(boardWidth - nextMove.X)
		} else if evaluateMove == "left" {
			checkMovesDistanceToWall[evaluateMove] = float64(boardWidth - (boardWidth - nextMove.X))
		} else if evaluateMove == "down" {
			checkMovesDistanceToWall[evaluateMove] = float64(boardHeight - (boardHeight - nextMove.Y))
		} else if evaluateMove == "up" {
			checkMovesDistanceToWall[evaluateMove] = float64(boardHeight - nextMove.Y)
		}
	}
	bestMove := getLongestDistance(checkMovesDistanceToWall)
	twoMoves := 2
	collidWithMyselfTwoMoves := movesToCollideWithMyself(bestMove, myHead, twoMoves, mybody)
	threeMoves := 3
	collidWithMyselfThreeMoves := movesToCollideWithMyself(bestMove, myHead, threeMoves, mybody)
	if !collidWithMyselfThreeMoves && !collidWithMyselfTwoMoves {
		return bestMove
	}
	return firstMove
}

func getPossibelMoves(myHead Coord, boardWidth int, boardHeight int, mybody []Coord, additionalMoves int) []string {
	safeMoves := []string{}
	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}
	for move, _ := range possibleMoves {
		lengthToNextMove := 1
		var newCoordinate = nextCoordinate(move, myHead, lengthToNextMove)
		moves := makeRange(additionalMoves)
		notHitMyself := collideWithMyself(move, newCoordinate, mybody, moves)
		noWalls := collideWithWalls(newCoordinate, myHead, boardWidth, boardHeight)

		if !notHitMyself && !noWalls {
			safeMoves = append(safeMoves, move)
		}
	}
	return safeMoves
}

func move(state GameState) BattlesnakeMoveResponse {
	myHead := state.You.Body[0]
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height
	mybody := state.You.Body
	health := state.You.Health

	safeMoves := getPossibelMoves(myHead, boardWidth, boardHeight, mybody, 1)

	var nextMove string

	if len(safeMoves) == 0 {
		safeMoves = getPossibelMoves(myHead, boardWidth, boardHeight, mybody, 0)
		if len(safeMoves) == 0 {
			nextMove = "up"
			log.Printf("%s MOVE %d: No safe moves detected! Moving %s\n", state.Game.ID, state.Turn, nextMove)
		} else {
			nextMove = moveWithLongestDistance(myHead, safeMoves, mybody, boardWidth, boardHeight)
		}
	} else {
		food := state.Board.Food
		nextMoveFood := shortestMoveToFood(myHead, food, safeMoves)
		nextMoveSpace := moveWithLongestDistance(myHead, safeMoves, mybody, boardWidth, boardHeight)

		if nextMoveFood == nextMoveSpace || health < 10 {
			nextMove = nextMoveFood
		} else {
			nextMove = moveWithLongestDistance(myHead, safeMoves, mybody, boardWidth, boardHeight)
		}
		log.Printf("%s MOVE %d: %s\n", state.Game.ID, state.Turn, nextMove)
	}
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}
