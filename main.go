// TODO ADD reading scene from file
//using stopwatch
//2000x2000 sphere example no goroutines ~1:04 with 2 threads 36.15 4 threads 34.85   8 threads ~34.84
//using time ./program   real
//using computer for other things at this point
//4000x4000 same scene   1 thread 5m02  threads 2m38  4 threads  2m39  8 threads 2m44
package main

import "fmt"
import "sync"
import "image"
import "image/png"
import "image/color"
import "os"

func main() {
	frames:=1
	threads := 4

	WIDTH, HEIGHT := 1600, 1600


	fmt.Println("Starting...")



	cam := Cam{Vec3{-10, 0, 0}, 0, 0}



	units := HEIGHT / threads
	fmt.Println("barSize",units)



	final := make([][]color.RGBA, HEIGHT)
	for r := range final {
		final[r] = make([]color.RGBA, WIDTH)
	}
for i:=0;i<frames;i++{


	p := yPlane{-10, color.RGBA{123, 186, 245, 255}}
	p2 := yPlane{10, color.RGBA{220, 88, 21, 255}}
	p3 := xPlane{7, color.RGBA{123, 186, 245, 255}}
	//p4 := zPlane{-10, color.RGBA{123, 186,245, 255}}
	//p5 := zPlane{10, color.RGBA{123, 186, 245, 255}}

	s:=Sphere{Vec3{6, -5, 0}, 0.5, color.RGBA{255, 255, 255,255}}
	s1:=Sphere{Vec3{6, -4.7, 0.4}, 0.3, color.RGBA{255, 255, 255,255}}
	s2:=Sphere{Vec3{6, -4.7, -0.4}, 0.3, color.RGBA{255, 255, 255,255}}



	l1 := Light{Vec3{-6, -6, 4}, 1}
	lights := []Light{l1} //for multiple lights combine em




	objects := []Shape{p, p2, p3, s,s1,s2}

	//,ulti threading yeehaw
	output:=make([][][]color.RGBA,threads)
	bar:=make([][]color.RGBA, units)
	line:=make([]color.RGBA, WIDTH)
	for j:= range bar{
		bar[j]=line
	}
	for k:= range output{
		output[k]=bar
	}

	var wg sync.WaitGroup
	wg.Add(threads) // we're going to wait for 1 person to finit

	for i:=0; i<threads; i++{
		start:=i*units
		end:=start+units
		fmt.Println(start,end)

		go func(num int,WIDTH int , HEIGHT int , start int, end int , cam Cam, objects []Shape, lights []Light) {
			output[num] = raymarch(WIDTH, HEIGHT, start, end, cam, objects, lights)
			wg.Done()//avoid race condition i think

		}(i,WIDTH, HEIGHT, start, end, cam, objects, lights)
	}
	wg.Wait() //when done all is finished-
	final:=join(output, WIDTH,HEIGHT,units)
	simage := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	//use join to combine list of sectioons to
	//or use img.set in all of them plus the offsett    <---- This is prolly faster use code from join tho

	//Writing to image
	//join final from multithreading

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			simage.Set(x, y, final[y][x])
		}
	}

	f, _ := os.Create( fmt.Sprintf("%s%d%s","video/image",i,".png") )
	png.Encode(f, simage)
	fmt.Println("Finished Frame ",i)
}

}
