// This package calculate R1CS matricies for input (string) function and roots of function
// Input string have to being in form, where first variable without space before and last variable without space after
// ohers variables and operators need to put before and after spaces
// Example: wrong input - " x+3^x+7"; correct input - "x + 3 ^ x + 7" or "( 3 + x ) ^ ( h + 2 )"

// Output of this package is R1CS matricies for left, right inputs and output; witness in numbers and formal forms
// Matricies in form of [][]int slice; witness in form of []int or []string slices (number and formal form)

package r1cs

import (
	"fmt"
	"strconv"
	"strings"
)

var input string
var roots_input string

// var input = "x ^ 3 + x + 5"
// var roots_input = "x = 3 y = 35" //эти две переменные мы получаем с сайта
//var roots_input = "x = 1 z = 2" //эти две переменные мы получаем с сайта

var mapOfRoots = make(map[string]int) // Store roots of our function and y value in form: key = value
var evaluationInput string            // String of function with replaced characters variables by numbers of private and public inputs

var constraints []string        // Store constraints in form of numbers (2*2=4 4+2=6 ...)
var constraintsFormall []string // Store constraints in form of characters (x*x=Con1 Con1+2=Con2 ...)

var witnes []int          // Store witness in form of numbers
var witnesFormal []string // Store witness in form of characters

var vectorsA [][]int // Store R1CS representation of left input (matrix)
var vectorsB [][]int // Store R1CS representation of right input (matrix)
var vectorsC [][]int // Store R1CS representation of output (matrix)

var witnesFormalChecker []string // Store witness in form of characters but in (leftInput*rightInput-Output). Uses for checking correctness of R1CS matricis
var witnesChecker []string       // Store witness in form of numbers but in (2*2-4). Uses for checking correctness of R1CS matricis

var indexOfYInWitness int // Store index of witness vector, where "y" value situated

// Struct of operators
var oper = map[string]struct {
	prec   int
	rAssoc bool
}{
	"^": {4, true},
	"*": {3, false},
	"/": {3, false},
	"+": {2, false},
	"-": {2, false},
}

// Convert function from infix view to postfix view (Reverse Polish Notation)
// On input: 3^2+5*x; Output: 32^5x*+
func ParseInfix(e string) (rpn string) {
	var stack []string // holds operators and left parenthesis
	for _, tok := range strings.Fields(e) {
		switch tok {
		case "(":
			stack = append(stack, tok) // push "(" to stack
		case ")":
			var op string
			for {
				// pop item ("(" or operator) from stack
				op, stack = stack[len(stack)-1], stack[:len(stack)-1]
				if op == "(" {
					break // discard "("
				}
				rpn += " " + op // add operator to result
			}
		default:
			if o1, isOp := oper[tok]; isOp {
				// token is an operator
				for len(stack) > 0 {
					// consider top item on stack
					op := stack[len(stack)-1]
					if o2, isOp := oper[op]; !isOp || o1.prec > o2.prec ||
						o1.prec == o2.prec && o1.rAssoc {
						break
					}
					// top item is an operator that needs to come off
					stack = stack[:len(stack)-1] // pop it
					rpn += " " + op              // add it to result
				}
				// push operator (the new one) to stack
				stack = append(stack, tok)
			} else { // token is an operand
				if rpn > "" {
					rpn += " "
				}
				rpn += tok // add operand to result
			}
		}
	}
	// drain stack to result
	for len(stack) > 0 {
		rpn += " " + stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
	return
}

// 1. Evaluate function wich given in postfix view
// 2. Calculate and saves the constraint with their "way" in number form
func constraintsEval(rpn string) (res int) {

	str := strings.Split(rpn, " ")

	for i := 0; i < len(str); i++ {
		if i+2 <= len(str)-1 {
			num1, _ := strconv.Atoi(str[i])
			num2, _ := strconv.Atoi(str[i+1])
			num3, _ := strconv.Atoi(str[i+2]) // If 3rd value equal to 0 then 3rd value is a operator
			if num3 == 0 {

				switch str[i+2] {
				case "+":

					num := num1 + num2                                                    // Calculate output of constraint
					constraintsSaver(str[i], str[i+2], str[i+1], strconv.Itoa(num), true) // Call func for saving constraint

					str[i+2] = strconv.Itoa(int(num)) // Replace operator by output of constraint

					str = append(str[:i], str[i+2:]...) // Delete inputs of constraints
					i = -1                              // Starting from begining of function in postfix view
				case "*":
					num := num1 * num2 // Calculate output of constraint
					constraintsSaver(str[i], str[i+2], str[i+1], strconv.Itoa(num), true)
					str[i+2] = strconv.Itoa(int(num)) // Replace operator by output of constraint

					str = append(str[:i], str[i+2:]...) // Delete inputs of constraints
					i = -1                              // Starting from begining of function in postfix view
				case "^":
					numTemp := num1 // Raising to degree = multiply numbers few times by themselves
					num := num1

					for j := num2; j > 1; j-- { // Num2 - degree, so we multiply "x" on "x"  Num2 times
						num = numTemp * num1

						constraintsSaver(strconv.Itoa(int(numTemp)), "*", str[i], strconv.Itoa(int(num)), true) // Call func for saving constraint for each multiplying

						numTemp = num
					}

					str[i+2] = strconv.Itoa(int(num)) // Replace operator by output of constraint

					str = append(str[:i], str[i+2:]...) // Delete inputs of constraints
					i = -1                              // Starting from begining of function in postfix view

					/*case "-": Not implemented yet

					num := num1 - num2

					str[i+2] = strconv.Itoa(int(num))

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
								//str[i+2] = strconv.FormatFloat(num, 'E', -1, 64)
								str[i+2] = fmt.Sprintf("%f", num)

							}
							str = append(str[:i], str[i+2:]...)
							i = -1*/
				}

			}
		}
	}
	res, _ = strconv.Atoi(str[0])
	return res
	// }
}

// Save constraint with their "way"
// Input true: save constraints in number view; Input false: save constrints in formal view
func constraintsSaver(lftInput string, operation string, rghtInput string, output string, types bool) {
	switch types {
	case true:
		if operation == "^" { // Needed cause in R1CS only addition or multiply are can be
			rghtInput = lftInput
			operation = "*"
		}
		constraints = append(constraints, lftInput) // Insert one by one
		constraints = append(constraints, operation)

		constraints = append(constraints, rghtInput)
		constraints = append(constraints, output)

	case false:
		if operation == "^" { // Needed cause in R1CS only addition or multiply are can be
			rghtInput = lftInput
			operation = "*"
		}
		constraintsFormall = append(constraintsFormall, lftInput) // Insert one by one
		constraintsFormall = append(constraintsFormall, operation)

		constraintsFormall = append(constraintsFormall, rghtInput)
		constraintsFormall = append(constraintsFormall, output)
	}

}

// Calculate and saves the constraint with their "way" in formal form
func constraintsFormalForm(rpn string) {
	str := strings.Split(rpn, " ")
	a := oper   // Operator checking
	numCon := 1 // Counter of constraint number
	var constrNum = []string{"Con", ""}

	for i := 0; i < len(str); i++ {

		for key := range a {
			constrNum[1] = strconv.Itoa(numCon)

			if str[i] == key && key == "^" {
				first := str[i-2]
				second := first
				num := 0

				isConstraint := strings.Split(str[i-1], "n") // Need it for check if there are constraint or a variable in degree
				_, err := strconv.Atoi(str[i-1])
				switch {
				case isConstraint[0] == "Co":
					num = evalOneConstraint(str[i-1])
				case err != nil:
					num = mapOfRoots[str[i-1]]
				default:
					num, _ = strconv.Atoi(str[i-1])
				}

				for j := num; j > 1; j-- {
					constraintsSaver(first, "*", second, strings.Join(constrNum, ""), false)
					numCon = numCon + 1
					first = strings.Join(constrNum, "")
					constrNum[1] = strconv.Itoa(numCon)
				}

				numCon = numCon - 1
				constrNum[1] = strconv.Itoa(numCon)

				str[i] = strings.Join(constrNum, "")
				str = append(str[:i-2], str[i:]...)
				numCon = numCon + 1
				i = 0

			} else if str[i] == key {
				constrNum[1] = strconv.Itoa(numCon)
				constraintsSaver(str[i-2], str[i], str[i-1], strings.Join(constrNum, ""), false)
				numCon = numCon + 1
				str[i] = strings.Join(constrNum, "")
				str = append(str[:i-2], str[i:]...)
				i = 0
				break
			}
		}

	}

}

// Get private and public inputs form general input
func rootsMap(roots string) {
	strRoots := strings.Split(roots, " ")

	for i := 0; i < len(strRoots); i++ {
		if strRoots[i] == "=" {
			var temp, _ = strconv.Atoi(strRoots[i+1])
			mapOfRoots[strRoots[i-1]] = temp // Placed roots in map
		}
	}
	fmt.Println(mapOfRoots)

}

// Calculated only one constraint.
// Input: Constraint's number = Con1, Con2... Oputput: con1=2+3=5 - 5 is output
func evalOneConstraint(constraintID string) int {
	var constraintFunctionTemp []string
	for i := 3; i < len(constraintsFormall); i++ { // Find constraint function
		if constraintsFormall[i] == constraintID {
			constraintFunctionTemp = append(constraintFunctionTemp, constraintsFormall[i-3], constraintsFormall[i-2], constraintsFormall[i-1]) // Copy constraint function
		}

	}
	constraintFunction := evalInput(strings.Join(constraintFunctionTemp, " "), true)

	constraintFunctionInfix := ParseInfix(constraintFunction)
	res := constraintsEval(constraintFunctionInfix)
	return res
}

// Replaced characters in input function (infix view) by corresponding values from mapOfRoots
func evalInput(function string, tag bool) string {
	strFunc := strings.Split(function, " ")
	for i := 0; i < len(strFunc); i++ {
		_, ok := mapOfRoots[strFunc[i]]
		if ok == true {
			strFunc[i] = strconv.Itoa(mapOfRoots[strFunc[i]])
		}
	}
	fmt.Println(strFunc)

	evaluationInput = strings.Join(strFunc, " ")
	return evaluationInput
}

// Generating witnes in formal and number views
func witnesInit() {

	witnes = append(witnes, 1)                           // Put 1 on first place in vector
	witnesFormal = append(witnesFormal, strconv.Itoa(1)) // Put in vector public and private inputs
	for key, val := range mapOfRoots {

		witnesFormal = append(witnesFormal, key)
		witnes = append(witnes, val)
	}

	for i := 1; i < len(witnesFormal); i++ { // Find index of "y" in witness vector
		if witnesFormal[i] == "y" {
			indexOfYInWitness = i
			break
		}
	}

	witnessAdd()    // Call func for add constraints variables in numbers view
	formalWitness() // Call func for add constraints variables in formal view
}

// Create witness in formal view
func formalWitness() {
	for i := 0; i < len(constraintsFormall); i++ { // Each 4th elements are equal to output of one corresponding constraint
		if (i+1)%4 == 0 {
			//temp, _ := strconv.Atoi(constraints[i])
			witnesFormal = append(witnesFormal, constraintsFormall[i])
		}
	}

	//witnesFormal[len(witnesFormal)-1]=
}

// Added constraints outputs in witness vector of numbers form
func witnessAdd() {
	for i := 0; i < len(constraints); i++ {
		if (i+1)%4 == 0 { // Each 4th elements are equal to output of one corresponding constraint
			temp, _ := strconv.Atoi(constraints[i])
			witnes = append(witnes, temp)
		}
	}

}

// Allocate memory for zero fulled two demensional array wich are matrix of R1CS representation for inputs and output
func ZeroOneVectorFulling() (zeroOneVector [][]int) {
	zeroOneVector = make([][]int, len(constraintsFormall)/4)
	for i := range zeroOneVector {
		zeroOneVector[i] = make([]int, len(witnesFormal)-1)
	}
	return zeroOneVector
}

// Determine wich func to call. Func for inputs or output R1CS
func operatorsPipe(operator string) {
	switch operator {
	case "a":
		r1CSCompilerOperatorA()
	case "b":
		r1CSCompilerOperatorB()
	case "c":
		r1CSCompilerOperatorC()
	}
}

// Filled two demensional array for output C of R1CS
func r1CSCompilerOperatorC() {

	r1csVector := ZeroOneVectorFulling()

	counter := 0 // Needed for rows number controlling

	for i := 3; i < len(constraintsFormall); i++ { // On each 3rd index the output is situated
		for j := 0; j < len(witnesFormal); j++ {
			if i == len(constraintsFormall)-1 && constraintsFormall[i] == witnesFormal[j] { // If index of witness vector last then put 1 in "y" index
				r1csVector[counter][indexOfYInWitness] = 1 // output of function  -  "y" value
				break

			} else if constraintsFormall[i] == witnesFormal[j] { // Put in index, where witness variable situated, value 1 (if statement true)
				r1csVector[counter][j] = 1
				break
			}

		}
		i = i + 3             // Next output
		counter = counter + 1 // Next row

	}
	vectorsC = r1csVector

}

// Filled two demensional array for right input B of R1CS
func r1CSCompilerOperatorB() {

	r1csVector := ZeroOneVectorFulling()
	counter := 0 // Needed for rows number controlling

	for i := 2; i < len(constraintsFormall); i++ { // On each 3nd index the output is situated
		for j := 0; j < len(witnesFormal); j++ {
			if constraintsFormall[i] == witnesFormal[j] && constraintsFormall[i-1] != "+" { // Put in index, where witness variable situated, value 1 (if statement true)
				r1csVector[counter][j] = 1
				break

			} else if constraintsFormall[i-1] == "+" { // If there "+" in constraint then we put 1 value in 0 index, cause constraints looks like this: (x+2)*1 (R1CS spec), so value 1 - right input
				r1csVector[counter][0] = 1
				break
			} else if j == len(witnesFormal)-1 { // If in witness veactor not stored the same value, then this value is number
				r1csVector[counter][0], _ = strconv.Atoi(constraintsFormall[i])
				break
			}

		}
		i = i + 3             // Next output
		counter = counter + 1 // Next row

	}
	vectorsB = r1csVector

}

// Filled two demensional array for left input A of R1CS
func r1CSCompilerOperatorA() {

	r1csVector := ZeroOneVectorFulling()
	//fmt.Println(r1csVector)
	counter := 0 // Needed for rows number controlling

	for i := 0; i < len(constraintsFormall); i++ {
		for j := 0; j < len(witnesFormal); j++ {
			if constraintsFormall[i] == witnesFormal[j] && constraintsFormall[i+1] != "+" { // Put in index, where witness variable situated, value 1 (if statement true)
				r1csVector[counter][j] = 1
				break

			} else if constraintsFormall[i] == witnesFormal[j] && constraintsFormall[i+1] == "+" { // If constraints describe addtion, then this constraint could be (x+2)*1 where x+2 have to store in left input matrix

				r1csVector[counter][j] = 1
				for k := 0; k < len(witnesFormal); k++ { // Check if there are knownable value in witness vector
					if witnesFormal[k] == constraintsFormall[i+2] { // If yes put 1 in corresponding index
						r1csVector[counter][k] = 1
						break
					} else if k == len(witnesFormal)-1 { // If no then it's a number and we have to put 1 to 0 index

						r1csVector[counter][0], _ = strconv.Atoi(constraints[i+2])
					}
				}

			} else if j == len(witnesFormal)-1 && r1csVector[counter][0] == 0 { // Check if we compare all variables but can't find any coincidences.
				r1csVector[counter][0], _ = strconv.Atoi(constraintsFormall[i])
				if constraintsFormall[i+1] == "+" { // and we know that in matrix A can't be all values of row equals 0, so, maybe we have 2 + x form of constraint
					for k := 0; k < len(witnes); k++ {
						if witnesFormal[k] == constraintsFormall[i+2] { // Need for 2+x cases, where unknow (didn't store in witness vector) variable at right input
							r1csVector[counter][k] = 1
						}
					}

				}
				break
			}

		}
		i = i + 3             // Next output
		counter = counter + 1 // Next row

	}

	vectorsA = r1csVector

}

// Store indexes numbers wich are not equal 0. This funct needed for next checking of correctness of R1CS matricis
func witnessReprIndexes(vectors [][]int) (indexes [][]int) {

	indexes = make([][]int, len(vectors))
	for i := range indexes {
		indexes[i] = make([]int, 3)
	}

	k := 0
	for i := 0; i < len(vectors); i++ {
		for j := 0; j < len(vectors[i]); j++ { // If value in index equal 1 then save index number
			if vectors[i][j] == 1 {
				indexes[i][k] = j
				k = k + 1
			} else if vectors[i][j] > 1 {
				indexes[i][k] = j
				indexes[i][2] = vectors[i][j]
				k = k + 1
			}
		}
		indexes[i][2] = vectors[i][0]
		k = 0
	}
	return indexes
}

// Print witness in formal form wich builded by R1CS matrcis. needed for next checking of correctness of R1CS matricis
func witnessReprCheckerFormal() {
	a := witnessReprIndexes(vectorsA)
	b := witnessReprIndexes(vectorsB)
	c := witnessReprIndexes(vectorsC)

	fmt.Println(len(a))
	for i := 0; i < len(a); i++ {
		counter := 0
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] >= 1 {
				counter = counter + 1
			}
		}
		if a[i][0] >= 0 && a[i][2] > 0 && counter > 1 {
			witnesFormalChecker = append(witnesFormalChecker, "("+witnesFormal[a[i][1]]+"+"+strconv.Itoa(a[i][2])+")"+"*"+witnesFormal[b[i][0]]+"-"+witnesFormal[c[i][0]])

		} else if a[i][0] > 0 && a[i][1] > 0 {
			witnesFormalChecker = append(witnesFormalChecker, "("+witnesFormal[a[i][0]]+"+"+witnesFormal[a[i][1]]+")"+"*"+witnesFormal[b[i][0]]+"-"+witnesFormal[c[i][0]])

		} else if a[i][2] >= 1 && counter == 1 {
			witnesFormalChecker = append(witnesFormalChecker, strconv.Itoa(a[i][2])+"*"+witnesFormal[b[i][0]]+"-"+witnesFormal[c[i][0]])
		} else {
			witnesFormalChecker = append(witnesFormalChecker, witnesFormal[a[i][0]]+"*"+witnesFormal[b[i][0]]+"-"+witnesFormal[c[i][0]])
		}
	}

}

// Print witness in numbers form wich builded by R1CS matrcis. needed for next checking of correctness of R1CS matricis
func witnessReprChecker() {
	a := witnessReprIndexes(vectorsA)
	b := witnessReprIndexes(vectorsB)
	c := witnessReprIndexes(vectorsC)

	fmt.Println(len(a))
	for i := 0; i < len(a); i++ {
		counter := 0
		for j := 0; j < len(a[i]); j++ {
			if a[i][j] >= 1 {
				counter = counter + 1
			}
		}
		if a[i][0] >= 0 && a[i][2] > 0 && counter > 1 {
			witnesChecker = append(witnesChecker, "("+strconv.Itoa(witnes[a[i][1]])+"+"+strconv.Itoa(a[i][2])+")"+"*"+strconv.Itoa(witnes[b[i][0]])+"-"+strconv.Itoa(witnes[c[i][0]]))

		} else if a[i][0] > 0 && a[i][1] > 0 {
			witnesChecker = append(witnesChecker, "("+strconv.Itoa(witnes[a[i][0]])+"+"+strconv.Itoa(witnes[a[i][1]])+")"+"*"+strconv.Itoa(witnes[b[i][0]])+"-"+strconv.Itoa(witnes[c[i][0]]))

		} else if a[i][2] >= 1 && counter == 1 {
			witnesChecker = append(witnesChecker, strconv.Itoa(a[i][2])+"*"+strconv.Itoa(witnes[b[i][0]])+"-"+strconv.Itoa(witnes[c[i][0]]))
		} else {
			witnesChecker = append(witnesChecker, strconv.Itoa(witnes[a[i][0]])+"*"+strconv.Itoa(witnes[b[i][0]])+"-"+strconv.Itoa(witnes[c[i][0]]))
		}
	}

}

// Return r1cs vectors and matrcis for others packages
func ReturnVectorsA() [][]int {
	return vectorsA
}

func ReturnVectorsB() [][]int {
	return vectorsB
}

func ReturnVectorsC() [][]int {
	return vectorsC
}

func ReturnWitnessFormal() []string {
	return witnesFormal
}

func ReturnWitnessNumbers() []int {
	return witnes
}

func ReturnWitnessFormalChecker() []string {
	return witnesFormalChecker
}

func ReturnWitnessNumberChecker() []string {
	return witnesChecker
}

func ClearAllVar() {
	for key := range mapOfRoots {
		delete(mapOfRoots, key)
	}

	input = ""
	roots_input = ""
	evaluationInput = ""
	constraints = nil
	constraintsFormall = nil
	witnes = nil
	witnesFormal = nil
	witnesChecker = nil
	witnesFormalChecker = nil

	vectorsA = nil
	vectorsB = nil
	vectorsC = nil
}

func Start(function string, roots string) {
	input = function
	roots_input = roots

	//input = "x ^ 3 + x + 5"
	//roots_input = "x = 2 y = 15"

	rootsMap(roots_input)
	evalInput(input, false)

	constraintsFormalForm(ParseInfix(input))

	evalInput(input, false)       // this needed cause here can be another evaluated function in evaluationInput variable
	constraints = constraints[:0] // this needed cause it may be some values in constraints already

	constraintsEval(ParseInfix(evaluationInput))

	witnesInit()

	fmt.Println("Witness in numers representation: ")
	fmt.Println(witnes)
	fmt.Println("Witness in formal representation: ")
	fmt.Println(witnesFormal)

	operatorsPipe("a")
	operatorsPipe("b")
	operatorsPipe("c")
	fmt.Println("R1CS matrix of values a: ")
	fmt.Println(vectorsA)
	fmt.Println("R1CS matrix of values b: ")
	fmt.Println(vectorsB)
	fmt.Println("R1CS matrix of values c: ")
	fmt.Println(vectorsC)

	witnessReprCheckerFormal()
	fmt.Println("Checker of satisfaction of R1CS conditions wich calculate by R1CS matrices (formal representation): ")
	fmt.Println(witnesFormalChecker)

	witnessReprChecker()
	fmt.Println("Checker of satisfaction of R1CS conditions wich calculate by R1CS matrices (numbers representation): ")
	fmt.Println(witnesChecker)

}
