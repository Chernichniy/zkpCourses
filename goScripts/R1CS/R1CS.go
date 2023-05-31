package r1cs

import (
	"fmt"
	"strconv"
	"strings"
)

var input string

// var input = "x ^ 3 + x + 5"
// var roots_input = "x = 3 y = 35" //эти две переменные мы получаем с сайта
var roots_input = "x = 1 z = 2" //эти две переменные мы получаем с сайта

var mapOfRoots = make(map[string]int) //обработанные значения корней в виде key = value
var evaluationInput string            // строка функции с подставленными private and public inputs

var constraints []string        //constraints в численной форме
var constraintsFormall []string //constraints в буквенной форме

var witnes []int          //witness в численной форме
var witnesFormal []string //witness в буквенной форме

var vectorsA [][]int
var vectorsB [][]int
var vectorsC [][]int

var witnesFromalChecker []string
var witnesChecker []string

//var zeroOneVector [][]int // хранит R1CS веткора (0,0,0,1,0,0)...

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

// переводит из инфиксного вида записи функции в постфиксный
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
	a := oper   //нужно для проверки операторов
	numCon := 1 //счетчик номера constraint
	var constrNum = []string{"Con", ""}

	for i := 0; i < len(str); i++ {

		for key := range a {
			constrNum[1] = strconv.Itoa(numCon)

			if str[i] == key && key == "^" {
				first := str[i-2]
				second := first
				num, _ := strconv.Atoi(str[i-1])

				for j := num; j > 1; j-- {
					//num = numTemp * num1
					constraintsSaver(first, "*", second, strings.Join(constrNum, ""), false)
					//numTemp = num
					numCon = numCon + 1
					first = strings.Join(constrNum, "")
					constrNum[1] = strconv.Itoa(numCon)
				}
				//constraintsSaver(str[i], str[i+2], str[i+1], strconv.Itoa(num))
				//str[i+2] = strconv.Itoa(num)
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
	witnesFormal = append(witnesFormal, strconv.Itoa(1))
	for key, val := range mapOfRoots {

		witnesFormal = append(witnesFormal, key)
		witnes = append(witnes, val)
	}

	witnessAdd()
	formalWitness()
}

// создает формальный вид witness
func formalWitness() {
	for i := 0; i < len(constraintsFormall); i++ {
		if (i+1)%4 == 0 {
			//temp, _ := strconv.Atoi(constraints[i])
			witnesFormal = append(witnesFormal, constraintsFormall[i])
		}
	}
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

// создает нулевой двумерный срез для хранения векторов (1,0,0,0,1)...
func zeroOneVectorFulling() (zeroOneVector [][]int) {
	zeroOneVector = make([][]int, len(constraintsFormall)/4)
	for i := range zeroOneVector {
		zeroOneVector[i] = make([]int, len(witnes))
	}
	return zeroOneVector
}

// распределяет какую функцию вызывать в зависимости от оператора  a*b-c=0 для заполнения двумерного среза
func operatorsPipe(operator string) {
	switch operator {
	case "a":
		R1CSCompilerOperatorA()
	case "b":
		R1CSCompilerOperatorB()
	case "c":
		R1CSCompilerOperatorC()
	}
}

// заполняет дмуерный срез для оператора c
func R1CSCompilerOperatorC() { //тут конкретно много if

	r1csVector := zeroOneVectorFulling()
	//fmt.Println(r1csVector)
	counter := 0

	for i := 3; i < len(constraintsFormall); i++ {
		for j := 0; j < len(witnesFormal); j++ {
			if constraintsFormall[i] == witnesFormal[j] { //проблема: если констраинт 2*x, а двойки в витнесе нет
				r1csVector[counter][j] = 1
				break
			}

		}
		i = i + 3
		counter = counter + 1

	}
	vectorsC = r1csVector

}

// заполняет дмуерный срез для оператора b
func R1CSCompilerOperatorB() { //тут конкретно много if

	r1csVector := zeroOneVectorFulling()
	//fmt.Println(r1csVector)
	counter := 0

	for i := 2; i < len(constraintsFormall); i++ {
		for j := 0; j < len(witnesFormal); j++ {
			if constraintsFormall[i] == witnesFormal[j] && constraintsFormall[i-1] != "+" { //проблема: если констраинт 2*x, а двойки в витнесе нет
				r1csVector[counter][j] = 1
				break

			} else if /*constraintsFormall[i] == witnesFormal[j] &&*/ constraintsFormall[i-1] == "+" {
				r1csVector[counter][0] = 1
				break
			} else if j == len(witnesFormal)-1 {
				r1csVector[counter][0], _ = strconv.Atoi(constraintsFormall[i])
				break
			}

		}
		i = i + 3
		counter = counter + 1

	}
	vectorsB = r1csVector

}

// заполняет дмуерный срез для оператора a
func R1CSCompilerOperatorA() { //тут конкретно много if

	r1csVector := zeroOneVectorFulling()
	//fmt.Println(r1csVector)
	counter := 0

	for i := 0; i < len(constraintsFormall); i++ {
		for j := 0; j < len(witnesFormal); j++ {
			if constraintsFormall[i] == witnesFormal[j] && constraintsFormall[i+1] != "+" {
				r1csVector[counter][j] = 1
				break

			} else if constraintsFormall[i] == witnesFormal[j] && constraintsFormall[i+1] == "+" {

				//num, _ := strconv.Atoi(constraints[i+2])
				r1csVector[counter][j] = 1
				for k := 0; k < len(witnesFormal); k++ {
					if witnesFormal[k] == constraintsFormall[i+2] {
						r1csVector[counter][k] = 1
						break
					} else if k == len(witnesFormal)-1 {

						r1csVector[counter][0], _ = strconv.Atoi(constraints[i+2])
					}
				}
				//r1csVector[counter][0] = num

			} else if j == len(witnesFormal)-1 && r1csVector[counter][0] == 0 {
				r1csVector[counter][0], _ = strconv.Atoi(constraintsFormall[i])
				break
			}

		}
		i = i + 3
		counter = counter + 1

	}

	vectorsA = r1csVector

}

// сохраняет значения индексов, в которых находятся ненулевые значения двумерного среза
func witnessReprIndexes(vectors [][]int) (indexes [][]int) {

	indexes = make([][]int, len(vectors))
	for i := range indexes {
		indexes[i] = make([]int, 3)
	}

	k := 0
	for i := 0; i < len(vectors); i++ {
		for j := 0; j < len(vectors[i]); j++ {
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

// выводит witness построенный на основе двумерного среза в формальном виде. Служит для проверки правильности вычеслений
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
			witnesFromalChecker = append(witnesFromalChecker, "("+witnesFormal[a[i][1]]+"+"+strconv.Itoa(a[i][2])+")"+"*"+witnesFormal[b[i][0]]+"-"+witnesFormal[c[i][0]])

		} else if a[i][0] > 0 && a[i][1] > 0 {
			witnesFromalChecker = append(witnesFromalChecker, "("+witnesFormal[a[i][0]]+"+"+witnesFormal[a[i][1]]+")"+"*"+witnesFormal[b[i][0]]+"-"+witnesFormal[c[i][0]])

		} else if a[i][2] >= 1 && counter == 1 {
			witnesFromalChecker = append(witnesFromalChecker, strconv.Itoa(a[i][2])+"*"+witnesFormal[b[i][0]]+"-"+witnesFormal[c[i][0]])
		} else {
			witnesFromalChecker = append(witnesFromalChecker, witnesFormal[a[i][0]]+"*"+witnesFormal[b[i][0]]+"-"+witnesFormal[c[i][0]])
		}
	}

}

// выводит witness построенный на основе двумерного среза в числовом виде. Служит для проверки правильности вычеслений
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

func Start(input string) {

	rootsMap(roots_input)
	evalInput(input)
	//fmt.Println("infix:  ", evaluationInput)
	//parseInfix(evaluationInput)
	//fmt.Println("postfix:", parseInfix(evaluationInput))
	//fmt.Println("infix:  ", input)
	//fmt.Println("postfix:", parseInfix(input))

	//fmt.Println("res:", circuitEval(parseInfix(input)))
	//fmt.Println("res:", circuitEval(parseInfix(input)))
	constraintsFormalForm(parseInfix(input))
	//fmt.Println(constraintsFormall)

	//fmt.Println("Evaluation input\n " + evaluationInput)
	constraintsEval(parseInfix(evaluationInput))
	//fmt.Println(constraints)

	witnesInit()
	/*
		fmt.Print("a")
		fmt.Println("infix:  ", input)
		fmt.Println("postfix:", parseInfix(input))
		fmt.Println("res:", binaryTree(parseInfix(input)))
	*/

	//witnessAdd()
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

	//fmt.Print(len(vectorsA))
	witnessReprCheckerFormal()
	fmt.Println("Checker of satisfaction of R1CS conditions wich calculate by R1CS matrices (formal representation): ")
	fmt.Println(witnesFromalChecker)

	witnessReprChecker()
	fmt.Println("Checker of satisfaction of R1CS conditions wich calculate by R1CS matrices (numbers representation): ")
	fmt.Println(witnesChecker)

}
