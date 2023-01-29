package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	// check args for filename
	if len(os.Args) != 2 {
		log.Fatalf("wrong number of argumnts provided: %d\n", len(os.Args))
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	original, err := png.Decode(file)
	if err != nil {
		log.Fatalln(err)
	}

	bounds := original.Bounds()
	x, y := bounds.Dx(), bounds.Dy()
	img := image.NewRGBA(bounds)

	wg := sync.WaitGroup{}
	wg.Add(y)
	for j := 0; j < y; j++ {
		go func(j int) {
			defer wg.Done()
			// reverse colors of the opposite points on the line
			for i := 0; i < x-i-1; i++ {
				a, b := original.At(i, j), original.At(x-i-1, j)

				img.Set(x-i-1, j, a)
				img.Set(i, j, b)
			}
		}(j)
	}

	dest, err := os.Create(fmt.Sprintf("%s_reverse.png", strings.TrimSuffix(os.Args[1], ".png")))
	if err != nil {
		log.Fatalln(err)
	}

	defer dest.Close()

	err = png.Encode(dest, img)
	if err != nil {
		log.Fatalf("enoding failed: %v\n", err)
	}

	wg.Wait()
}
