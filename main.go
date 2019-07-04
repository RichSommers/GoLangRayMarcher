// TODO ADD reading scene from file

package main

import "fmt"
import "image"
import "image/png"
import "image/color"
import "os"

func main() {
	//a:=Vec3{0,0,0}
	//b:=Vec3{1,0,1}
	//fmt.Println(genAngles(a,b))

	//return

	WIDTH, HEIGHT := 1366, 768

	fmt.Println("Starting...")

	s := Sphere{Vec3{1, 0, 0}, 2, color.RGBA{225, 225, 225, 255}}

	p := yPlane{-7, color.RGBA{0, 0, 255, 255}}
	p2 := yPlane{3, color.RGBA{255, 0, 0, 255}}
	p3 := xPlane{7, color.RGBA{0, 255, 0, 255}}
	p4 := zPlane{-7, color.RGBA{255, 0, 255, 255}}
	p5 := zPlane{7, color.RGBA{0, 255, 255, 255}}



	objects := []Shape{p, p2, p3, p4, p5, s}

	l1 := Light{Vec3{-6, -6, 2}, 1}
	lights := []Light{l1}

	cam := Cam{Vec3{-10, 0, 0}, 0, 0}

	threads := 4

	fmt.Println("DONE=========================")

	units := HEIGHT / threads
	fmt.Println(units)

	final := make([][]color.RGBA, HEIGHT)
	for r := range final {
		final[r] = make([]color.RGBA, WIDTH)
	}

	final = raymarch(WIDTH, HEIGHT, 0, HEIGHT, cam, objects, lights)

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

	for y := 0; y < HEIGHT; y++ {
		for x := 0; x < WIDTH; x++ {
			simage.Set(x, y, final[y][x])
		}
	}

	f, _ := os.Create("images/image.png")
	png.Encode(f, simage)
}
