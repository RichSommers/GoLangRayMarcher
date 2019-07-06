// TODO ADD reading scene from file

package main

import "fmt"
import "image"
import "image/png"
import "image/color"
import "os"

func main() {
	frames:=1

	WIDTH, HEIGHT := 400, 400


	fmt.Println("Starting...")



	cam := Cam{Vec3{-10, 0, 0}, 0, 0}



	threads := 4

	fmt.Println("DONE=========================")

	units := HEIGHT / threads
	fmt.Println("barSize",units)



	final := make([][]color.RGBA, HEIGHT)
	for r := range final {
		final[r] = make([]color.RGBA, WIDTH)
	}
for i:=0;i<frames;i++{


	p := yPlane{-10, color.RGBA{0, 0, 255, 255}}
	p2 := yPlane{10, color.RGBA{255, 0, 0, 255}}
	p3 := xPlane{7, color.RGBA{0, 255, 0, 255}}
	p6 := xPlane{-22, color.RGBA{255, 255, 0, 255}}
	p4 := zPlane{-10, color.RGBA{255, 0, 255, 255}}
	p5 := zPlane{10, color.RGBA{0, 255, 255, 255}}

	s:=Sphere{Vec3{1, -5-(float64(i)/-20), 0}, 2, color.RGBA{100, 100, 100,255}}

	s1:=Sphere{Vec3{-18, -5-(float64(i)/-20), 0}, 2, color.RGBA{100, 100, 100,255}}


	l1 := Light{Vec3{-6, -6, 4}, 1}
	lights := []Light{l1} //for multiple lights combine em




	objects := []Shape{p, p2, p3, p4, p5,p6, s,s1}

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
	for i:=0; i<threads; i++{
		start:=i*units
		end:=start+units
		fmt.Println(start,end)
		output[i] = raymarch(WIDTH, HEIGHT, start, end, cam, objects, lights)
	}

	final:=join(output, WIDTH,HEIGHT,units)
	simage := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	//use join to combine list of sectioons to
	//or use img.set in all of them plus the offsett

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
