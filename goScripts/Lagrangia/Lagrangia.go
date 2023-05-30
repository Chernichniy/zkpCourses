// TO DO
// Сделать проверку этих полиномов (в файле R1CS) постфиксный вид и рассчет
package Lagrangia

import (
	"fmt"
	"strconv"
	"strings"
)

var inputLag = "x = 3 x2 = 1 x3 = 2 y = 5 y2 = 7 y3 = 3"
var mapOfRootsLag = make(map[string]int) //map корней и знчений функций

var mapOfPolynomials = make(map[string]string) //полиномы для одной точки
var coordinates []int                          //хранит координаты

var mapOfMultiPolynomials = make(map[string]string) //полиномы, созданные из нескольких точек

var mapOfLagrangiaBasis = make(map[int]string) //хранит базисы

var mapOfBarricNumerator = make(map[int]string)   //хранит числитель конечного баррицентричного вида
var mapOfBarricDenumerator = make(map[int]string) //хранит знаменатель конечного баррицентричного вида

func rootsMapLag(roots string) { // это нужно закинуть в отдеьлный файл
	strRoots := strings.Split(roots, " ")

	for i := 0; i < len(strRoots); i++ {
		if strRoots[i] == "=" {
			var temp, _ = strconv.Atoi(strRoots[i+1])
			mapOfRootsLag[strRoots[i-1]] = temp //помещает корни в key=value форму
		}
	}

}

// вытаскивает одну пару (x,y)
func xy(roots map[string]int, key string) (int, int) {
	temp := strings.Split(key, "")
	temp[0] = "y"
	y := strings.Join(temp, "")
	fmt.Println("Point " + "(" + strconv.Itoa(roots[key]) + "," + strconv.Itoa(roots[y]) + ")")
	return roots[key], roots[y]
}

// сохраняет полиномы в mapOfPolynomials, созданные из одной точки
func savePolynomial(x int, y int, counter int) {
	tempKey := "pol"
	tempKey = tempKey + strconv.Itoa(counter)

	tempPol := strconv.Itoa(y) + " x" + " - " + strconv.Itoa(x+y) + " / " + strconv.Itoa(x) + " - " + strconv.Itoa(x+y)
	mapOfPolynomials[tempKey] = tempPol

}

// сохраняет полиномы в mapOfMultiPolynomials, созданные из нескольких точек
func saveMultPolynomial(xKeys []int, x int, y int, counter int) { //попробовать через поинтеры улучшить (тоесть уменьшить количество переменных)
	tempKey := "pol"
	tempKey = tempKey + strconv.Itoa(counter)

	keysForNumerator := make([]int, len(xKeys))
	keysForDenominator := make([]int, len(xKeys))
	copy(keysForNumerator, xKeys)
	copy(keysForDenominator, xKeys)

	tempVal := "( x - "

	tempValStr := strconv.Itoa(y) + " * "

	for j := 0; j < len(xKeys)-1; j++ {
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

	tempValStr = tempValStr + " / "
	tempVal = "( " + strconv.Itoa(x) + " - "

	for j := 0; j < len(xKeys)-1; j++ {
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

	mapOfMultiPolynomials[tempKey] = tempValStr

}

// создает полиномы лишь для одной точки
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

// создает полиномы из нескольких точек
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
		basisCalc(keys, value, mapOfRootsLag[y], counter2)

		counter2++

	}

}

// Создает следующий по счету key. Генерирует соответсвующую key значение y для key значения x. Получается пара значений (xCOUNTER,yCOUNTER)
// Так же используется в генерации соответсвующего yCOUNTER для key значения map, что бы пропустить значения key=yCOUNTER
func keyGenerate(key string, counter int) (keyNew string) {
	temp := strings.Split(key, "")
	temp[0] = "y"
	keyNew = strings.Join(temp, "")
	return keyNew
}

// считает базис Лагранжа и сохраняет его в
func basisCalc(xKeys []int, x int, y int, counter int) {

	keysForNumerator := make([]int, len(xKeys))
	keysForDenominator := make([]int, len(xKeys))
	copy(keysForNumerator, xKeys)
	copy(keysForDenominator, xKeys)

	tempValStr := "1 / "
	tempVal := "( " + strconv.Itoa(x) + " - "

	for j := 0; j < len(xKeys)-1; j++ {
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

	mapOfLagrangiaBasis[y] = tempValStr

	tempVal = "( x - "

	tempValStr = ""

	for j := 0; j < len(xKeys)-1; j++ {
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

	mapOfBarricNumerator[y] = strconv.Itoa(y) + " * " + mapOfLagrangiaBasis[y] + " / " + tempValStr
	mapOfBarricDenumerator[y] = mapOfLagrangiaBasis[y] + " / " + tempValStr

}

// умножение базиса на значение функции
func numeratorBarricentricCals(map[int]string) map[int]string { //DELETE
	//var str []string
	var basisPlusFunc = make(map[int]string)

	for key, _ := range mapOfLagrangiaBasis {

		basisPlusFunc[key] = strconv.Itoa(key) + " * " + mapOfLagrangiaBasis[key]
	}

	return basisPlusFunc
}

// Возвращает интерполяционный полином Лагранжа в баррицентрическом виде
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

	return "( " + tempNumerator + " ) / ( " + tempDenumerator + " )"
}

// выводит полиномы в виде суммы полиномов(обычный вид), но не считает их
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

// Выводит базисы Лагранжа для каждой точки
func printLagrangiaBasis() {
	var temp []string
	for key, value := range mapOfLagrangiaBasis {
		x := strings.Split(value, "")
		fmt.Println("Point: " + "(" + x[6] + "," + strconv.Itoa(key) + ")")
		fmt.Println(value)
		temp = append(temp, value)
	}
	fmt.Println("\n")
}

func Start() {
	rootsMapLag(inputLag)

	polynomialByOnePoint(mapOfRootsLag)
	polynomialByMultiPoints()
	//fmt.Println(mapOfMultiPolynomials)

	//numeratorBarricentricCals(mapOfLagrangiaBasis)
	resultNormalForm()

	//fmt.Println(mapOfLagrangiaBasis)
	printLagrangiaBasis()
	fmt.Println(resultBarricentricForm())
	//resultNormalForm()
}
