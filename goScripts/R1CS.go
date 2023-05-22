package main

import (
	"fmt"
	"strconv"
	"strings"
)

var input = "x ^ 2 + 2 * x * z + z ^ 2 + z + 1"
var roots_input = "x = 1 z = 2" //эти две переменные мы получаем с сайта

var mapOfRoots = make(map[string]int) //обработанные значения корней в виде key = value
var evaluationInput string            // строка функции с подставленными private and public inputs

var constraints []string

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

func constraintsSaver(lftInput string, operation string, rghtInput string, output string) { //constraints evaluation way
	if operation == "^" { // приводим к виду только сумм и произведения
		rghtInput = lftInput
		operation = "*"
	}
	constraints = append(constraints, lftInput) // помещаем поочередно
	constraints = append(constraints, operation)

	constraints = append(constraints, rghtInput)
	constraints = append(constraints, output)
}

func binaryTree(rpn string) (res int) { //Здесь нужно изменить под contsraints (функция просто считает результат)
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
					constraintsSaver(str[i], str[i+2], str[i+1], strconv.Itoa(num))
					str[i+2] = strconv.Itoa(num)
					str = append(str[:i], str[i+2:]...)
					i = -1
				case "*":
					num := num1 * num2
					constraintsSaver(str[i], str[i+2], str[i+1], strconv.Itoa(num))
					str[i+2] = strconv.Itoa(num)
					str = append(str[:i], str[i+2:]...)
					i = -1
				case "^":
					num := num1 * num2
					constraintsSaver(str[i], str[i+2], str[i+1], strconv.Itoa(num))
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

func rootsMap(roots string) { // обрабатывает private и public inputs
	strRoots := strings.Split(roots, " ")

	for i := 0; i < len(strRoots); i++ {
		if strRoots[i] == "=" {
			var temp, _ = strconv.Atoi(strRoots[i+1])
			mapOfRoots[strRoots[i-1]] = temp //помещает корни в key=value форму
		}
	}
	fmt.Println(mapOfRoots)

}

func evalInput(function string) { // записывает функцию уже с корнями
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

func witnessGen() {
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
	fmt.Println("infix:  ", evaluationInput)
	parseInfix(evaluationInput)
	fmt.Println("postfix:", parseInfix(evaluationInput))
	fmt.Println("res:", binaryTree(parseInfix(evaluationInput)))
	/*
		fmt.Print("a")
		fmt.Println("infix:  ", input)
		fmt.Println("postfix:", parseInfix(input))
		fmt.Println("res:", binaryTree(parseInfix(input)))
	*/

	fmt.Println(constraints)
	witnessGen()
	fmt.Println(witnes)

}
