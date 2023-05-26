package Lagrangia

import (
	"fmt"
	"strconv"
	"strings"
)

var inputLag = "x = 1 x2 = 2 x3 = 3 y = 3 y2 = 2 y3 = 1"
var mapOfRootsLag = make(map[string]int)

func rootsMapLag(roots string) {
	strRoots := strings.Split(roots, " ")

	for i := 0; i < len(strRoots); i++ {
		if strRoots[i] == "=" {
			var temp, _ = strconv.Atoi(strRoots[i+1])
			mapOfRootsLag[strRoots[i-1]] = temp //помещает корни в key=value форму
		}
	}
	fmt.Println(mapOfRootsLag)

}

func mainLag() {
	rootsMapLag(inputLag)
}
