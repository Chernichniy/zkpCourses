package QAP

import (
	"github.com/Chernichniy/zkpCourses/goScripts/Lagrangia"
	r1cs "github.com/Chernichniy/zkpCourses/goScripts/R1CS"
)

func main() {
	r1cs.Start()

	Lagrangia.Start()
}
