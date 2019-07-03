package main

//import "fmt"
import "math"
import "image/color"

type Vec3 struct {
	x float64
	y float64
	z float64
}

func mV(a Vec3, s float64) Vec3 { //multiply vector by scalar
	return Vec3{a.x * s, a.y * s, a.z * s}
}
func aV(a Vec3, b Vec3) Vec3 {
	return Vec3{a.x + b.x, a.y + b.y, a.z + b.z}
}
func Length(a Vec3, b Vec3) float64 {
	return math.Sqrt(math.Pow(a.x-b.x, 2) + math.Pow(a.y-b.y, 2) + math.Pow(a.z-b.z, 2))
}

type Shape interface { //Shpere floor plane(on x and z hasy va,lue)  to test multiple shaes
	DE(o Vec3) (float64, color.RGBA)
}

type Sphere struct {
	pos Vec3
	r   float64
	c   color.RGBA
}
type yPlane struct {
	y float64
	c color.RGBA
}

/*
type Torus struct {
	pos  Vec3
	r, R float64
}
*/

func (s Sphere) DE(o Vec3) (float64, color.RGBA) {
	return Length(o, s.pos) - s.r, s.c
}
func (p yPlane) DE(o Vec3) (float64, color.RGBA) {
	return math.Abs(o.y - p.y), p.c
}

func radians(d float64) float64 {
	return d * (math.Pi / 180)
}
func genVec(ax float64, ay float64) Vec3 {
	rax, ray := radians(ax), radians(ay)
	I := math.Cos(ray)
	return Vec3{I * math.Cos(rax), math.Sin(ray), I * math.Sin(rax)}
}
func raymarch(width int, height int, hS int, hE int, cam Vec3, objects []Shape) [][]color.RGBA { //split by horizontal bars for less arguements      height start height stop
	FOV := 90
	MINDIST := 0.01
	MAXDIST := 2000.0
	MAXSTEPS := 200

	voidColor := color.RGBA{0, 0, 0, 0}

	//create slice
	imgSlice := make([][]color.RGBA, height)
	for r := range imgSlice {
		imgSlice[r] = make([]color.RGBA, width)
	}
	for y := hS; y < hE; y++ {
		ay := (float64(y) * (float64(FOV) / float64(height))) - (float64(FOV) / 2.0)
		for x := 0; x < width; x++ {

			ax := (float64(x) * (float64(FOV) / float64(width))) - (float64(FOV) / 2.0)
			av := genVec(ax, ay)
			ov := cam //origin vector
			fullDist := 0.0
			for i := 0; i < MAXSTEPS; i++ {
				currentDist, currentColor := objects[0].DE(ov) //remember add , color in des
				//fmt.Println(ov)
				//fmt.Println(av)

				for _, o := range objects {
					tmpDist, _ := o.DE(ov)
					if currentDist > tmpDist { //new dist is shorter
						currentDist, currentColor = o.DE(ov)
					}
				}
				if fullDist > MAXDIST {
					imgSlice[y][x] = voidColor
					break
				} else if currentDist <= MINDIST { //stop advancing  use color   add shadows if necessary
					imgSlice[y][x] = currentColor
					//fmt.Println("HIt")
					break
				} else {
					//continue advancing
					ov = aV(ov, mV(av, currentDist))
					fullDist += currentDist
				}
			}

		}
	}
	return imgSlice
} // for goroutines threading

func join(s [][][]color.RGBA, units int,threads int, w int, h int) [][]color.RGBA { // list of 2d arrays    join them to 1 2d array for image
	//create 2d image
	img := make([][]color.RGBA, h)
	for r := range img {
		img[r] = make([]color.RGBA, w)
	}

	//dont use length of arrays because they are larger than needed change that in the future
	for i:=0; i<threads; i++{
		for y:=0; y<units; y++{
		for x:=0; x<w; x++{
			img[(i*units)+y][x]=s[i][y][x]
		}
		}
	}
	return img
}
