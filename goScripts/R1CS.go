package main

import (
	"fmt"
	"strconv"
	"strings"
)

// var input = "x ^ 2 + 2 * x * z + z ^ 2 + z + 1"
var input = "x ^ 3 + x + 5"
var roots_input = "x = 3 y = 35" //эти две переменные мы получаем с сайта
//var roots_input = "x = 1 z = 2" //эти две переменные мы получаем с сайта

var mapOfRoots = make(map[string]int) //обработанные значения корней в виде key = value
var evaluationInput string            // строка функции с подставленными private and public inputs

var constraints []string
var constraintsFormall []string

var witnes []int

/*
	type Circuit struct {
		m int //number of inputs
		n int //number of gates

}
*/
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

func parseInfix(e string) (rpn string) { //попробовать свою реализацию сддеалть https://github.com/codefreezr/rosettacode-to-go/blob/df006db732e5/tasks/Parsing-Shunting-yard-algorithm/parsing-shunting-yard-algorithm.go
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

// constraints evaluation way: true for evaluation form(with roots), false for formal view
func constraintsSaver(lftInput string, operation string, rghtInput string, output string, types bool) {
	switch types {
	case true:
		if operation == "^" { // приводим к виду только сумм и произведения
			rghtInput = lftInput
			operation = "*"
		}
		constraints = append(constraints, lftInput) // помещаем поочередно
		constraints = append(constraints, operation)

		constraints = append(constraints, rghtInput)
		constraints = append(constraints, output)

	case false:
		if operation == "^" { // приводим к виду только сумм и произведения
			rghtInput = lftInput
			operation = "*"
		}
		constraintsFormall = append(constraintsFormall, lftInput) // помещаем поочередно
		constraintsFormall = append(constraintsFormall, operation)

		constraintsFormall = append(constraintsFormall, rghtInput)
		constraintsFormall = append(constraintsFormall, output)
	}

}

// Create a constraints with their "way" in formal form
func constraintsFormalForm(rpn string) {
	str := strings.Split(rpn, " ")
	a := oper
	numCon := 1
	var constrNum = []string{"Con", ""}

	for i := 0; i < len(str); i++ {

		for key := range a {
			if str[i] == key {
				constrNum[1] = strconv.Itoa(numCon)
				constraintsSaver(str[i-2], str[i], str[i-1], strings.Join(constrNum, ""), false)
				numCon = numCon + 1
				str[i] = strings.Join(constrNum, "")
				str = append(str[:i-2], str[i:]...)
				i = -1
				break
			}
		}

	}

}

// Create a constraint with their "way" in evaluate form
func constraintsEval(rpn string) (res int) {
	str := strings.Split(rpn, " ")
	//str = append(str, "s")
	//	var flag bool //0 означает уровень остается, 1 - перехд на уровень выше
	//for j:= range strings.Fields(rpn){
	for i := 0; i < len(str); i++ /*range str*/ {
		if i+2 <= len(str)-1 {
			num1, _ := strconv.Atoi(str[i])
			num2, _ := strconv.Atoi(str[i+1])
			num3, _ := strconv.Atoi(str[i+2])
			if num3 == 0 { //нужно перегрузка операторов или что-то типо лямба функции

				switch str[i+2] {
				case "+":
					num := num1 + num2
					constraintsSaver(str[i], str[i+2], str[i+1], strconv.Itoa(num), true)
					str[i+2] = strconv.Itoa(num)
					str = append(str[:i], str[i+2:]...)
					i = -1
				case "*":
					num := num1 * num2
					constraintsSaver(str[i], str[i+2], str[i+1], strconv.Itoa(num), true)
					str[i+2] = strconv.Itoa(num)
					str = append(str[:i], str[i+2:]...)
					i = -1
				case "^":
					numTemp := num1
					num := num1

					for j := num2; j > 1; j-- {
						num = numTemp * num1
						constraintsSaver(strconv.Itoa(numTemp), "*", str[i], strconv.Itoa(num), true)
						numTemp = num
					}
					//constraintsSaver(str[i], str[i+2], str[i+1], strconv.Itoa(num))
					str[i+2] = strconv.Itoa(num)
					str = append(str[:i], str[i+2:]...)
					i = -1
				}
			}
		}
	}
	res, _ = strconv.Atoi(str[0])
	return res
	// }
}

// вытаскивает private и public inputs из programm input
func rootsMap(roots string) {
	strRoots := strings.Split(roots, " ")

	for i := 0; i < len(strRoots); i++ {
		if strRoots[i] == "=" {
			var temp, _ = strconv.Atoi(strRoots[i+1])
			mapOfRoots[strRoots[i-1]] = temp //помещает корни в key=value форму
		}
	}
	fmt.Println(mapOfRoots)

}

// записывает функцию уже с подставленными корнями
func evalInput(function string) {
	strFunc := strings.Split(function, " ")
	for i := 0; i < len(strFunc); i++ {
		_, ok := mapOfRoots[strFunc[i]]
		if ok == true {
			strFunc[i] = strconv.Itoa(mapOfRoots[strFunc[i]])
		}
	}
	fmt.Println(strFunc)
	evaluationInput = strings.Join(strFunc, " ")
}

// генерирует witness
func witnesInit() {
	witnes = append(witnes, 1)
	for _, val := range mapOfRoots {

		witnes = append(witnes, val)
	}
	witnessAdd()
	formalWitness()
}

// создает формальный вид witness
func formalWitness() {

}

// записывает constraint`s output witnes
func witnessAdd() {
	for i := 0; i < len(constraints); i++ {
		if (i+1)%4 == 0 {
			temp, _ := strconv.Atoi(constraints[i])
			witnes = append(witnes, temp)
		}
	}
}

func main() {
	rootsMap(roots_input)
	evalInput(input)
	//fmt.Println("infix:  ", evaluationInput)
	//parseInfix(evaluationInput)
	fmt.Println("postfix:", parseInfix(evaluationInput))
	//fmt.Println("infix:  ", input)
	fmt.Println("postfix:", parseInfix(input))

	//fmt.Println("res:", circuitEval(parseInfix(input)))
	//fmt.Println("res:", circuitEval(parseInfix(input)))
	constraintsFormalForm(parseInfix(input))
	fmt.Println(constraintsFormall)

	fmt.Println("Evaluation input\n " + evaluationInput)
	constraintsEval(parseInfix(evaluationInput))
	fmt.Println(constraints)

	witnesInit()
	/*
		fmt.Print("a")
		fmt.Println("infix:  ", input)
		fmt.Println("postfix:", parseInfix(input))
		fmt.Println("res:", binaryTree(parseInfix(input)))
	*/

	//witnessAdd()
	fmt.Println(witnes)

}
