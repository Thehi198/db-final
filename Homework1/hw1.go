// --------------------------------------------------
//
// HOMEWORK 1
//
// Due: Fri, Feb 7, 2025 (23h59)
//
// Name: Mateo Brancoli
//
// Email: mbrancoli@olin.edu
//
// Remarks, if any:
//
//
// --------------------------------------------------
//
// Please fill in this file with your solutions and submit it
//
// The functions below are stubs that you should replace with your
// own implementation.
//
// PLEASE DO NOT CHANGE THE SIGNATURE IN THE STUBS BELOW.
// Doing so makes it impossible for me to test your code.
//
// --------------------------------------------------

package main

import (
	"fmt"
)

func Clamp(min int, max int, v int) int {
	if v > min && v < max {
		return v
	} else {
		return min
	}
}

func Interpolate(min float32, max float32, v float32) float32 {
	if min < max {
		return min + v*min
	} else {
		return min - v*max
	}
}

func Spaces(n int) string {
	space := " "
	string_ := ""
	for i := 0; i < n; i++ {
		string_ = string_ + space
	}
	return string_
}

func PadLeft(s string, n int) string {
	space := " "
	string_ := ""
	if len(s) > n {
		return s
	} else {
		for i := 0; i < n-len(s); i++ {
			string_ = string_ + space
		}
		return string_ + s
	}
}

func PadRight(s string, n int) string {
	space := " "
	string_ := ""
	if len(s) > n {
		return s
	} else {
		for i := 0; i < n-len(s); i++ {
			string_ = string_ + space
		}
		return s + string_
	}
}

func PadBoth(s string, n int) string {
	space := " "
	string_ := ""
	if len(s) > n {
		return s
	} else {
		if (n-len(s))%2 == 0 {
			for i := 0; i < (n-len(s))/2; i++ {
				string_ = string_ + space
			}
			return string_ + s + string_
		} else {
			for i := 0; i < (n-len(s)-1)/2; i++ {
				string_ = string_ + space
			}
			return string_ + s + string_ + space
		}
	}
}

func NewVec(n1 float32, n2 float32, n3 float32) [3]float32 {
	var output = [3]float32{n1, n2, n3}
	return output
}

func ScaleVec(sc float32, v1 [3]float32) [3]float32 {
	var output [3]float32
	for i := 0; i < len(v1); i++ {
		output[i] = v1[i] * sc
	}
	return output
}

func AddVec(v1 [3]float32, v2 [3]float32) [3]float32 {
	var output [3]float32
	for i := 0; i < len(v1); i++ {
		output[i] = v1[i] + v2[i]
	}
	return output
}

func DotProd(v1 [3]float32, v2 [3]float32) float32 {
	var output float32 = 0
	for i := 0; i < len(v1); i++ {
		output = output + v1[i]*v2[i]
	}
	return output
}

func NewMat(r1 [3]float32, r2 [3]float32, r3 [3]float32) [9]float32 {
	var output [9]float32
	var space int = len(r1)
	for i := 0; i < len(r1); i++ {
		output[i] = r1[i]
		output[i+space] = r2[i]
		output[i+2*space] = r3[i]
	}
	return output
}

func ScaleMat(sc float32, m [9]float32) [9]float32 {
	for i := 0; i < len(m); i++ {
		m[i] = m[i] * sc
	}
	return m
}

func TransposeMat(m [9]float32) [9]float32 {
	//var m1 [9]float32 = [9]float32{1, 0, 0, 0, 1, 0, 0, 0, 1}
	numbers := [9]float32{m[0], m[3], m[6], m[1], m[4], m[7], m[2], m[5], m[8]}
	return numbers
}

func AddMat(m1 [9]float32, m2 [9]float32) [9]float32 {
	var m [9]float32
	for i := 0; i < len(m1); i++ {
		m[i] = m1[i] + m2[i]
	}
	return m
}

func MultMat(m1 [9]float32, m2 [9]float32) [9]float32 {
	var m [9]float32
	for i := 0; i < len(m1); i++ {
		for j := 0; j < 3; j++ {
			m[i+3*j] = m1[3*j+0]*m2[0*3+i] + m1[3*j+1]*m2[1*3+i] + m1[3*j+2]*m2[2*3+i]
		}
	}
	return m
}

func main() {
	var f string = "Expected = %v\n     Got = %v\n"
	var wrap func(string) string = func(s string) string { return "'" + s + "'" }
	fmt.Println("You can write some sample tests for yourself here. Here are some to get started.")

	fmt.Println("****** Clamp ***************************************")
	fmt.Printf(f, 10.0, Clamp(10.0, 20.0, 5.0))

	fmt.Println("****** Interpolate ***************************************")
	fmt.Printf(f, 10.0, Interpolate(10.0, 20.0, 0))
	fmt.Printf(f, 20.0, Interpolate(10.0, 20.0, 1))

	fmt.Println("****** Spaces ***************************************")
	fmt.Printf(f, "'          '", wrap(Spaces(10)))

	fmt.Println("****** PadLeft ***************************************")
	fmt.Printf(f, "'      test'", wrap(PadLeft("test", 10)))

	fmt.Println("****** PadRight ***************************************")
	fmt.Printf(f, "'test      '", wrap(PadRight("test", 10)))

	fmt.Println("****** PadBoth ***************************************")
	fmt.Printf(f, "'   test   '", wrap(PadBoth("test", 10)))
	fmt.Printf(f, "'this is a longer test'", wrap(PadBoth("this is a longer test", 10)))

	var v1 [3]float32 = [3]float32{1.0, 2.0, 3.0}

	fmt.Println("****** NewVec ***************************************")
	fmt.Printf(f, [3]float32{0, 10, 20}, NewVec(0, 10, 20))
	fmt.Printf(f, [3]float32{42, 42, 42}, NewVec(42, 42, 42))

	fmt.Println("****** ScaleVec ***************************************")
	fmt.Printf(f, [3]float32{2, 4, 6}, ScaleVec(2.0, v1))

	fmt.Println("****** AddVec ***************************************")
	fmt.Printf(f, [3]float32{2, 4, 6}, AddVec(v1, v1))

	fmt.Println("****** DotProd ***************************************")
	fmt.Printf(f, 14, DotProd(v1, v1))

	var v31 [3]float32 = [3]float32{1, 2, 3}
	var v32 [3]float32 = [3]float32{4, 5, 6}
	var v33 [3]float32 = [3]float32{7, 8, 9}
	var m1 [9]float32 = [9]float32{1, 2, 3, 4, 5, 6, 7, 8, 9}

	fmt.Println("****** NewMat ***************************************")
	fmt.Printf(f, [9]float32{1, 2, 3, 4, 5, 6, 7, 8, 9}, NewMat(v31, v32, v33))

	fmt.Println("****** ScaleMat ***************************************")
	fmt.Printf(f, [9]float32{2, 4, 6, 8, 10, 12, 14, 16, 18}, ScaleMat(2.0, m1))

	fmt.Println("****** TransposeMat ***************************************")
	fmt.Printf(f, [9]float32{1, 4, 7, 2, 5, 8, 3, 6, 9}, TransposeMat(m1))

	fmt.Println("****** AddMat ***************************************")
	fmt.Printf(f, [9]float32{2, 4, 6, 8, 10, 12, 14, 16, 18}, AddMat(m1, m1))

	fmt.Println("****** MultMat ***************************************")
	fmt.Printf(f, [9]float32{30, 36, 42, 66, 81, 96, 102, 126, 150}, MultMat(m1, m1))

}
