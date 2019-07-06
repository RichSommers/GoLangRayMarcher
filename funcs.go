package main

import "fmt"
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

type Cam struct {
	pos Vec3
	ax  float64
	ay  float64
}

type Light struct {
	pos      Vec3
	strength float64
	//c color.RGBA
}

type Shape interface { //Shpere floor plane(on x and z hasy va,lue)  to test multiple shaes
	DE(o Vec3) (float64, color.RGBA)
}

type Sphere struct {
	pos Vec3
	r   float64
	c   color.RGBA
}

func (s Sphere) DE(o Vec3) (float64, color.RGBA) {
	return Length(o, s.pos) - s.r, s.c
}

type yPlane struct {
	y float64
	c color.RGBA
}

func (p yPlane) DE(o Vec3) (float64, color.RGBA) {
	return math.Abs(o.y - p.y), p.c
}

type xPlane struct {
	x float64
	c color.RGBA
}

func (p xPlane) DE(o Vec3) (float64, color.RGBA) {
	return math.Abs(o.x - p.x), p.c
}

type zPlane struct {
	z float64
	c color.RGBA
}

func (p zPlane) DE(o Vec3) (float64, color.RGBA) {
	return math.Abs(o.z - p.z), p.c
}

/*
type Torus struct {
	pos  Vec3
	r, R float64
}
*/
func shadowColor(c color.RGBA, s float64) color.RGBA {
	//s=s/2 //blend shadow with environment
	nc := color.RGBA{uint8(float64(c.R) * s), uint8(float64(c.G) * s), uint8(float64(c.B) * s), 255}
	return nc
}
func combineColor(a color.RGBA, b color.RGBA) color.RGBA {
	fmt.Println("to not uncomment fmt every time i need to dn")
	return color.RGBA{(a.R + b.R) / 2, (a.G + b.G) / 2, (a.B + b.B) / 2, (a.A + b.A) / 2}
}

func radians(d float64) float64 {
	return d * (math.Pi / 180)
}
func degrees(r float64) float64 {
	return r * (180 / math.Pi)
}
func min(v ...float64) float64 {
	min := v[0]
	for _, i := range v {
		if min > i {
			min = i
		}
	}
	return min
}

func genVec(ax float64, ay float64) Vec3 {
	//rax, ray := radians(ax), radians(ay)
	I := math.Cos(ay)
	return Vec3{I * math.Cos(ax), math.Sin(ay), I * math.Sin(ax)}
}

func genAngles(a Vec3, b Vec3) Vec3 { //incorrectly mnamed but wjayever   also this is wrong
	t := Length(a, b)
	ax := math.Atan2((b.z - a.z), (b.x - a.x)) //changed from .x
	ay := math.Asin((b.y - a.y) / t)
	return genVec(ax, ay)
}

func marchShadow(ov Vec3, objects []Shape, lights []Light) float64 {
	k := 16 //smoth shadows    higher is harder lower is sofgter shadow
	MINDIST := 0.000000001
	//0-1 0 or 1 for hard shadows do soft shad0ws soon
	//get advance ray for each light source
	//return min of all light sources (0 or 1 for hard shadows)
	currentDist, _ := objects[0].DE(ov)
	res := 1.0
	s := make([]float64, len(lights))

	for i, l := range lights {
		MAXDIST := Length(ov, l.pos)
		fullDist := 0.0
		av := genAngles(ov, l.pos)
		for {

			if currentDist <= MINDIST { //hitting an object
				s[i] = 0 //shadowed
				break
			} else if fullDist >= MAXDIST {
				s[i] = res //not blocked
				break
			} else {
				currentDist, _ = objects[0].DE(ov)
				for _, o := range objects {
					tmpDist, _ := o.DE(ov)
					if currentDist > tmpDist { //new dist is shorter
						currentDist, _ = o.DE(ov)
					}
					res = min((float64(k) * currentDist / fullDist), res)
				}
				fullDist += currentDist
				ov = aV(ov, mV(av, currentDist))

			}
		}
	}
	sum := 0.0
	for _, f := range s {
		sum += f
	}
	return sum / float64(len(s))
}

func raymarch(width int, height int, hS int, hE int, cam Cam, objects []Shape, lights []Light, FOV int) [][]color.RGBA { //split by horizontal bars for less arguements      height start height stop
	FOVF := float64(FOV)
	FOVX, FOVY := 0.0, 0.0
	if width > height {
		FOVX, FOVY = FOVF, (float64(height)/float64(width))*FOVF
	} else {
		FOVY, FOVX = FOVF, (float64(width)/float64(height))*FOVF
	}

	MINDIST := 0.0001
	MAXDIST := 20000.0
	MAXSTEPS := 1000

	voidColor := color.RGBA{0, 0, 0, 0}

	//create slice
	imgSlice := make([][]color.RGBA, hE-hS)
	for r := range imgSlice {
		imgSlice[r] = make([]color.RGBA, width)
	}
	for y := hS; y < hE; y++ {
		ay := (float64(y) * (float64(FOVY) / float64(height))) - (float64(FOVY) / 2.0)
		ay += cam.ay
		for x := 0; x < width; x++ {

			ax := (float64(x) * (float64(FOVX) / float64(width))) - (float64(FOVX) / 2.0)
			ax += cam.ax
			av := genVec(radians(ax), radians(ay))
			ov := cam.pos //origin vector
			fullDist := 0.0
			for i := 0; i < MAXSTEPS; i++ {
				currentDist, currentColor := objects[0].DE(ov) //remember add , color in des

				for _, o := range objects {
					tmpDist, _ := o.DE(ov)
					if currentDist > tmpDist { //new dist is shorter
						currentDist, currentColor = o.DE(ov)
					}
				}

				if fullDist > MAXDIST { //too far
					imgSlice[y-hS][x] = voidColor
					break
				} else if currentDist <= MINDIST { //stop advancing  use color   add shadows if necessary
					sh := (marchShadow(ov, objects, lights) / 2) + 0.5 //shadow amount
					newColor := shadowColor(currentColor, sh)          //make new with shadow
					imgSlice[y-hS][x] = newColor                       // set color
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
}

// for goroutines threading

func join(in [][][]color.RGBA, w int, h int, units int) [][]color.RGBA {
	out := make([][]color.RGBA, h)
	for i := range out {
		out[i] = make([]color.RGBA, w)
	}
	for i, bar := range in {
		for y, line := range bar {
			out[y+(i*units)] = line

		}
	}

	return out
}
