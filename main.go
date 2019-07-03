package main

import "fmt"
import "image"
import "image/png"
import "image/color"
import "os"

func main() {
WIDTH,HEIGHT:=400,400

	fmt.Println("Starting...")



	s := Sphere{Vec3{4, 0, 0}, 2, color.RGBA{125, 125, 125,255}}
	p := yPlane{-2, color.RGBA{0, 0, 255,255}}
	p2 := yPlane{2, color.RGBA{255, 0, 0,255}}


	objects := []Shape{p,p2,s}

	cam:=Vec3{-10,0,0}

	//later sp;it for threads
	threads:=4

	//var final [][]color.RGBA
	//final:=make( [][][]color.RGBA, threads)
	//fmt.Println(final)
	fmt.Println("DONE=========================")

	units:=HEIGHT/threads
	fmt.Println(units)
/*
	for thread:=0; thread<threads; thread++ {
		fmt.Println("thread")

		start:=(units*thread)
		end:=start+units

		final[thread]=raymarch(WIDTH,HEIGHT,start,end, cam, objects)

	}
*/
	final := make([][]color.RGBA, HEIGHT)
	for r := range final {
		final[r] = make([]color.RGBA, WIDTH)
	}

	final=raymarch(WIDTH,HEIGHT,0,HEIGHT, cam, objects)

    simage := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))



	//use join to combine list of sectioons to 
	//or use img.set in all of them plus the offsett

	//Writing to image
	//join final from multithreading
/*
	out:=join(final, units, threads,WIDTH, HEIGHT)
	fmt.Println(out[399][399])
	-fmt.Println(final[3][399])


*/
/*
	fmt.Println(len(final))

	out := make([][]color.RGBA, HEIGHT)
	for r := range out {
		out[r] = make([]color.RGBA, WIDTH)
	}

	for i:=0; i<len(final); i++{
		start:=i*units
		for y:=0; y<units; y++ {
			fmt.Println(final[0][0][0])
			out[start+y]=final[i][y]
	}}
*/








	for y:=0; y<HEIGHT; y++{
	for x:=0; x<WIDTH; x++{
		simage.Set(x,y,final[y][x])
	}}


	f, _ := os.Create("image.png")
	png.Encode(f, simage)
}
