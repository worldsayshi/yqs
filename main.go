package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

/*
This is a go program that does the following:
* expects a path to a yaml file as input
* expects an yq expression as input
* generate potential continuations of the yq expression by looking at the output and potentially the grammar
* test these continuations against the input and print those that give output
*/

// getKeysFromYQOutput extracts keys from the output of a YQ expression
func getKeysFromYQOutput(yamlPath, baseExpression string) ([]string, error) {
	// Try to get keys using YQ
	keysCmd := exec.Command("yq", baseExpression+" | keys", yamlPath)
	keysOutput, err := keysCmd.CombinedOutput()

	var keys []string

	if err == nil && len(keysOutput) > 0 && !strings.Contains(string(keysOutput), "Error") {
		// Parse the keys from the output
		lines := strings.Split(strings.TrimSpace(string(keysOutput)), "\n")
		for _, line := range lines {
			trimmedLine := strings.TrimSpace(line)
			// Skip empty lines and array indicators like "-"
			if trimmedLine != "" && trimmedLine != "-" {
				// Remove any leading "- " that yq might output for arrays
				key := strings.TrimPrefix(trimmedLine, "- ")
				keys = append(keys, key)
			}
		}
	}

	// Try another approach if the first didn't work - check if it's an object directly
	if len(keys) == 0 {
		// Sometimes direct properties are better detected with this approach
		fieldsCmd := exec.Command("yq", baseExpression+" | to_entries | .[] | .key", yamlPath)
		fieldsOutput, err := fieldsCmd.CombinedOutput()

		if err == nil && len(fieldsOutput) > 0 && !strings.Contains(string(fieldsOutput), "Error") {
			lines := strings.Split(strings.TrimSpace(string(fieldsOutput)), "\n")
			for _, line := range lines {
				trimmedLine := strings.TrimSpace(line)
				if trimmedLine != "" && trimmedLine != "-" {
					key := strings.TrimPrefix(trimmedLine, "- ")
					keys = append(keys, key)
				}
			}
		}
	}

	return keys, nil
}

// joinYQExpressions properly joins a base YQ expression with a sub-path or operation
func joinYQExpressions(base, subPath string) string {
	// If subPath starts with a pipe or other operator, just concatenate
	if strings.HasPrefix(subPath, "|") || strings.HasPrefix(subPath, "[") {
		return base + subPath
	}

	// Handle the root expression case
	if base == "." && strings.HasPrefix(subPath, ".") {
		// Avoid having ".." in the result
		return subPath
	} else if base == "." && !strings.HasPrefix(subPath, ".") {
		// Add dot before a key name if needed
		return base + subPath
	} else if !strings.HasPrefix(subPath, ".") && !strings.HasPrefix(subPath, "[") {
		// Add separator dot between base and subPath
		return base + "." + subPath
	}

	// Default concatenation for other cases
	return base + subPath
}

// suggestContinuations generates potential continuations for a YQ expression
func suggestContinuations(baseExpression, yamlPath string) []string {
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
		fullContinuations = append(fullContinuations, joinYQExpressions(baseExpression, cont))
	}

	// Add continuations based on the keys in the output
	keys, err := getKeysFromYQOutput(yamlPath, baseExpression)
	if err == nil {
		for _, key := range keys {
			// For numeric or simple keys
			fullContinuations = append(fullContinuations, joinYQExpressions(baseExpression, "."+key))

			// For keys that might contain special characters
			fullContinuations = append(fullContinuations, joinYQExpressions(baseExpression, fmt.Sprintf(".[%q]", key)))

			// Common operations on specific keys
			fullContinuations = append(fullContinuations, fmt.Sprintf("%s | has(%q)", baseExpression, key))
		}
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
	if len(os.Args) < 2 {
		fmt.Println("Usage: yq-continuation-generator <yaml_file_path> <base_yq_expression>")
		os.Exit(1)
	}
	yamlPath := os.Args[1]
	log.Println("length of args", len(os.Args))
	var baseExpression string
	if len(os.Args) < 3 {
		baseExpression = "."
	} else {
		baseExpression = os.Args[2]
	}
	// baseExpression := os.Args[2]
	// if the base expression is empty, set it to "."
	// if baseExpression == "" {
	// 	baseExpression = "."
	// }

	// Generate potential continuations
	continuations := suggestContinuations(baseExpression, yamlPath)

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
