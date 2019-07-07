package main

import(
	"fmt"
	"strings"
	"strconv"
	"image/color"

)


func process(file string) ([][]Shape, [][]Light, []Cam, int) {
	//remove comments
	tFil:= strings.Split(file,"\n")
	for i,line:= range tFil{
		if len(line)!=0{
		if line[0]=='#'{
			tFil[i]=""
		}
		}
	}
	fileFrames:=strings.Split(file,"FRAME")[1:]

	for i := range fileFrames{fmt.Println(fileFrames[i])}

	allObjs:=make([][]Shape,len(fileFrames))
	allLights:=make([][]Light,len(fileFrames))
	allCams:=make([]Cam, len(fileFrames))

	for j,frameData := range fileFrames{
		splitData:=strings.Split(frameData,"\n")
		amts:=strings.Fields(splitData[0])
		oAmt,_:=strconv.Atoi(amts[0])
		lAmt,_:=strconv.Atoi(amts[1])
		allObjs[j]=make([]Shape, oAmt)
		allLights[j]=make([]Light, lAmt)


		lNum:=0
		oNum:=0
		for _,line := range strings.Split(frameData,"\n"){
			if len(line)>5{
				fmt.Println(line[0:3])
				if line[0:6]=="sphere"{
					tNs:=strings.Fields(line)[1:]//temp nums
					tN:=make([]float64, len(tNs))
					for i:= range tN{tN[i],_=strconv.ParseFloat(tNs[i],64)}
					allObjs[j][oNum]=Sphere{Vec3{tN[0],tN[1],tN[2]},tN[3],color.RGBA{uint8(tN[4]),uint8(tN[5]),uint8(tN[6]),255}}
					oNum++
				}
				if line[0:3]=="cam"{
					tNs:=strings.Fields(line)[1:]//temp nums
					tN:=make([]float64, len(tNs))
					for i:= range tN{tN[i],_=strconv.ParseFloat(tNs[i],64)}
					allCams[j]=Cam{Vec3{tN[0],tN[1],tN[2]},tN[3],tN[4]}
				}
				if line[0:5]=="light"{
					tNs:=strings.Fields(line)[1:]//temp nums
					tN:=make([]float64, len(tNs))
					for i:= range tN{tN[i],_=strconv.ParseFloat(tNs[i],64)}
					allLights[j][lNum]=Light{Vec3{tN[0],tN[1],tN[2]},tN[3]}
					lNum++
				}
			}
		}
	}
	fmt.Println(allObjs)
	return allObjs, allLights, allCams,len(fileFrames)
}
