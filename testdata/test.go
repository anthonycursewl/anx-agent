package main

import "fmt"

func FakeFunctionOne() {
}

func FakeFunctionTwo(inputValue int) int {
	return inputValue * 2
}

func FakeFunctionThree(arg1, arg2 int) string {
	result := arg1 + arg2
	return fmt.Sprintf("Processed: %d", result)
}

func AnotherFakeFunction() string {
	return "Fake function executed."
}

func GetFullName(firstName, lastName string) string {
	return firstName + " " + lastName
}

func CalculateArea(length, width float64) float64 {
	return length * width
}

func CompareNumbers(a, b int) string {
	if a > b {
		return fmt.Sprintf("%d is greater than %d", a, b)
	} else if a < b {
		return fmt.Sprintf("%d is less than %d", a, b)
	} else {
		return fmt.Sprintf("%d is equal to %d", a, b)
	}
}