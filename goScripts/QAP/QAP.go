package main

import (
	"github.com/Chernichniy/zkpCourses/goScripts/Lagrangia"
	r1cs "github.com/Chernichniy/zkpCourses/goScripts/R1CS"
)

func main() {
	r1cs.Start("x ^ 2 + 2 * x * z + z ^ 2 + z + 1")
	Lagrangia.Start("x = 3 x2 = 1 x3 = 2 y = 5 y2 = 7 y3 = 3")
}
