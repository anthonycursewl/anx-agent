package main

import (
	"fmt"
	"strings"
)

type User struct {
	ID       int
	Username string
	Email    string
	Active   bool
	Roles    []string
	Metadata map[string]interface{}
}

func NewUser(id int, username, email string) *User {
	return &User{
		ID:       id,
		Username: username,
		Email:    email,
		Active:   true,
		Roles:    []string{"user"},
		Metadata: make(map[string]interface{}),
	}
}

func (u *User) AddRole(role string) {
	for _, r := range u.Roles {
		if r == role {
			return
		}
	}
	u.Roles = append(u.Roles, role)
}

func (u *User) Greet() string {
	return fmt.Sprintf("Hello, %s! (ID: %d)", u.Username, u.ID)
}
func ProcessData(input []int) (sum, product int, err error) {
	if len(input) == 0 {
		return 0, 0, fmt.Errorf("input slice is empty")
	}

	sum = 0
	product = 1

	for _, num := range input {
		sum += num
		product *= num
	}

	return sum, product, nil
}

type StringProcessor struct {
	transformer func(string) string
}

func NewStringProcessor(transformer func(string) string) *StringProcessor {
	return &StringProcessor{
		transformer: transformer,
	}
}

func (sp *StringProcessor) Process(input string) string {
	if sp.transformer == nil {
		return input
	}
	return sp.transformer(input)
}

func Example() {
	user := NewUser(1, "testuser", "test@example.com")
	user.AddRole("admin")
	user.AddRole("editor")
	fmt.Println(user.Greet())
	fmt.Println("Roles:", strings.Join(user.Roles, ", "))

	nums := []int{1, 2, 3, 4, 5}
	sum, product, err := ProcessData(nums)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Sum: %d, Product: %d\n", sum, product)
	}

	upper := NewStringProcessor(strings.ToUpper)
	fmt.Println("Uppercase:", upper.Process("hello, world!"))

	reverse := NewStringProcessor(func(s string) string {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	})
	fmt.Println("Reversed:", reverse.Process("hello, world!"))
}

func main() {
	fmt.Println("ANX Agent Test Data")
	fmt.Println("===================")
	Example()
}
