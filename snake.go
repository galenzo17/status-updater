package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
)

type Point struct {
	x, y int
}

type Snake struct {
	body      []Point
	direction Point
}

var (
	snake        Snake
	food         Point
	boardWidth   = 20
	boardHeight  = 10
	gameOver     bool
	tickDuration = 150 * time.Millisecond
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	startGame()
	mainLoop()
}

func startGame() {
	snake = Snake{
		body:      []Point{{boardWidth / 2, boardHeight / 2}},
		direction: Point{0, -1},
	}
	spawnFood()
	gameOver = false
}

func spawnFood() {
	food = Point{rand.Intn(boardWidth), rand.Intn(boardHeight)}
}

func mainLoop() {
	ticker := time.NewTicker(tickDuration)
	defer ticker.Stop()

	for !gameOver {
		select {
		case <-ticker.C:
			update()
			render()
		default:
			handleInput()
		}
	}
}

func update() {
	newHead := Point{
		x: snake.body[0].x + snake.direction.x,
		y: snake.body[0].y + snake.direction.y,
	}

	// Check boundaries
	if newHead.x < 0 || newHead.x >= boardWidth || newHead.y < 0 || newHead.y >= boardHeight {
		gameOver = true
		return
	}

	// Check collision with itself
	for _, part := range snake.body {
		if part == newHead {
			gameOver = true
			return
		}
	}

	// Check if food is eaten
	if newHead == food {
		snake.body = append([]Point{newHead}, snake.body...)
		spawnFood()
	} else {
		// Move snake
		snake.body = append([]Point{newHead}, snake.body[:len(snake.body)-1]...)
	}
}

func render() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	// Draw food
	termbox.SetCell(food.x, food.y, 'F', termbox.ColorRed, termbox.ColorDefault)

	// Draw snake
	for _, part := range snake.body {
		termbox.SetCell(part.x, part.y, 'O', termbox.ColorGreen, termbox.ColorDefault)
	}

	// Draw border
	for x := 0; x < boardWidth; x++ {
		termbox.SetCell(x, 0, '-', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(x, boardHeight-1, '-', termbox.ColorWhite, termbox.ColorDefault)
	}
	for y := 0; y < boardHeight; y++ {
		termbox.SetCell(0, y, '|', termbox.ColorWhite, termbox.ColorDefault)
		termbox.SetCell(boardWidth-1, y, '|', termbox.ColorWhite, termbox.ColorDefault)
	}

	// Refresh screen
	termbox.Flush()
}

func handleInput() {
	switch ev := termbox.PollEvent(); ev.Type {
	case termbox.EventKey:
		switch ev.Key {
		case termbox.KeyArrowUp:
			if snake.direction.y == 0 {
				snake.direction = Point{0, -1}
			}
		case termbox.KeyArrowDown:
			if snake.direction.y == 0 {
				snake.direction = Point{0, 1}
			}
		case termbox.KeyArrowLeft:
			if snake.direction.x == 0 {
				snake.direction = Point{-1, 0}
			}
		case termbox.KeyArrowRight:
			if snake.direction.x == 0 {
				snake.direction = Point{1, 0}
			}
		case termbox.KeyEsc:
			gameOver = true
		}
	}
}

func endScreen() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	printMessage(boardWidth/2-5, boardHeight/2, "Game Over")
	printMessage(boardWidth/2-7, boardHeight/2+1, "Press ESC to quit")
	termbox.Flush()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			break
		}
	}
}

func printMessage(x, y int, msg string) {
	for i, c := range msg {
		termbox.SetCell(x+i, y, c, termbox.ColorWhite, termbox.ColorDefault)
	}
}