package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// suggestContinuations generates potential continuations for a YQ expression
func suggestContinuations(baseExpression string) []string {
	// List of common YQ operators and selectors
	continuations := []string{
		// Filtering and selection
		".[]",
		".[0]",
		".[*]",
		".select()",

		// Transformation
		".map()",
		".flatten()",
		".sort()",
		".reverse()",

		// String operations
		".to_string()",
		".to_json()",
		".to_yaml()",

		// Arithmetic and comparison
		"| length",
		"| keys",
		"| has()",

		// Logical operators
		"| contains()",
		"| any()",
		"| all()",
	}

	var fullContinuations []string
	for _, cont := range continuations {
		fullContinuations = append(fullContinuations, baseExpression+cont)
	}

	return fullContinuations
}

// testYQExpression runs the YQ expression and returns its output
func testYQExpression(yamlPath, expression string) (string, error) {
	cmd := exec.Command("yq", expression, yamlPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func main() {
	// Check for correct number of arguments
	if len(os.Args) != 3 {
		fmt.Println("Usage: yq-continuation-generator <yaml_file_path> <base_yq_expression>")
		os.Exit(1)
	}

	yamlPath := os.Args[1]
	baseExpression := os.Args[2]

	// Generate potential continuations
	continuations := suggestContinuations(baseExpression)

	fmt.Println("Testing potential YQ expression continuations:")
	fmt.Println("----------------------------------------------")

	output, err := testYQExpression(yamlPath, baseExpression)
	if err != nil {
		fmt.Println("Base expression is invalid")
		os.Exit(1)
	}
	fmt.Println("Output of base expression:")
	fmt.Println(output)
	// Test each continuation
	for _, cont := range continuations {
		output, err := testYQExpression(yamlPath, cont)

		if err == nil && output != "" && output != "null" {
			fmt.Printf("Continuation: %s\n", cont)
		}
	}
}
