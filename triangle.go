package main

import "math"
import "image/color"



func subVec(a, b Vec3) Vec3 {
	return Vec3{a.x - b.x, a.y - b.y, a.z - b.z}
}
func dot(a, b Vec3) float64 {
	return (a.x * b.x) + (a.y * b.y) + (a.z * b.z)
}


type Tri struct {
	pts []Vec3
	c color.RGBA
}

func(tri Tri) DE(o Vec3) (float64, color.RGBA) {
	//translated from https://gist.githubusercontent.com/joshuashaffer/99d58e4ccbd37ca5d96e/raw/8662caaaf344dbe4bc89bcff21f12fc116eda93f/pointTriangleDistance.py
	B := tri.pts[0]
	E0 := subVec(tri.pts[1], B)
	E1 := subVec(tri.pts[2], B)
	D := subVec(B, o)

	a := dot(E0, E0)
	b := dot(E0, E1)
	c := dot(E1, E1)
	d := dot(E0, D)
	e := dot(E1, D)
	f := dot(D, D)

	det := (a * c) - (b * b)
	s := (b * e) - (c * d)
	t := (b * d) - (a * e)
	sqrdistance := 0.0
	invDet:=0.0
	tmp0:=0.0
	tmp1:=0.0
	numer:=0.0
	denom:=0.0

	//here i quote "Terible tree of conditionals to determine in which region of the diagram shown above the projection of the point into the triangle-plane lies"

	if (s + t) <= det {
		if s < 0.0 {
			if t < 0.0 {
				// region4
				if d < 0 {
					t = 0.0
					if -d >= a {
						s = 1.0
						sqrdistance = a + 2.0*d + f
					} else {
						s = -d / a
						sqrdistance = d*s + f
					}
				} else {
					s = 0.0
					if e >= 0.0 {
						t = 0.0
						sqrdistance = f
					} else {
						if -e >= c {
							t = 1.0
							sqrdistance = c + 2.0*e + f
						} else {
							t = -e / c
							sqrdistance = e*t + f
						}
					}
				}
				// of region 4
			} else {
				// region 3
				s = 0
				if e >= 0 {
					t = 0
					sqrdistance = f
				} else {
					if -e >= c {
						t = 1
						sqrdistance = c + 2.0*e + f
					} else {
						t = -e / c
						sqrdistance = e*t + f
						// of region 3
					}
				}
			}
		} else {
			if t < 0 {
				// region 5
				t = 0
				if d >= 0 {
					s = 0
					sqrdistance = f
				} else {
					if -d >= a {
						s = 1
						sqrdistance = a + 2.0*d + f // GF 20101013 fixed typo d*s ->2*d
					} else {
						s = -d / a
						sqrdistance = d*s + f
					}
				}
			} else {
				// region 0
				invDet = 1.0 / det
				s = s * invDet
				t = t * invDet
				sqrdistance = s*(a*s+b*t+2.0*d) + t*(b*s+c*t+2.0*e) + f
			}
		}
	} else {
		if s < 0.0 {
			// region 2
			tmp0 = b + d
			tmp1 = c + e
			if tmp1 > tmp0 { // minimum on edge s+t=1
				numer = tmp1 - tmp0
				denom = a - 2.0*b + c
				if numer >= denom {
					s = 1.0
					t = 0.0
					sqrdistance = a + 2.0*d + f // GF 20101014 fixed typo 2*b -> 2*d
				} else {
					s = numer / denom
					t = 1 - s
					sqrdistance = s*(a*s+b*t+2*d) + t*(b*s+c*t+2*e) + f
				}
			} else { // minimum on edge s=0
				s = 0.0
				if tmp1 <= 0.0 {
					t = 1
					sqrdistance = c + 2.0*e + f
				} else {
					if e >= 0.0 {
						t = 0.0
						sqrdistance = f
					} else {
						t = -e / c
						sqrdistance = e*t + f
						// of region 2
					}
				}
			}
		} else {
			if t < 0.0 {
				// region6
				tmp0 = b + e
				tmp1 = a + d
				if tmp1 > tmp0 {
					numer = tmp1 - tmp0
					denom = a - 2.0*b + c
					if numer >= denom {
						t = 1.0
						s = 0
						sqrdistance = c + 2.0*e + f
					} else {
						t = numer / denom
						s = 1 - t
						sqrdistance = s*(a*s+b*t+2.0*d) + t*(b*s+c*t+2.0*e) + f
					}
				} else {
					t = 0.0
					if tmp1 <= 0.0 {
						s = 1
						sqrdistance = a + 2.0*d + f
					} else {
						if d >= 0.0 {
							s = 0.0
							sqrdistance = f
						} else {
							s = -d / a
							sqrdistance = d*s + f
						}
					}
				}
			} else {
				// region 1
				numer = c + e - b - d
				if numer <= 0 {
					s = 0.0
					t = 1.0
					sqrdistance = c + 2.0*e + f
				} else {
					denom = a - 2.0*b + c
					if numer >= denom {
						s = 1.0
						t = 0.0
						sqrdistance = a + 2.0*d + f
					} else {
						s = numer / denom
						t = 1 - s
						sqrdistance = s*(a*s+b*t+2.0*d) + t*(b*s+c*t+2.0*e) + f
					}
				}
			}
		}
	}
	if sqrdistance < 0 {
		sqrdistance = 0
	}

	dist := math.Sqrt(sqrdistance)
	return dist,tri.c
}



