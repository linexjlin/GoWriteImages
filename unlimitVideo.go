package main

import (
	"bufio"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	ffplay := exec.Command("./ffplay", "-i", "pipe:0", "-pixel_format", "bgr24", "-video_size", strconv.Itoa(100)+"x"+strconv.Itoa(80), "-f", "rawvideo") //nolint
	ffplayIn, _ := ffplay.StdinPipe()
	//ffplayOut, _ := ffplay.StdoutPipe()
	ffplayErr, _ := ffplay.StderrPipe()

	if err := ffplay.Start(); err != nil {
		panic(err)

	}

	go func() {
		scanner := bufio.NewScanner(ffplayErr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	rgbframe := make([]uint8, 100*80*3)

	i := 0
	for {
		img := image.NewRGBA(image.Rect(0, 0, 100, 80))
		addLabel(img, 30, 20, fmt.Sprintf("%d", i))
		i++
		rgba2rgba(img.Pix, rgbframe)
		ffplayIn.Write(rgbframe)
		time.Sleep(time.Millisecond * (1000 / 24))
	}
}

func rgba2rgba(rgba []uint8, rgb []uint8) {
	j := 0
	var r, g, b uint8
	for i := 0; i < len(rgba); i += 4 {
		r, g, b, _ = rgba[i], rgba[i+1], rgba[i+2], rgba[i+3]
		rgb[j], rgb[j+1], rgb[j+2] = b, g, r
		j += 3
	}
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{200, 100, 1, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
