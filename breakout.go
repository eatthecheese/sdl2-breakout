package main

/* High level overview -
for {
	Update
		get all input
		update all things
	Draw
		draw out the updates
}
*/

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const winWidth, winHeight int = 800, 600
const paddleWidth, paddleHeight int = 100, 20
const brickWidth, brickHeight int = 50, 25
const blockWidthNum, blockHeightNum int = 14, 6

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius float32
	xv, yv float32
	color  color
}

func (ball *ball) draw(pixels []byte) {
	for y := -ball.radius; y < ball.radius; y++ {
		for x := -ball.radius; x < ball.radius; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x+x), int(ball.y+y), color{255, 255, 255}, pixels)
			}
		}
	}
}

func (ball *ball) update(paddle *paddle, bricks bricks) {
	ball.x = ball.x + ball.xv
	ball.y = ball.y + ball.yv

	// handle collisions
	if ball.x-ball.radius < 0 || ball.x+ball.radius > float32(winWidth) {
		ball.xv = -ball.xv
	}
	if ball.y-ball.radius > float32(winHeight) {
		ball.y = 200
	}
	if ball.y-ball.radius < 0 {
		ball.yv = -ball.yv
	}

	// Collision with paddle
	if ball.y+ball.radius > paddle.y-float32(paddle.h/2) && ball.y+ball.radius < paddle.y+float32(paddle.h/2) {
		if ball.x-ball.radius < paddle.x+float32(paddle.w/2) && ball.x+ball.radius > paddle.x-float32(paddle.w/2) {
			//ball.xv = -ball.xv
			ball.yv = -ball.yv
		}
	}

	// Collision with bricks
	for i := range bricks {
		if bricks[i].isAlive {
			if ball.y-ball.radius < bricks[i].y+bricks[i].h/2 {
				if ball.x-ball.radius < bricks[i].x+bricks[i].w/2 && ball.x+ball.radius > bricks[i].x-bricks[i].w/2 {
					//ball.xv = -ball.xv
					ball.yv = -ball.yv
					bricks[i].isAlive = false
				}
			}
			/*else if ball.x-ball.radius < bricks[i].x+bricks[i].w/2 || ball.x+ball.radius > bricks[i].x-bricks[i].w/2 {
				if ball.y-ball.radius < bricks[i].y+bricks[i].h/2 && ball.y+ball.radius > bricks[i].y-bricks[i].h/2 {
					ball.xv = -ball.xv
					bricks[i].isAlive = false
				}
			}
			*/

		}
	}
}

type brick struct {
	pos
	w       float32
	h       float32
	isAlive bool
	color   color
}

type bricks []brick

func (brick *brick) draw(pixels []byte) {

	for y := 0; y < int(brick.h); y++ {
		for x := 0; x < int(brick.w); x++ {
			setPixel(int(brick.x)+x, int(brick.y)+y, color{255, 255, 255}, pixels)
		}
	}

}

type paddle struct {
	pos
	w     int
	h     int
	color color
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.w/2
	startY := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, color{255, 255, 255}, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_RIGHT] != 0 {
		paddle.x = paddle.x + 10
	}
	if keyState[sdl.SCANCODE_LEFT] != 0 {
		paddle.x = paddle.x - 10
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4

	if index < len(pixels)-4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func main() {
	window, err := sdl.CreateWindow("Breakout", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	checkErr(err)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	checkErr(err)
	defer renderer.Destroy()

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	checkErr(err)
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	for y := 0; y < winHeight; y++ {
		for x := 0; x < winWidth; x++ {
			setPixel(x, y, color{0, 0, 0}, pixels)
		}
	}

	player := paddle{pos{400, 550}, paddleWidth, paddleHeight, color{255, 255, 255}}
	ball := ball{pos{400, 200}, 10, 10, 10, color{255, 255, 255}}
	//brick := brick{pos{50, 25}, brickWidth, brickHeight, color{255, 255, 255}}
	var bricks bricks = bricks{}
	// initialise bricks array
	for i := 0; i < blockWidthNum; i++ {
		for j := 0; j < blockHeightNum; j++ {
			bricks = append(bricks, brick{pos{float32(50 + brickWidth*i), float32(25 + brickHeight*j)}, float32(brickWidth), float32(brickHeight), true, color{255, 255, 255}})
		}
	}

	keystate := sdl.GetKeyboardState()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		clear(pixels)

		player.update(keystate)
		ball.update(&player, bricks)

		for i := range bricks {
			if bricks[i].isAlive {
				bricks[i].draw(pixels)
			}
		}
		//bricks.draw(pixels)
		player.draw(pixels)
		ball.draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()
		sdl.Delay(16)
	}

}
