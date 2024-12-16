package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Использование: go run main.go input.txt output.txt")
		return
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	// Открываем входной файл
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла %s: %v\n", inputFile, err)
		return
	}
	defer file.Close()

	// Очищаем выходной файл или создаем его, если он не существует
	err = ioutil.WriteFile(outputFile, []byte{}, 0644)
	if err != nil {
		fmt.Printf("Ошибка при очистке/создании файла %s: %v\n", outputFile, err)
		return
	}

	output, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Ошибка при открытии файла %s для записи: %v\n", outputFile, err)
		return
	}
	defer output.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`^(\d+)([+\-*/])(\d+)=\?$`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := re.FindStringSubmatch(line)

		if len(matches) == 4 {
			num1, _ := strconv.Atoi(matches[1])
			operator := matches[2]
			num2, _ := strconv.Atoi(matches[3])
			var result int

			switch operator {
			case "+":
				result = num1 + num2
			case "-":
				result = num1 - num2
			case "*":
				result = num1 * num2
			case "/":
				if num2 != 0 {
					result = num1 / num2
				} else {
					fmt.Printf("Ошибка: деление на ноль в выражении %s\n", line)
					continue
				}
			default:
				fmt.Printf("Неизвестный оператор %s в выражении %s\n", operator, line)
				continue
			}

			outputLine := fmt.Sprintf("%s%d\n", line[:len(line)-1], result) // Убираем "=?"
			if _, err := output.WriteString(outputLine); err != nil {
				fmt.Printf("Ошибка при записи в файл %s: %v\n", outputFile, err)
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Ошибка при чтении файла %s: %v\n", inputFile, err)
	}
}
