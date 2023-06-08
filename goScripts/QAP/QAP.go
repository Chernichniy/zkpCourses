package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Chernichniy/zkpCourses/goScripts/Lagrangia"
	r1cs "github.com/Chernichniy/zkpCourses/goScripts/R1CS"
)

var qapVectorA [][]float64
var qapVectorB [][]float64
var qapVectorC [][]float64
var qapVectorTemp [][]float64

var indexOfRow int = 0

// Check if variable float64 is int variable
// Example: isInt(2,0) == true; isInt(2,1)==false
func isInt(val float64) bool { //поместить в отдельный пакет
	return val == float64(int(val))
}

// Check if variable is finit
// Finite = true; Infinite = false.
func isFinite(num float64) bool {
	return !math.IsInf(num, 0) && !math.IsNaN(num)
}

// Initialized two-demesional slice
func vectorQAPSizeAllocate(r1csVecor [][]int, vectName string) {
	var temp = make([][]float64, len(r1csVecor[0]))
	for i := range temp {
		temp[i] = make([]float64, len(r1csVecor))
	}

	switch vectName {
	case "A":
		qapVectorA = temp
	case "B":
		qapVectorB = temp
	case "C":
		qapVectorC = temp
	case "temp":
		qapVectorTemp = temp
	}

}

// Calculating infix view statement
// Example: input(34+56+*) -> 756+* -> 711* -> output(77)
func calculateInfixView(rpn string) (res int) {
	str := strings.Split(rpn, " ")

	for i := 0; i < len(str); i++ {
		if i+2 <= len(str)-1 {
			num1, _ := strconv.ParseFloat(str[i], 64)
			num2, _ := strconv.ParseFloat(str[i+1], 64)
			num3, _ := strconv.ParseFloat(str[i+2], 64)
			if num3 == 0 {

				switch str[i+2] {
				case "+":
					num := num1 + num2
					if isInt(num) {
						str[i+2] = strconv.Itoa(int(num))
					} else {
						str[i+2] = fmt.Sprintf("%f", num)
					}
					str = append(str[:i], str[i+2:]...)
					i = -1
				case "*":
					num := num1 * num2
					if isInt(num) {
						str[i+2] = strconv.Itoa(int(num))
					} else {
						str[i+2] = fmt.Sprintf("%f", num)
					}

					str = append(str[:i], str[i+2:]...)
					i = -1
				case "^":
					numTemp := num1
					num := num1

					for j := num2; j > 1; j-- {
						num = numTemp * num1
						numTemp = num
					}
					str[i+2] = strconv.FormatFloat(num, 'E', -1, 8)
					str = append(str[:i], str[i+2:]...)
					i = -1
				case "-":

					num := num1 - num2

					if isInt(num) { // If Int or float then...
						str[i+2] = strconv.Itoa(int(num))
					} else {
						str[i+2] = fmt.Sprintf("%f", num)
					}
					str = append(str[:i], str[i+2:]...)
					i = -1
				case "/":
					num := num1 / num2
					if isFinite(num) == false {
						num = 0
					}
					if isInt(num) {
						str[i+2] = strconv.Itoa(int(num))
					} else {
						str[i+2] = fmt.Sprintf("%f", num)
					}
					str = append(str[:i], str[i+2:]...)
					i = -1
				}

			}
		}
	}
	res, _ = strconv.Atoi(str[0])
	return res
}

// Deleted polynommials from map with interpolated polynomials wich coordinates are (x,0)
func deleteZerosMapValue(numerator map[string]string, denumerator map[string]string) {
	for key, value := range numerator {
		tempStr := strings.Split(value, " ")
		if tempStr[1] == "0" {
			delete(numerator, key)
			delete(denumerator, key)
		}
	}
	//return functions
}

// Multiplies left and right inputs by multiplying summing degrees and multiplyign coefficients
// Example: [34][x][2] * [2][x][4] = [68][x][8]
func leftRightInputMult(left []string, right []string, indexXLeft int, indexXRight int) string {
	var coeffLeftForLeftInput float64 = 1
	var coeffRightForLeftInput float64 = 0

	var coeffLeftForRightInput float64 = 1
	var coeffRightForRightInput float64 = 0

	var tempLeftCoefForLeft string = ""
	var tempRightCoefForLeft string = ""

	var tempLeftCoefForRight string = ""
	var tempRightCoefForRight string = ""

	// Units numbers parts: [3][4][x][2][1] to [34][x][21]
	switch left[indexXLeft] {
	case "x":

		for i := 0; i < indexXLeft; i++ {
			tempLeftCoefForLeft = tempLeftCoefForLeft + left[i]
		}

		for i := indexXLeft + 1; i < len(left); i++ {
			tempRightCoefForLeft = tempRightCoefForLeft + left[i]
		}

	default:
		for i := 0; i < len(left); i++ {
			tempLeftCoefForLeft = tempLeftCoefForLeft + left[i]
		}
	}

	switch right[indexXRight] {
	case "x":

		for i := 0; i < indexXRight; i++ {
			tempLeftCoefForRight = tempLeftCoefForRight + right[i]
		}

		for i := indexXRight + 1; i < len(right); i++ {
			tempRightCoefForRight = tempRightCoefForRight + right[i]
		}
	default:
		for i := 0; i < len(right); i++ {
			tempLeftCoefForRight = tempLeftCoefForRight + right[i]
		}
	}

	// String -> float64: Assign coefficients and degrees values of left and right multipliers
	switch {
	case len(tempLeftCoefForLeft) > 0:
		coeffLeftForLeftInput, _ = strconv.ParseFloat(tempLeftCoefForLeft, 64)
	}

	switch {
	case len(tempRightCoefForLeft) > 0:
		coeffRightForLeftInput, _ = strconv.ParseFloat(tempRightCoefForLeft, 64)

	}

	switch {
	case len(tempLeftCoefForRight) > 0:
		coeffLeftForRightInput, _ = strconv.ParseFloat(tempLeftCoefForRight, 64)
	}

	switch {
	case len(tempRightCoefForRight) > 0:
		coeffRightForRightInput, _ = strconv.ParseFloat(tempRightCoefForRight, 64)
	}

	// Here float64*int and float64*int issues are solves
	switch {

	case isInt(coeffLeftForLeftInput*coeffRightForLeftInput) == true && isInt(coeffLeftForRightInput*coeffRightForRightInput) == true:
		return strconv.Itoa(int(coeffLeftForLeftInput*coeffLeftForRightInput)) + "x" + strconv.Itoa(int(coeffRightForLeftInput+coeffRightForRightInput))

	case isInt(coeffLeftForLeftInput*coeffRightForLeftInput) == false && isInt(coeffLeftForRightInput*coeffRightForRightInput) == false:
		return fmt.Sprintf("%f", coeffLeftForLeftInput*coeffLeftForRightInput) + "x" + fmt.Sprintf("%f", coeffRightForLeftInput+coeffRightForRightInput)

	case isInt(coeffLeftForLeftInput*coeffRightForLeftInput) == false && isInt(coeffLeftForRightInput*coeffRightForRightInput) == true:
		return fmt.Sprintf("%f", coeffLeftForLeftInput*coeffLeftForRightInput) + "x" + strconv.Itoa(int(coeffRightForLeftInput+coeffRightForRightInput))

	case isInt(coeffLeftForLeftInput*coeffRightForLeftInput) == true && isInt(coeffLeftForRightInput*coeffRightForRightInput) == false:
		return strconv.Itoa(int(coeffLeftForLeftInput*coeffLeftForRightInput)) + "x" + fmt.Sprintf("%f", coeffRightForLeftInput+coeffRightForRightInput)

	}

	return "" // If there's no interpolated polynomial for proccess of function
}

// 1. Doing this thing: [3][4][x][4][6] -> [34][x][46]
// 2. Find index of variable "x" in monomial slice representation
func leftAndRightOfOneValueInputAppend(input []string) ([]string, int) {
	var tempValLeftSide []string
	var tempValRightSide []string
	var index int

	for i := 0; i < len(input); i++ { // [3][4][x][4][6] -> [34][x][46]
		if input[i] == "x" && i > 1 {

			for j := 0; j < i; j++ {
				tempValLeftSide = append(tempValLeftSide, input[j])
			}
			leftInputStr := strings.Join(tempValLeftSide, "") //ТУТ
			input[i-1] = leftInputStr
			input = append(input[:0], input[i-1:]...)
			i = -1

		} else if input[i] == "x" && i+2 != len(input) || input[0] == "x" && i+2 != len(input) {

			for j := i + 1; j < len(input); j++ {
				tempValRightSide = append(tempValRightSide, input[j])
			}
			leftInputStr := strings.Join(tempValRightSide, "") //ТУТ
			input[i+1] = leftInputStr
			input = append(input[:i+2], input[len(input):]...)

		}
	}
	for i := range input { // Loop for index finding
		if input[i] == "x" {
			index = i
		}
	}

	return input, index
}

// 1. Multiplies calculated coeeficint of monomial(after brackets opening) by "y" coordinate
// 2. Sums coefficient of corresponding "x" degrees. In other words, simplifies the expression.
// Example: input (x^3-3x^2+2x-3x^2+9x-6) -> output (x^3-6x^2+11x-6)
func additionMonoms(monoms []string) map[string]float64 {
	copiedMonoms := make([]string, len(monoms))
	copy(copiedMonoms, monoms)

	var mapOfCoeff = make(map[string]float64)

	for j := 1; j < len(copiedMonoms); j++ { // Devided cooeficients of variables by degree of "x", then summing coefficients by corresponding degrees.
		// If there no x variable in value, the value writing like this: [coefficient][x][0]
		// Example: value "6" writes like 6x0
		tempMonomSlice := strings.Split(copiedMonoms[j], "x")
		val, isExist := mapOfCoeff[tempMonomSlice[1]]

		if isExist {
			tempValueFromMap, _ := strconv.ParseFloat(tempMonomSlice[0], 64)
			mapOfCoeff[tempMonomSlice[1]] = val + tempValueFromMap
		} else {
			mapOfCoeff[tempMonomSlice[1]], _ = strconv.ParseFloat(tempMonomSlice[0], 64)
		}
	}

	yValue, _ := strconv.ParseFloat(copiedMonoms[0], 64)
	for key, value := range mapOfCoeff { // Multiply coefficients on "y" coordinate

		mapOfCoeff[key] = value * yValue
	}

	return mapOfCoeff
}

// Start procces of multiplying monomials from brackets (x-1)(x-2)...
func multiplyMonoms(str string, signForLeftNum string, signForRightNum string) string {

	strTempSlice := strings.Split(str, "*")

	tempLeftInput := strings.Split(strTempSlice[0], "")  // Store left monomial
	tempRightInput := strings.Split(strTempSlice[1], "") // Store right monomial

	// (For left monomial) Checks if there x3-* are "-" so 3 => -3 and so on
	if signForLeftNum == "-" && tempLeftInput[0] != "x" {
		tempLeftInput[0] = "-" + tempLeftInput[0]
	} else if signForLeftNum == "-" && tempLeftInput[0] == "x" {
		tempLeftInput = append(tempLeftInput[:1], tempLeftInput[0:]...)
		tempLeftInput[0] = "-1"
	}
	// (For right monomial) Checks if there x3-* are "-" so 3 => -3 and so on
	if signForRightNum == "-" && tempRightInput[0] != "x" {
		tempRightInput[0] = "-" + tempRightInput[0]
	} else if signForRightNum == "-" && tempRightInput[0] == "x" {
		tempRightInput = append(tempRightInput[:1], tempRightInput[0:]...)
		tempRightInput[0] = "-1"
	}

	var indexOfValueXForLeftInput int
	var indexOfValueXForRightInput int

	// Need for storing index of "x" value from slice. It's needed for next calculations
	tempLeftInput, indexOfValueXForLeftInput = leftAndRightOfOneValueInputAppend(tempLeftInput)
	tempRightInput, indexOfValueXForRightInput = leftAndRightOfOneValueInputAppend(tempRightInput)

	newVar := leftRightInputMult(tempLeftInput, tempRightInput, indexOfValueXForLeftInput, indexOfValueXForRightInput)

	return newVar
}

// Opened bracketts from numerator of interpolated polynomials
// Example(normal view): input ((x-1)(x-2))= x^2-x-2x+2 = output (x^2-3x+2)
// Example(code interpritated view): x1-x2-* =
func calcFunc(function string) (result map[string]float64) {
	tempSliceStr := strings.Split(function, " ")
	for i := range tempSliceStr { // Loop interpritate normal "x" view into string [coefficieent of x][x][degree of x]
		if tempSliceStr[i] == "x" {
			tempSliceStr[i] = "1x1"
		}
	}

	var isEnd bool
	tempVal := ""
	var tempValSlice []string
	count := 0 // Uses to calculate the difference between j and i. This value is necessary when we multiply the bracket by the bracket
	n := 2     // It is necessary to take all the new monomials. When you multiply the braket by bracket,
	// you create a number of monomes that you can subtract by the formula (n*2)+1, where n is the number of elements in the 1st bracket
	// Example: (x-1)(x-2), in the first bracket are 3 elements, thats n=3. Open the brackets: x^2-x-2x+2 - here n=3*2+1=7 and so on

	// (x-1)(x-2) = x*x - x -2x + 2. j is responsible for the values in the first bracket, i for the second
	// Firslty, multiplies the elements of the first bracket by the 1st element of the second bracket, then the elements of the 1st bracket ultiplies by the 2nd element of the 2nd bracket
	for j := 1; j < n+1; j++ {
		for i := 1; i < n+1; i++ { // Cycle, which saves multiplier
			tempVal = tempSliceStr[i]
			tempValSlice = append(tempValSlice, multiplyMonoms(tempVal+"*"+tempSliceStr[j+n*2-count], tempSliceStr[i+1], tempSliceStr[j+n*2-count+1]))
		}

		if len(tempValSlice) == n*2 { // If we calculated all of possible values after removing the brackets
			tempSliceStr = append(tempSliceStr[:1], tempSliceStr[1+n:]...) // Delete n elements form slice. This two values corresponding to *- or similar values in infix representation

			for k := 1; k < n*2+1; k++ { // Putting calculated values in slice of infix representation for next calculations
				//Example: (x-1)(x-2)(x-3) -> (x^2-2x-x+2)(x-3) Here we put values form left brackets into slice for next calculating.
				//Next calculating: (x^2-2x-x+2)(x-3) -> x^3-2x^2-x^2+2x-3x^2+6x+3x-6 (Then put this values into slice)

				if len(tempSliceStr) == k { // This condition needed, when our slice of infix representation are full but we have to put more calculated variablse in it

					for i := 0; i < len(tempValSlice)-len(tempSliceStr)+2; i++ { //+2 cause in [0] position situated "y" coordinate value +2
						tempSliceStr = append(tempSliceStr, tempValSlice[k-1])
						k++
					}
					isEnd = true
					break
				}
				tempSliceStr[k] = tempValSlice[k-1]
			}

			if isEnd {
				break
			}

			tempSliceStr[n*2+1] = "*"
			tempValSlice = tempValSlice[:0]
			count = n
			n = n * 2
			j = 0
		}
	}

	result = additionMonoms(tempSliceStr) // Adding cooeficients of same degrees of "x" from interpolated polynomial wich corresponding to one point
	return result
}

// Saves calculated cooeficients of "x" form calculated interpolated polynomials
func saveQAP(partQAP map[string]map[string]float64) {

	if len(partQAP) == 0 { // If "y" values of all coordintes for interpolation are equals zero -> putting zero values in row
		for i := 0; i < len(qapVectorTemp[0]); i++ {
			qapVectorTemp[indexOfRow][i] = 0
		}
	}

	for key1, _ := range partQAP { // Putting one by one calculated values of interpolated polynomials, where each polynomial are interpolated by one r1cs.vector coloumn coordinates
		for key2, value := range partQAP[key1] {

			key2Temp, _ := strconv.Atoi(key2)
			qapVectorTemp[indexOfRow][key2Temp] = value
		}

	}
	indexOfRow++ // Needed for indexing of row, where we have to put calculated values.
}

// 1.Separate numerators and denominators of interpolated Lagrangia polynomials.
// 2.Dividing numerators by corresponding denominators.
// 3.Calculate one general Lagrangia polynomials
func QapCreate(numerator map[string]string, denumerator map[string]string, vectName string) {
	deleteZerosMapValue(numerator, denumerator) // Deleting unusefull interpolated values by point coordinates with y=0 (x,0)

	numeratorEvMap := make(map[string]string) // Postfix to infix view of numerator. Numerator takes from interpolated Lagrangia polynomial
	for key, _ := range numerator {
		numeratorEvMap[key] = r1cs.ParseInfix(numerator[key])
	}

	denumeratorEvMap := make(map[string]float64)
	for key, _ := range numerator { // Postfix to infix view of denominator(then calculate result number of denominator). Denominator takes from interpolated Lagrangia polynomial
		denumeratorEvMap[key] = float64(calculateInfixView(r1cs.ParseInfix(denumerator[key])))
	}

	//fmt.Println(denumeratorValue)
	numeratorCalculatedValue := make(map[string]map[string]float64) // map with keys = degree of "x", values = coefficient of "x"

	for key, value := range numeratorEvMap { // Calculating numerator. Getiing a map where key = degree of "x"; value = coefficient of corresposponding degree of "x"
		numeratorCalculatedValue[key] = calcFunc(value)
	}

	for key1, _ := range numeratorCalculatedValue { // One by one dividing coefficients of "x[degree]", which situated numenator, by denominator (denomiantor is number. always)
		for key2, value := range numeratorCalculatedValue[key1] {
			numeratorCalculatedValue[key1][key2] = value / denumeratorEvMap[key1]
		}
	}

	if len(numeratorCalculatedValue) > 1 { // If in interpolation more then 1 point with y!=0 then summ their interpolated values.
		resultMapOfCoeff := make(map[string]map[string]float64)

		keyForResultMap := ""

		for key, value := range numeratorCalculatedValue { // Copying in result map keys and values of one point interpolated values
			resultMapOfCoeff[key] = value
			keyForResultMap = key
			delete(numeratorCalculatedValue, key)
			break
		}

		// While there are exist values of interpolated points
		for key1, _ := range numeratorCalculatedValue { // Added coefficient of corresponding degrees of "x"
			for key2, value2 := range numeratorCalculatedValue[key1] {
				resultMapOfCoeff[keyForResultMap][key2] = resultMapOfCoeff[keyForResultMap][key2] + value2
			}
			delete(numeratorCalculatedValue, key1)
		}

		saveQAP(resultMapOfCoeff)
	} else {
		saveQAP(numeratorCalculatedValue)
	}

}

// 1.Puts values from r1cs vectors for Lagrangia interpolation
// 2.Save calculated results of QAP vectors
func xyCreate(vectors [][]int, vectName string) {

	var coordin []string
	x := []string{"x"}
	y := []string{"y"}
	for j := 0; j < len(vectors[1]); j++ { // Get coordinates for interpolation
		for i := 0; i < len(vectors); i++ {
			x = append(x, strconv.Itoa(i+1), strconv.Itoa(j+1))
			coordin = append(coordin, strings.Join(x, ""), "=", strconv.Itoa(i+1))
			x = append(x[:0], x[0])

			y = append(y, strconv.Itoa(i+1), strconv.Itoa(j+1))
			coordin = append(coordin, strings.Join(y, ""), "=", strconv.Itoa(vectors[i][j]))
			y = append(y[:0], y[0])

		}
		tempStr := strings.Join(coordin, " ")
		Lagrangia.Start(tempStr)

		QapCreate(Lagrangia.ReturnNumeratorNormal(), Lagrangia.ReturnDenumeratorNormal(), vectName)
		coordin = coordin[:0]
		Lagrangia.ClearAllVar()
	}

	switch vectName { //Save results
	case "A":
		copySlice(qapVectorTemp, qapVectorA)
	case "B":
		copySlice(qapVectorTemp, qapVectorB)
	case "C":
		copySlice(qapVectorTemp, qapVectorC)
	}
}

// Copied vlues from original slice to destination slice. Needed cause native function copy() doesn't work correctly
func copySlice(originalSlice [][]float64, destinationSlice [][]float64) {
	for j := 0; j < len(originalSlice); j++ {
		for i := 0; i < len(originalSlice[0]); i++ {
			destinationSlice[j][i] = originalSlice[j][i]
		}
	}
}

// function wich return QAP vectors to others packages
func QAPVectAReturn() [][]float64 {
	return qapVectorA
}

func QAPVectBReturn() [][]float64 {
	return qapVectorB
}

func QAPVectCReturn() [][]float64 {
	return qapVectorC
}

// Function for startiing calculating
func main() {
	r1cs.Start("x ^ 3 + x + 5", "x = 3 y = 35")

	// Allocate memory for QAP vectors
	vectorQAPSizeAllocate(r1cs.ReturnVectorsA(), "A")
	vectorQAPSizeAllocate(r1cs.ReturnVectorsB(), "B")
	vectorQAPSizeAllocate(r1cs.ReturnVectorsC(), "C")
	vectorQAPSizeAllocate(r1cs.ReturnVectorsC(), "temp")

	// ПРОБЛЕМА С КОПИРОВАНИЕМ СЛАЙСА В ФИНАЛЬНЫЙ ВЕКТОР
	// start point of QAP calculations
	xyCreate(r1cs.ReturnVectorsA(), "A")
	indexOfRow = 0
	xyCreate(r1cs.ReturnVectorsB(), "B")
	indexOfRow = 0
	xyCreate(r1cs.ReturnVectorsC(), "C")

	// Print results
	fmt.Println("QAP reprresentation of vector A:")
	fmt.Println(qapVectorA)
	fmt.Println("QAP reprresentation of vector B:")
	fmt.Println(qapVectorB)
	fmt.Println("QAP reprresentation of vector C:")
	fmt.Println(qapVectorC)

}
