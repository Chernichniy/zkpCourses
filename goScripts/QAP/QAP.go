// This package calculate QAP polynomial for input function. Input function in form like R1CS package.
// Don't forget about module!

package qap

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Chernichniy/zkpCourses/goScripts/Lagrangia"
	r1cs "github.com/Chernichniy/zkpCourses/goScripts/R1CS"
)

var module int

var qapVectorA [][]int
var qapVectorB [][]int
var qapVectorC [][]int
var qapVectorFinall []int
var qapVectorZ []int
var qapVectorTemp [][]int
var quotientOfFruct []int

var indexOfRow int = 0
var isQAPCorrect bool

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
func vectorQAPSizeAllocate(r1csVecor [][]int, vectName string) [][]int {

	var temp = make([][]int, len(r1csVecor[0]))

	switch { // Needed when we going to interpolate polynomials by only one point, cause numerator
	// of those interpolated polynoail are consost of two values: value "x" and value "number"
	case len(r1csVecor) == 1:
		for i := range temp {
			temp[i] = make([]int, 2)
		}
	default:
		for i := range temp {
			temp[i] = make([]int, len(r1csVecor))
		}
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

	return temp

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
					str[i+2] = fmt.Sprintf("%f", num)
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
					if !isFinite(num) {
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
	var coeffLeftForLeftInput int = 1
	var coeffRightForLeftInput int = 0

	var coeffLeftForRightInput int = 1
	var coeffRightForRightInput int = 0

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

	// String -> int: Assign coefficients and degrees values of left and right multipliers
	switch {
	case len(tempLeftCoefForLeft) > 0:
		coeffLeftForLeftInput, _ = strconv.Atoi(tempLeftCoefForLeft)
	}

	switch {
	case len(tempRightCoefForLeft) > 0:
		coeffRightForLeftInput, _ = strconv.Atoi(tempRightCoefForLeft)

	}

	switch {
	case len(tempLeftCoefForRight) > 0:
		coeffLeftForRightInput, _ = strconv.Atoi(tempLeftCoefForRight)
	}

	switch {
	case len(tempRightCoefForRight) > 0:
		coeffRightForRightInput, _ = strconv.Atoi(tempRightCoefForRight)
	}

	return strconv.Itoa(coeffLeftForLeftInput*coeffLeftForRightInput) + "x" + strconv.Itoa(coeffRightForLeftInput+coeffRightForRightInput)

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
			leftInputStr := strings.Join(tempValLeftSide, "")
			input[i-1] = leftInputStr
			input = append(input[:0], input[i-1:]...)
			i = -1

		} else if input[i] == "x" && i+2 != len(input) || input[0] == "x" && i+2 != len(input) {

			for j := i + 1; j < len(input); j++ {
				tempValRightSide = append(tempValRightSide, input[j])
			}
			leftInputStr := strings.Join(tempValRightSide, "")
			input[i+1] = leftInputStr
			input = append(input[:i+2], input[len(input):]...)

		}
	}
	for i := range input { // Loop for index finding
		if input[i] == "x" {
			index = i
		}
	}

	return input, index //(x-1)(x-2) = x-3x+2;  x 1 - x 2 - *;  [1][x][1] [1][x][0]
}

// 1. Multiplies calculated coeeficint of monomial(after brackets opening) by "y" coordinate
// 2. Sums coefficient of corresponding "x" degrees. In other words, simplifies the expression.
// Example: input (x^3-3x^2+2x-3x^2+9x-6) -> output (x^3-6x^2+11x-6)
func additionMonoms(monoms []string) map[string]int {
	copiedMonoms := make([]string, len(monoms))
	copy(copiedMonoms, monoms)

	var mapOfCoeff = make(map[string]int)

	for j := 1; j < len(copiedMonoms); j++ { // Devided cooeficients of variables by degree of "x", then summing coefficients by corresponding degrees.
		// If there no x variable in value, the value writing like this: [coefficient][x][0]
		// Example: value "6" writes like 6x0
		tempMonomSlice := strings.Split(copiedMonoms[j], "x")
		val, isExist := mapOfCoeff[tempMonomSlice[1]]

		if isExist { // If coefficient by degree "x" alredy exist
			tempValueFromMap, _ := strconv.Atoi(tempMonomSlice[0])
			mapOfCoeff[tempMonomSlice[1]] = val + tempValueFromMap
		} else {
			mapOfCoeff[tempMonomSlice[1]], _ = strconv.Atoi(tempMonomSlice[0])
		}
	}

	yValue, _ := strconv.Atoi(copiedMonoms[0])
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
func calcFunc(function string) (result map[string]int) {
	tempSliceStr := strings.Split(function, " ")
	var mapForTwoOrLessPointsInterp = make(map[string]int)

	if len(tempSliceStr) < 6 { // If polynomial interpolated only from 2 (or less) points
		mapForTwoOrLessPointsInterp["1"] = 1
		mapForTwoOrLessPointsInterp["0"], _ = strconv.Atoi(tempSliceStr[2])

		for key, value := range mapForTwoOrLessPointsInterp {
			if key == "0" {
				value = value * -1
			}
			yCoordinateInt, _ := strconv.Atoi(tempSliceStr[0])
			mapForTwoOrLessPointsInterp[key] = numberByModule(value * yCoordinateInt)
		}
		result = mapForTwoOrLessPointsInterp
		return result
	}

	for i := range tempSliceStr { // Loop interpritate normal "x" view into string [coefficieent of x][x][degree of x]
		if tempSliceStr[i] == "x" {
			tempSliceStr[i] = "1x1"
		}
	}

	//var isEnd bool
	tempVal := ""
	var tempValSlice []string
	//countMultiplies := 0
	count := 0 // Uses to calculate the difference between j and i. This value is necessary when we multiply the bracket by the bracket
	n := 2     // It is necessary to take all the new monomials. When you multiply the braket by bracket,
	// you create a number of monomes that you can subtract by the formula (n*2)+1, where n is the number of elements in the 1st bracket
	// Example: (x-1)(x-2), in the first bracket are 3 elements, thats n=3. Open the brackets: x^2-x-2x+2 - here n=3*2+1=7 and so on

	for j := 1; j < n+1; j++ {
		for i := 1; i < n+1; i++ { // Cycle, which saves multiplier
			tempVal = tempSliceStr[i]
			tempValSlice = append(tempValSlice, multiplyMonoms(tempVal+"*"+tempSliceStr[j+n*2-count], tempSliceStr[i+1], tempSliceStr[j+n*2+1-count]))
		}
	}
	var resultsTempSlice = make([]string, len(tempValSlice))
	copy(resultsTempSlice, tempValSlice)
	//resultsTempSlice := tempValSlice
	tempValSlice = tempValSlice[:0]
	// (x-1)(x-2) = x*x - x -2x + 2. j is responsible for the values in the first bracket, i for the second
	// Firslty, multiplies the elements of the first bracket by the 1st element of the second bracket, then the elements of the 1st bracket ultiplies by the 2nd element of the 2nd bracket
	for j := 9; j < len(tempSliceStr); j++ {
		for i := 0; i < len(resultsTempSlice); i++ { // Cycle, which saves multiplier
			tempVal = resultsTempSlice[i]
			tempValSlice = append(tempValSlice, multiplyMonoms(tempVal+"*"+tempSliceStr[j], "", tempSliceStr[j+1]))
		}

		if len(tempValSlice) == len(resultsTempSlice)*2 {
			resultsTempSlice = make([]string, len(tempValSlice))
			//resultsTempSlice = append(resultsTempSlice[:1], resultsTempSlice[:]...)
			copy(resultsTempSlice, tempValSlice)
			tempValSlice = tempValSlice[:0]
		}

		if j+3 == len(tempSliceStr) {
			break
		} else if j%2 == 0 {
			j = j + 2
		}

	}

	resultsTempSlice = append(resultsTempSlice[:1], resultsTempSlice[0:]...)
	resultsTempSlice[0] = tempSliceStr[0]

	result = additionMonoms(resultsTempSlice) // Adding cooeficients of same degrees of "x" from interpolated polynomial wich corresponding to one point

	for key, value := range result {
		result[key] = numberByModule(value)
	}
	return result
}

// Saves calculated cooeficients of "x" form calculated interpolated polynomials
func saveQAP(partQAP map[string]map[string]int) {

	if len(partQAP) == 0 { // If "y" values of all coordinates for interpolation are equals zero -> putting zero values in row
		for i := 0; i < len(qapVectorTemp[0]); i++ {
			qapVectorTemp[indexOfRow][i] = 0
		}
	}

	_, isExist := partQAP["polZ"]

	if isExist {
		qapVectorZ = make([]int, len(partQAP["polZ"]))
		for key1, _ := range partQAP { // Putting one by one calculated coefficient values polynomial Z

			for key2, value := range partQAP[key1] {

				key2Temp, _ := strconv.Atoi(key2)
				qapVectorZ[key2Temp] = value
			}

		}
		return
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

	denumeratorEvMap := make(map[string]int)
	for key, _ := range numerator { // Postfix to infix view of denominator(then calculate result number of denominator). Denominator takes from interpolated Lagrangia polynomial
		denumeratorEvMap[key] = calculateInfixView(r1cs.ParseInfix(denumerator[key]))
	}

	//fmt.Println(denumeratorValue)
	numeratorCalculatedValue := make(map[string]map[string]int) // map with keys = degree of "x", values = coefficient of "x"

	for key, value := range numeratorEvMap { // Calculating numerator. Getiing a map where key = degree of "x"; value = coefficient of corresposponding degree of "x"
		numeratorCalculatedValue[key] = calcFunc(value)
	}

	for key1, _ := range numeratorCalculatedValue { // One by one dividing coefficients of "x[degree]", which situated numenator, by denominator (denomiantor is number. always)
		for key2, value := range numeratorCalculatedValue[key1] {
			numeratorCalculatedValue[key1][key2] = fructionsByModule(value, denumeratorEvMap[key1]) /*int(value)  int(denumeratorEvMap[key1])) */ /*value / denumeratorEvMap[key1]*/
		}
	}

	if len(numeratorCalculatedValue) > 1 { // If in interpolation more then 1 point with y!=0 then summ their interpolated values.
		resultMapOfCoeff := make(map[string]map[string]int)

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
				resultMapOfCoeff[keyForResultMap][key2] = (resultMapOfCoeff[keyForResultMap][key2] + value2) % module
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
	for j := 0; j < len(vectors[0]); j++ { // Get coordinates for interpolation
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

		if len(vectors) == 1 { // If interpolate polynomial by one point
			polynomialByOnePoint := Lagrangia.ReturnPolynomialByOnePoint()
			tempPol := strings.Split(polynomialByOnePoint["pol0"], " / ")

			tempMapNumerator := make(map[string]string)
			tempMapDenumerator := make(map[string]string)

			tempMapNumerator["pol0"] = tempPol[0]
			tempMapDenumerator["pol0"] = tempPol[1]
			QapCreate(tempMapNumerator, tempMapDenumerator, "pol0")
		} else {
			QapCreate(Lagrangia.ReturnNumeratorNormal(), Lagrangia.ReturnDenumeratorNormal(), vectName)
		}

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
func copySlice(originalSlice [][]int, destinationSlice [][]int) {
	for j := 0; j < len(originalSlice); j++ {
		for i := 0; i < len(originalSlice[0]); i++ {
			destinationSlice[j][i] = originalSlice[j][i]
		}
	}
}

// Doing division by module
// Example: 3/4 mod 15 -> while 3+module/4 without reminder -> 48/4 -> result number is 48/4 = 12 => 3/4 mod 15 = 12
func fructionsByModule(numerator int, denominator int) int {

	for numerator%denominator != 0 {
		numerator = numerator + module
	}

	result := numerator / denominator
	for result < 0 {
		result = result + module
	}

	return result % module

}

// Convert number by module
// Examples: 4 mod 3 = 1; -4 mod 3 = 2
func numberByModule(number int) int {

	switch {
	case number > 0:
		return number % module
	case number < 0:
		for number < 0 {
			number = number + module
		}
		return number % module
	}
	return number
}

// Multiply one value of witness vector on corresponding row of QAP representations of matrices A,B,C
// Example: w1*A1; w2*A2; w3*A3 ... wn*An, n - number of row
func witnessMultOnQAP(vector [][]int) [][]int {
	witness := r1cs.ReturnWitnessNumbers()
	resultVector := vectorQAPSizeAllocate(r1cs.ReturnVectorsA(), "")

	for j := 0; j < len(vector); j++ {
		for i := 0; i < len(vector[0]); i++ {
			resultVector[j][i] = (witness[j] * vector[j][i]) % module
		}
	}
	return resultVector
}

// Addition rows of QAP representation of R1CS matrices after multiplying their rows on witness vector
// Eexample: A1+A2+A3+...+An, n - number of rows
func summOfQAPParts(vector [][]int) []int {
	vectorAfterWitnessMult := witnessMultOnQAP(vector)

	resultVector := vectorAfterWitnessMult[0][:]

	for j := 1; j < len(vectorAfterWitnessMult); j++ {
		for i := 0; i < len(resultVector); i++ {
			resultVector[i] = (resultVector[i] + vectorAfterWitnessMult[j][i]) % module
		}
	}

	return resultVector
}

// Multiply QAP representation of matrices A and B (A*B)
func qapVectorsMult(vectorA []int, vectorB []int) []int {

	var expressionStrForCals []string

	for i := 0; i < len(vectorA); i++ {
		expressionStrForCals = append(expressionStrForCals, strconv.Itoa(vectorA[i])+"x"+strconv.Itoa(i))
	}
	expressionStrForCals = append(expressionStrForCals, "*", "*") // Standartization of multiplyMonoms func input

	for i := 0; i < len(vectorB); i++ {
		expressionStrForCals = append(expressionStrForCals, strconv.Itoa(vectorB[i])+"x"+strconv.Itoa(i))
	}

	tempVal := ""
	var tempValSlice []string

	for j := len(vectorA) + 2; j < len(expressionStrForCals); j++ {
		for i := 0; i < len(vectorA); i++ { // Cycle, wich saves result of multiplying
			tempVal = expressionStrForCals[i]
			tempValSlice = append(tempValSlice, multiplyMonoms(tempVal+"*"+expressionStrForCals[j], "", ""))
		}

	}

	tempValSlice = append(tempValSlice[:1], tempValSlice[0:]...)
	tempValSlice[0] = "1"
	mapOfCoefficientByDegree := additionMonoms(tempValSlice)

	var result = make([]int, len(mapOfCoefficientByDegree))

	for key, value := range mapOfCoefficientByDegree {
		index, _ := strconv.Atoi(key)
		result[index] = value % module
	}
	return result

}

// Calculate vanish (Z) polynomial
func polynomialZCreate() {
	constraints := r1cs.ReturnWitnessFormal()
	fmt.Println(constraints)
	var index int

	for i := range constraints {
		if constraints[i] == "Con1" {
			index = i
			break
		}
	}

	constraintsNumber := len(constraints) - index
	var zPolynomialSlice []string

	zPolynomialSlice = append(zPolynomialSlice, "( 1 * ")

	for i := 1; i < constraintsNumber+1; i++ {
		if i == constraintsNumber {
			zPolynomialSlice = append(zPolynomialSlice, "( x - "+strconv.Itoa(i)+" ) ) ")
			break
		}
		zPolynomialSlice = append(zPolynomialSlice, "( x - "+strconv.Itoa(i)+" ) * ")
	}

	var mapOfZPolynomialNumerator = make(map[string]string)
	var mapOfZPolynomialDenumerator = make(map[string]string)

	mapOfZPolynomialNumerator["polZ"] = strings.Join(zPolynomialSlice, "") // Needed for QapCreate func call
	mapOfZPolynomialDenumerator["polZ"] = "1"                              // Same as previous

	QapCreate(mapOfZPolynomialNumerator, mapOfZPolynomialDenumerator, "Z")
	fmt.Println(qapVectorZ)

}

// Devided QAP polynomial on vanish (Z) polynomial
// Input numerator = QAP polynomial coefficient and degrees
// Input denominator = vanish polynomial coefficient and degrees
func polynomialsDevide(numerator []int, denominator []int) (bool, []int) {

	var tempDenominator = make([]int, len(denominator))
	var tempNumerator = make([]int, len(numerator))
	copy(tempDenominator, denominator)
	copy(tempNumerator, numerator)

	var coefficientForMatchModule int // 5x^4+2/x (mod 11) -> 6x^4+2/6x

	var result []string

	for len(tempNumerator) != 0 {

		for i := len(tempNumerator) - 1; i >= 0; i-- {
			if tempNumerator[i] == 0 {
				continue
			}
			coefficientForMatchModule = -1 * (tempNumerator[i] - module)
			break
		}

		for i := range tempDenominator { //эта штука умножает знаменатель на число, что бы потом отнять по модулю знаменатель

			degreeDifferencies := (len(tempNumerator) - 1) - (len(tempDenominator) - 1)
			strForLeftInput := strconv.Itoa(denominator[i]) + "x" + strconv.Itoa(i)
			strForRightInput := strconv.Itoa(tempDenominator[len(tempDenominator)-1]*coefficientForMatchModule) + "x" + strconv.Itoa(degreeDifferencies)

			result = append(result, multiplyMonoms(strForLeftInput+"*"+strForRightInput, "", ""))

			fmt.Println(result)
		}

		for i := range result {
			strTransitTempValueOfResult := strings.Split(result[i], "x")
			indexForNumerator, _ := strconv.Atoi(strTransitTempValueOfResult[1])
			valueForIndexNumerator, _ := strconv.Atoi(strTransitTempValueOfResult[0])
			tempNumerator[indexForNumerator] = numberByModule(tempNumerator[indexForNumerator] + valueForIndexNumerator)
		}

		for i := len(tempNumerator) - 1; i >= 0; i-- { // Delete all zero`s coefficients
			if tempNumerator[i] == 0 {
				tempNumerator = tempNumerator[:i]
			} else {
				break
			}
		}

		if len(tempNumerator) < len(tempDenominator) && len(tempNumerator) != 0 {
			return false, tempNumerator
		}

		result = result[:0]
	}
	return true, tempNumerator
}

// Created QAP polynomial
func fullQAPPolynomialCalc() {
	vectorA := summOfQAPParts(qapVectorA)
	vectorB := summOfQAPParts(qapVectorB)
	vectorC := summOfQAPParts(qapVectorC)

	vectorAMultB := qapVectorsMult(vectorA, vectorB)

	resultVector := make([]int, len(vectorAMultB))
	qapVectorFinall = make([]int, len(resultVector)) // Save variable

	for i := 0; i < len(vectorA); i++ {
		resultVector[i] = numberByModule(vectorAMultB[i] - vectorC[i])
	}

	for i := len(vectorA); i < len(vectorAMultB); i++ {
		resultVector[i] = vectorAMultB[i]
	}

	copy(qapVectorFinall, resultVector)
	polynomialZCreate()

	isValid, quotient := polynomialsDevide(resultVector, qapVectorZ)
	isQAPCorrect = isValid

	if isValid {
		fmt.Println("Full QAP vector coefficients:\n", resultVector)
		fmt.Println("Vanish vector (Z polynomial) coefficients:\n", qapVectorZ)
		fmt.Println("QAP representation is correct\n")
	} else {
		fmt.Println("Full QAP vector coefficients:\n", resultVector)
		fmt.Println("Vanish vector (Z polynomial) coefficients:\n", qapVectorZ)
		fmt.Println("QAP representation isn`t correct\n Quotient = ")
		fmt.Println(quotient)
		quotientOfFruct = make([]int, len(quotient))
		copy(quotientOfFruct, quotient) // Saving quotient
	}

}

// function wich return QAP vectors to others packages
func QAPVectAReturn() [][]int {
	return qapVectorA
}

func QAPVectBReturn() [][]int {
	return qapVectorB
}

func QAPVectCReturn() [][]int {
	return qapVectorC
}

func QAPVectFinallReturn() []int {
	return qapVectorFinall
}

func QAPQuiotientOfFinallFructReturn() []int {
	return quotientOfFruct
}

func QAPVanishVectorReturn() []int {
	return qapVectorZ
}

func QAPCorrectCalcOrNot() bool {
	return isQAPCorrect
}

func ClearAllVar() {
	qapVectorA = nil
	qapVectorB = nil
	qapVectorC = nil
	qapVectorZ = nil

	qapVectorFinall = nil
	qapVectorTemp = nil
	quotientOfFruct = nil

	indexOfRow = 0
	isQAPCorrect = false

	Lagrangia.ClearAllVar()
	r1cs.ClearAllVar()
}

/*
func main() {
	Start()
	ClearAllVar()
	Start()
	Start()
}
*/

// Function for starting calculating
func Start(function string, roots string, mod int) {

	//Lagrangia.ClearAllVar()
	//r1cs.ClearAllVar()

	//module = 11
	//r1cs.Start("x ^ 3 + x + 5", "x = 2 y = 15")
	module = mod
	r1cs.Start(function, roots)
	//r1cs.Start("( x + z ) ^ 2 + z + 1", "x = 1 z = 2 y = 12")
	//r1cs.Start("x ^ ( g ^ 2 + 1 ) + 7", "x = 2 g = 2 y = 39")
	//r1cs.Start("x ^ ( 2 + z ) * g", "x = 2 z = 1 g = 2 y = 16")

	// Allocate memory for QAP vectors
	vectorQAPSizeAllocate(r1cs.ReturnVectorsA(), "A")
	vectorQAPSizeAllocate(r1cs.ReturnVectorsB(), "B")
	vectorQAPSizeAllocate(r1cs.ReturnVectorsC(), "C")
	vectorQAPSizeAllocate(r1cs.ReturnVectorsC(), "temp")

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

	fullQAPPolynomialCalc()

}
