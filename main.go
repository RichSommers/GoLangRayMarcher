// TODO ADD reading scene from file
//using stopwatch
//2000x2000 sphere example no goroutines ~1:04 with 2 threads 36.15 4 threads 34.85   8 threads ~34.84
//using time ./program   real
//using computer for other things at this point
//4000x4000 same scene   1 thread 5m02  threads 2m38  4 threads  2m39  8 threads 2m44
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv" //for sys.args  and file
//	"strings"
	"sync"
)

func main() {

	if len(os.Args) != 7 {

		fmt.Println("Incorrect Arguments")
		fmt.Println("raymarch width, height, FOV, # of threads")
		os.Exit(1)
	}
	args := os.Args[1:] //without program
	WIDTH, _ := strconv.Atoi(args[0])
	HEIGHT, _ := strconv.Atoi(args[1])
	FOV, _ := strconv.Atoi(args[2])
	threads, _ := strconv.Atoi(args[3])
	outFilePath := args[4]
	inFilePath := args[5]


	fmt.Println("Starting...")

	//get file pass to process()
	//split by "FRAME"[1:]
	//[][]arrayof all (cam,objects,lights) of length array split by frame
	//then splitlines  then first line.Fields Atoi == numObjects,numLights
	//this frame of bigger array = make([]Shapes,numObjects)

	file, err := os.Open(inFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fil, err := ioutil.ReadAll(file)

	allObjects, allLights, allCams,frames := process(string(fil))




	units := HEIGHT / threads
	final := make([][]color.RGBA, HEIGHT)
	for r := range final {
		final[r] = make([]color.RGBA, WIDTH)
	}
	for i := 0; i < frames; i++ {

		objects := allObjects[i]
		lights := allLights[i]
		cam:=allCams[i]

		//multi threading yeehaw
		output := make([][][]color.RGBA, threads)
		bar := make([][]color.RGBA, units)
		line := make([]color.RGBA, WIDTH)
		for j := range bar {
			bar[j] = line
		}
		for k := range output {
			output[k] = bar
		}

		var wg sync.WaitGroup
		wg.Add(threads) // we're going to wait for 1 person to finit

		for i := 0; i < threads; i++ {
			start := i * units
			end := start + units

			go func(num int, WIDTH int, HEIGHT int, start int, end int, cam Cam, objects []Shape, lights []Light, FOV int, threadNum int) {
				output[num] = raymarch(WIDTH, HEIGHT, start, end, cam, objects, lights, FOV, threadNum)
				wg.Done() //avoid race condition i think

			}(i, WIDTH, HEIGHT, start, end, cam, objects, lights, FOV, i)
		}
		wg.Wait() //when done all is finished-

		simage := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

		//join final from multithreading

		for i, bar := range output {
			for y, line := range bar {
				ypos := y + (i * units)
				for x, col := range line {
					simage.Set(x, ypos, col)
				}
			}
		}

		nowOutFilePath:=outFilePath
		if frames>1 {
			nowOutFilePath=fmt.Sprintf("%s%d%s",outFilePath,i,".png")
		}


		f, _ := os.Create(nowOutFilePath)
		png.Encode(f, simage)
		fmt.Println("Finished Frame ", i)
		fmt.Println("Saved at",nowOutFilePath)
	}

}
