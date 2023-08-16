// This package doesn`t calculate polynomial. It only interpolate in formal view.
// For calculating this polynomial we need to transform infix view to postfix and then calculate.
// How calculation works, you can check in QAP implementation.

package Lagrangia

import (
	"fmt"
	"strconv"
	"strings"
)

var inputLag string
var mapOfRootsLag = make(map[string]int) // Store coordinates (x,y) in form: key = x or y; value = coordinate value

var mapOfPolynomials = make(map[string]string) // Store interpolated polynomials for 1 point
var coordinates []int                          // Stor coordinates in form [0]=x coord. [1]=y coord [2]=x1 coord. [3]=y1 coord ...

var mapOfMultiPolynomials = make(map[string]string) // Store interpolated polynomials for multi ponts

var mapOfLagrangiaBasis = make(map[int]string) // Store basises of polynomials (needed for barricentric view)

var mapOfNormalNumerator = make(map[string]string)   // Store numerator of interpolated polynomials (multi points interpolation)
var mapOfNormalDenominator = make(map[string]string) // Store denominator of interpolated polynomials (multi points interpolation)

var mapOfBarricNumerator = make(map[int]string)   // Store numerator of interpolated polynomials for barricentric view (multi points interpolation)
var mapOfBarricDenumerator = make(map[int]string) // Store denominator of interpolated polynomials for barricentric view (multi points interpolation)

var multiPolynomialsBarric string

// Saves point coordinates
func rootsMapLag(roots string) {
	strRoots := strings.Split(roots, " ")

	for i := 0; i < len(strRoots); i++ {
		if strRoots[i] == "=" {
			var temp, _ = strconv.Atoi(strRoots[i+1])
			mapOfRootsLag[strRoots[i-1]] = temp
		}
	}

}

// Take only one coordinate pair from coordinates map
func xy(roots map[string]int, key string) (int, int) {
	temp := strings.Split(key, "")
	temp[0] = "y"
	y := strings.Join(temp, "")
	fmt.Println("Point " + "(" + strconv.Itoa(roots[key]) + "," + strconv.Itoa(roots[y]) + ")")
	return roots[key], roots[y]
}

// Saves polynomial wich interpolated from 1 point in mapOfPolynomials variable
func savePolynomial(x int, y int, counter int) {
	tempKey := "pol"
	tempKey = tempKey + strconv.Itoa(counter)

	tempPol := "( " + strconv.Itoa(y) + " * ( x" + " - " + strconv.Itoa(x+y) + " ) ) / ( ( " + strconv.Itoa(x) + " - " + strconv.Itoa(x+y) + " ) )"
	mapOfPolynomials[tempKey] = tempPol

}

// Saves polynomials wich interpolated from multi points in mapOfMultiPolynomials variable
func saveMultPolynomial(xKeys []int, x int, y int, counter int) { //попробовать через поинтеры улучшить (тоесть уменьшить количество переменных)
	tempKey := "pol"
	tempKey = tempKey + strconv.Itoa(counter)

	keysForNumerator := make([]int, len(xKeys))
	keysForDenominator := make([]int, len(xKeys))
	copy(keysForNumerator, xKeys)
	copy(keysForDenominator, xKeys)

	tempVal := "( x - "

	tempValStr := "( " + strconv.Itoa(y) + " * "

	for j := 0; j < len(xKeys)-1; j++ { // Numerator part
		for i := 0; i < len(keysForNumerator); i++ {
			if x == keysForNumerator[i] {
				keysForNumerator = append(keysForNumerator[:i], keysForNumerator[i+1:]...)
				i = i - 1
				continue
			}
			tempVal1 := tempVal + strconv.Itoa(keysForNumerator[i]) + " )"

			keysForNumerator = append(keysForNumerator[:i], keysForNumerator[i+1:]...)

			tempValStr = tempValStr + tempVal1
			break
		}
		if j != len(xKeys)-2 {
			tempValStr = tempValStr + " * "
		}
	}

	mapOfNormalNumerator[tempKey] = tempValStr + " )"
	tempValStr = tempValStr + " )" + " / ( "
	tempVal = "( " + strconv.Itoa(x) + " - "

	for j := 0; j < len(xKeys)-1; j++ { // Denominator part
		for i := 0; i < len(keysForDenominator); i++ {
			if x == keysForDenominator[i] {
				keysForDenominator = append(keysForDenominator[:i], keysForDenominator[i+1:]...)
				i = i - 1
				continue
			}
			tempVal1 := tempVal + strconv.Itoa(keysForDenominator[i]) + " )"

			keysForDenominator = append(keysForDenominator[:i], keysForDenominator[i+1:]...)
			tempValStr = tempValStr + tempVal1
			break
		}
		if j != len(xKeys)-2 {
			tempValStr = tempValStr + " * "
		}
	}

	denum := strings.Split(tempValStr, "/")
	mapOfNormalDenominator[tempKey] = denum[1] + " )"
	tempValStr = tempValStr + " )"
	mapOfMultiPolynomials[tempKey] = tempValStr

}

// Interpolate polynomial by one point
func polynomialByOnePoint(roots map[string]int) {
	fmt.Println("Polynomials by one point: \n")
	counter := 0
	for i := range mapOfRootsLag {
		key := strings.Split(i, "")
		if key[0] == "y" {
			continue
		}

		x, y := xy(roots, i)
		fmt.Println("(" + strconv.Itoa(y) + "(" + "x" + "-" + strconv.Itoa(x+y) + ")" + ")" + "/" + "(" + strconv.Itoa(x) + "-" + strconv.Itoa(x+y) + ")")
		savePolynomial(x, y, counter)

		coordinates = append(coordinates, x)
		coordinates = append(coordinates, y)

		counter++
	}
	fmt.Println("\n")
}

// Interpolated polynomials from multi points
func polynomialByMultiPoints() {

	arraySize := len(mapOfRootsLag) / 2
	keys := make([]int, arraySize)

	counter1 := 0
	for key, value := range mapOfRootsLag { //вытаскиваем x значения
		y := keyGenerate(key, counter1)
		if key == y {
			continue
		}
		keys[counter1] = value
		counter1++
	}

	counter2 := 0
	for key, value := range mapOfRootsLag {
		y := keyGenerate(key, counter2)
		if key == y {
			continue
		}

		saveMultPolynomial(keys, value, mapOfRootsLag[y], counter2)
		//basisCalc(keys, value, mapOfRootsLag[y], counter2)

		counter2++

	}

}

// Create a string wich point to coordinate index. Needed for well coordinates mapping
// Input: "x22". Output: "y22"
func keyGenerate(key string, counter int) (keyNew string) {
	temp := strings.Split(key, "")
	temp[0] = "y"
	keyNew = strings.Join(temp, "")
	return keyNew
}

// Calculate Lagrangia basis and save it.
func basisCalc(xKeys []int, x int, y int, counter int) {

	keysForNumerator := make([]int, len(xKeys))
	keysForDenominator := make([]int, len(xKeys))
	copy(keysForNumerator, xKeys)
	copy(keysForDenominator, xKeys)

	tempValStr := "1 / ( "
	tempVal := "( " + strconv.Itoa(x) + " - "

	for j := 0; j < len(xKeys)-1; j++ { // Denominator part
		for i := 0; i < len(keysForDenominator); i++ {
			if x == keysForDenominator[i] {
				keysForDenominator = append(keysForDenominator[:i], keysForDenominator[i+1:]...)
				i = i - 1
				continue
			}
			tempVal1 := tempVal + strconv.Itoa(keysForDenominator[i]) + " )"

			keysForDenominator = append(keysForDenominator[:i], keysForDenominator[i+1:]...)
			tempValStr = tempValStr + tempVal1
			break
		}
		if j != len(xKeys)-2 {
			tempValStr = tempValStr + " * "
		}
	}
	tempValStr = tempValStr + " )"
	mapOfLagrangiaBasis[y] = tempValStr

	tempVal = "( x - "

	tempValStr = ""

	for j := 0; j < len(xKeys)-1; j++ { // Numerator part
		for i := 0; i < len(keysForNumerator); i++ {
			if x == keysForNumerator[i] {
				keysForNumerator = append(keysForNumerator[:i], keysForNumerator[i+1:]...)
				i = i - 1
				continue
			}
			tempVal1 := tempVal + strconv.Itoa(keysForNumerator[i]) + " )"

			keysForNumerator = append(keysForNumerator[:i], keysForNumerator[i+1:]...)

			tempValStr = tempValStr + tempVal1
			break
		}
		if j != len(xKeys)-2 {
			tempValStr = tempValStr + " * "
		}
	}

	mapOfBarricNumerator[y] = "( " + strconv.Itoa(y) + " * " + mapOfLagrangiaBasis[y] + " / " + "( " + tempValStr + " ) )"
	mapOfBarricDenumerator[y] = mapOfLagrangiaBasis[y] + " / " + "( " + tempValStr + " ) "

}

// Multiply Lagrangia basis on corresponding "y" coordinate
func numeratorBarricentricCals(map[int]string) map[int]string { //DELETE
	//var str []string
	var basisPlusFunc = make(map[int]string)

	for key, _ := range mapOfLagrangiaBasis {

		basisPlusFunc[key] = strconv.Itoa(key) + " * " + mapOfLagrangiaBasis[key]
	}

	return basisPlusFunc
}

// Create Lagrangia interpolate polynomial in barricentric view
func resultBarricentricForm() string {

	fmt.Println("Polynomial by multi points (barricentric view):\n")
	var tempNumerator string
	var tempDenumerator string

	counter := 0
	for _, value := range mapOfBarricNumerator {
		if counter == 0 {
			tempNumerator = value
		} else {
			tempNumerator = tempNumerator + " + " + value
		}
		counter++
	}
	counter = 0
	for _, value := range mapOfBarricDenumerator {
		if counter == 0 {
			tempDenumerator = value
		} else {
			tempDenumerator = tempDenumerator + " + " + value
		}
		counter++
	}
	multiPolynomialsBarric = "( " + tempNumerator + " ) / ( " + tempDenumerator + " )"
	return multiPolynomialsBarric
}

// Print interpolated polynomial in sum form of each polnomial
// Example: output: pol1 + pol2 + pol3...
func resultNormalForm() {
	var strTempResult []string
	fmt.Println("Polynomial by multi points:\n")

	for _, value := range mapOfMultiPolynomials {
		strTempResult = append(strTempResult, value)
		strTempResult = append(strTempResult, " + ")
	}
	strTempResult = strTempResult[:len(strTempResult)-1]
	strResult := strings.Join(strTempResult, "")
	fmt.Println(strResult + "\n")
}

// Print Lagrangia basises for each interpolated poplynomails
func printLagrangiaBasis() {
	fmt.Println("\nLagrangia`s basisis")
	var temp []string
	for key, value := range mapOfLagrangiaBasis {
		x := strings.Split(value, "")
		fmt.Println("Point: " + "(" + x[6] + "," + strconv.Itoa(key) + ")")
		fmt.Println(value)
		temp = append(temp, value)
	}
	fmt.Println("\n")

}

// Return some variables for others packages
func ReturnMultiPolynomial() map[string]string {
	return mapOfMultiPolynomials
}

func ReturnMultiPolynomialBarric() string {

	return multiPolynomialsBarric
}

func ReturnNumeratorNormal() map[string]string {
	return mapOfNormalNumerator
}

func ReturnDenumeratorNormal() map[string]string {
	return mapOfNormalDenominator
}

func ReturnPolynomialByOnePoint() map[string]string {
	return mapOfPolynomials
}

func ReturnBaricBasis() map[int]string {
	return mapOfLagrangiaBasis
}

// Clear all variables
func ClearAllVar() {
	for key := range mapOfBarricDenumerator {
		delete(mapOfBarricDenumerator, key)
	}
	for key := range mapOfRootsLag {
		delete(mapOfRootsLag, key)
	}
	for key := range mapOfPolynomials {
		delete(mapOfPolynomials, key)
	}
	for key := range mapOfNormalNumerator {
		delete(mapOfNormalNumerator, key)
	}
	for key := range mapOfNormalDenominator {
		delete(mapOfNormalDenominator, key)
	}
	for key := range mapOfMultiPolynomials {
		delete(mapOfMultiPolynomials, key)
	}
	for key := range mapOfLagrangiaBasis {
		delete(mapOfLagrangiaBasis, key)
	}
	for key := range mapOfBarricNumerator {
		delete(mapOfBarricNumerator, key)
	}

	inputLag = ""
	multiPolynomialsBarric = ""
	coordinates = nil
}

// Starting function
func Start(input string) {
	inputLag = input

	rootsMapLag(inputLag)

	polynomialByOnePoint(mapOfRootsLag)
	polynomialByMultiPoints()
	//fmt.Println(mapOfMultiPolynomials)

	//numeratorBarricentricCals(mapOfLagrangiaBasis)
	//resultNormalForm()

	//fmt.Println(mapOfLagrangiaBasis)
	//printLagrangiaBasis()
	//fmt.Println(resultBarricentricForm())
	//resultNormalForm()
}
