package main

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testNeckHit(move string, snake []Coord) bool {
	me := Battlesnake{
		Head: snake[0],
		Body: snake,
	}
	amount := 1
	nextMove := nextCoordinate(move, me.Head, amount)
	moves := makeRange(amount)
	return collideWithMyself(move, nextMove, me.Body, moves)
}

func testDistanceToFoodPerMove(move string, snake []Coord, food []Coord) float64 {
	lengthToNextCoord := 1
	nextMove := nextCoordinate(move, snake[0], lengthToNextCoord)
	return distanceToFoodPerMove(nextMove, food)
}

func testWallHit(move string, snake []Coord, boardWith int, boardHeight int) bool {
	me := Battlesnake{
		Head: snake[0],
		Body: snake,
	}
	amount := 1
	nextMove := nextCoordinate(move, me.Head, amount)
	return collideWithWalls(nextMove, me.Head, boardWith, boardHeight)
}

func validatePossibleNeckMoves(snake []Coord, possibleMoves []string, t *testing.T) {
	for _, move := range possibleMoves {
		collidedWithNeck := testNeckHit(move, snake)
		if collidedWithNeck {
			t.Error("This should be a possible move:", move)
		}
	}
}

func TestNextCoordLeft(t *testing.T) {
	snake := []Coord{{X: 2, Y: 1}}
	me := Battlesnake{
		Head: snake[0],
		Body: snake,
	}
	move := "left"
	lengthToNextMove := 1
	coordinate := nextCoordinate(move, me.Head, lengthToNextMove)
	assert.Equal(t, coordinate.X, me.Head.X-1)
	assert.Equal(t, coordinate.Y, me.Head.Y)
}

func TestNextCoordUp(t *testing.T) {
	snake := []Coord{{X: 2, Y: 1}}
	me := Battlesnake{
		Head: snake[0],
		Body: snake,
	}
	move := "up"
	lengthToNextMove := 1
	coordinate := nextCoordinate(move, me.Head, lengthToNextMove)
	assert.Equal(t, coordinate.X, me.Head.X)
	assert.Equal(t, coordinate.Y, me.Head.Y+1)
}

func TestNextCoordDown(t *testing.T) {
	snake := []Coord{{X: 2, Y: 1}}
	me := Battlesnake{
		Head: snake[0],
		Body: snake,
	}
	move := "down"
	lengthToNextMove := 1
	coordinate := nextCoordinate(move, me.Head, lengthToNextMove)
	assert.Equal(t, coordinate.X, me.Head.X)
	assert.Equal(t, coordinate.Y, me.Head.Y-1)
}

func TestNextCoordRight(t *testing.T) {
	snake := []Coord{{X: 2, Y: 1}}
	me := Battlesnake{
		Head: snake[0],
		Body: snake,
	}
	move := "right"
	lengthToNextMove := 1
	coordinate := nextCoordinate(move, me.Head, lengthToNextMove)
	assert.Equal(t, coordinate.X, me.Head.X+1)
	assert.Equal(t, coordinate.Y, me.Head.Y)
}

// AVOID HITTING ITSELF IN THE NECK
func TestNeckAvoidanceLeftMove(t *testing.T) {
	// Length 3, facing right
	snake := []Coord{{X: 2, Y: 1}, {X: 1, Y: 1}, {X: 0, Y: 1}}
	impossibleMove := "left"
	possibleMoves := []string{"right", "up", "down"}
	validatePossibleNeckMoves(snake, possibleMoves, t)

	collidedWithNeck := testNeckHit(impossibleMove, snake)
	if !collidedWithNeck {
		t.Errorf("Collided with its neck %v", impossibleMove)
	}
}

func TestNeckAvoidanceRightMove(t *testing.T) {
	// Length 3, facing left
	snake := []Coord{{X: 1, Y: 1}, {X: 2, Y: 1}, {X: 3, Y: 1}}
	impossibleMove := "right"
	possibleMoves := []string{"up", "left", "down"}
	validatePossibleNeckMoves(snake, possibleMoves, t)

	collidedWithNeck := testNeckHit(impossibleMove, snake)
	if !collidedWithNeck {
		t.Errorf("Collided with its neck %v", impossibleMove)
	}
}

func TestNeckAvoidanceUpMove(t *testing.T) {
	// Length 3, facing down
	snake := []Coord{{X: 1, Y: 1}, {X: 1, Y: 2}, {X: 1, Y: 3}}
	impossibleMove := "up"
	possibleMoves := []string{"down", "left", "right"}
	validatePossibleNeckMoves(snake, possibleMoves, t)

	collidedWithNeck := testNeckHit(impossibleMove, snake)
	if !collidedWithNeck {
		t.Errorf("Collided with its neck %v", impossibleMove)
	}
}

func TestNeckAvoidanceDownMove(t *testing.T) {
	// Length 3, facing up
	snake := []Coord{{X: 2, Y: 3}, {X: 2, Y: 2}, {X: 2, Y: 1}}
	impossibleMove := "down"
	possibleMoves := []string{"up", "left", "right"}
	validatePossibleNeckMoves(snake, possibleMoves, t)

	collidedWithNeck := testNeckHit(impossibleMove, snake)
	if !collidedWithNeck {
		t.Errorf("Collided with its neck %v", impossibleMove)
	}
}

// AVOID HITTING WALLS
func TestDoNotHitWallsDown(t *testing.T) {
	// Head by the bottom left of the wall
	snake := []Coord{{X: 0, Y: 0}, {X: 0, Y: 1}, {X: 0, Y: 2}}
	impossibleMoves := []string{"down", "left"}
	boardWith, boardHeight := 3, 3

	for _, move := range impossibleMoves {
		collidingWithWall := testWallHit(move, snake, boardWith, boardHeight)
		if !collidingWithWall {
			t.Errorf("Snake moved into the walls, %s", move)
		}
	}
}

func TestDoNotHitWallsUpOrRight(t *testing.T) {
	// Head by the top right of the wall
	snake := []Coord{{X: 3, Y: 3}, {X: 3, Y: 2}, {X: 3, Y: 1}}
	impossibleMoves := []string{"top", "right"}
	boardWith, boardHeight := 3, 3

	for _, move := range impossibleMoves {
		collidingWithWall := testWallHit(move, snake, boardWith, boardHeight)
		if !collidingWithWall {
			t.Errorf("Snake moved into the walls, %s", move)
		}
	}
}

func TestCalculateDistanceToFood(t *testing.T) {
	food := []Coord{{X: 3, Y: 5}, {X: 3, Y: 7}}
	snake := []Coord{{X: 3, Y: 3}, {X: 3, Y: 2}, {X: 3, Y: 1}}
	// func distanceToFoodPerMove(nextMove Coord, food []Coord) map[int]float64{
	moveUp := "up"
	nextMoveUp := nextCoordinate(moveUp, snake[0], 1)

	moveDown := "down"
	nextMoveDown := nextCoordinate(moveDown, snake[0], 1)

	distanceUp := distanceToFoodPerMove(nextMoveUp, food)
	distanceDown := distanceToFoodPerMove(nextMoveDown, food)

	assert.Greater(t, distanceDown, distanceUp)
}

func TestBestMoveToFood(t *testing.T) {
	food := []Coord{{X: 3, Y: 5}, {X: 3, Y: 7}}
	snake := []Coord{{X: 3, Y: 3}, {X: 3, Y: 2}, {X: 3, Y: 1}}
	safeMoves := []string{"right", "up"}
	moveUp := "up"

	bestMove := shortestMoveToFood(snake[0], food, safeMoves)
	assert.Equal(t, bestMove, moveUp)
}

func TestFindShortestMove(t *testing.T) {
	optimalMoves := make(map[string]float64)
	optimalMoves["left"] = 1.77
	optimalMoves["right"] = 1.3

	bestMove := "right"

	shortestMove := getShortestDistance(optimalMoves)
	assert.Equal(t, shortestMove, bestMove)
}

func TestFindShortestMoveMoreOptions(t *testing.T) {
	optimalMoves := make(map[string]float64)
	optimalMoves["left"] = 1.77
	optimalMoves["right"] = 1.3
	optimalMoves["down"] = 4.3
	optimalMoves["up"] = 5.3

	bestMove := "right"

	shortestMove := getShortestDistance(optimalMoves)
	assert.Equal(t, shortestMove, bestMove)
}

func TestFindShortestMoveOneOption(t *testing.T) {
	optimalMoves := make(map[string]float64)
	optimalMoves["up"] = 5.3

	bestMove := "up"
	log.Println(len(optimalMoves))
	shortestMove := getShortestDistance(optimalMoves)
	assert.Equal(t, shortestMove, bestMove)
}

func TestIndexOfShortestDistance(t *testing.T) {
	testMap := make(map[int]float64)
	testMap[0] = 1.2
	testMap[1] = 1.1
	testMap[2] = 1.0

	validate := indexOfShortestDistance(testMap)
	assert.Equal(t, validate, 2)

}

func optimalMovesMapForTest() map[string]float64 {
	testMap := make(map[string]float64)
	testMap["right"] = 1.2
	testMap["up"] = 0.2
	testMap["left"] = 2.0
	testMap["down"] = 1.4
	return testMap
}

func TestFindMostSpaceToMove(t *testing.T) {
	testMap := optimalMovesMapForTest()
	expectedResult := "left"

	validate := getLongestDistance(testMap)
	log.Println(validate)
	assert.Equal(t, validate, expectedResult)
}

func TestFindShortestWay(t *testing.T) {
	testMap := optimalMovesMapForTest()
	expectedResult := "up"

	validate := getShortestDistance(testMap)
	assert.Equal(t, validate, expectedResult)
}
