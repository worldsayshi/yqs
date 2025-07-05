package main

import (
	"flag"
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

// getKeysFromYQOutput extracts keys from the output of a YQ expression using YAML content
func getKeysFromYQOutput(yamlContent, baseExpression string) ([]string, error) {
	// Try to get keys using YQ with piped content
	keysCmd := exec.Command("yq", baseExpression+" | keys")
	keysCmd.Stdin = strings.NewReader(yamlContent)
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
		fieldsCmd := exec.Command("yq", baseExpression+" | to_entries | .[] | .key")
		fieldsCmd.Stdin = strings.NewReader(yamlContent)
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
	// If subPath starts with a pipe, add a space before concatenating
	if strings.HasPrefix(subPath, "|") {
		return base + " " + subPath
	}

	// If subPath starts with array selector, just concatenate
	if strings.HasPrefix(subPath, "[") {
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

// suggestContinuations generates potential continuations for a YQ expression using YAML content
func suggestContinuations(baseExpression, yamlContent string) []string {
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
	keys, err := getKeysFromYQOutput(yamlContent, baseExpression)
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

// testYQExpressionFromContent runs the YQ expression on YAML content and returns its output
func testYQExpressionFromContent(yamlContent, expression string) (string, error) {
	cmd := exec.Command("yq", expression)
	cmd.Stdin = strings.NewReader(yamlContent)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// testYQExpression runs the YQ expression and returns its output (file-based version)
func testYQExpression(yamlPath, expression string) (string, error) {
	cmd := exec.Command("yq", expression, yamlPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// runFzfSelection feeds a list of options to fzf and returns the selected option
func runFzfSelection(options []string) string {
	cmd := exec.Command("fzf", "--header", "Select a YQ expression continuation (ESC to finish)")
	cmd.Stdin = strings.NewReader(strings.Join(options, "\n"))
	output, err := cmd.CombinedOutput()
	// If error status is 130, it means the user cancelled the fzf selection
	if err != nil && err.Error() != "exit status 130" {
		log.Printf("Error running fzf: %v", err)
		return ""
	}
	return strings.TrimSpace(string(output))
}

const bashInstallScript = `#!/bin/bash
suggest_command() {
	if ! command -v yqs &> /dev/null
	then
		echo "yqs command not found. Please install yqs first."
		return
	fi
    all_yaml_files=$(find . -type f -name "*.yaml" -o -name "*.yml")
    if [ -z "$all_yaml_files" ]; then
        echo "No YAML files found in the current directory."
        return
    fi
    yaml_path=$(fzf --header "Select a YAML file" --height 40% --preview 'cat {}' <<< "$all_yaml_files")
    if [ -z "$yaml_path" ]; then
		echo "No YAML file selected."
		return
	fi
	cmd=$(yqs $yaml_path)
    if [ -z "$cmd" ]; then
        echo "No command generated."
        return
    fi
    READLINE_LINE="$cmd"
    READLINE_POINT=${#READLINE_LINE}
}

bind -x '"\C-g": suggest_command'`

func main() {
	// Set up flag usage
	flag.Usage = func() {
		fmt.Println("Usage: " + os.Args[0] + " <yaml_file_path> [base_yq_expression]")
		fmt.Println("Or use `source <(yqs --command-installation bash)` to install the command in your shell.")
		fmt.Println()
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	// Add command installation flag
	installCommand := flag.String("command-installation", "", "Output shell configuration (supported: bash)")

	// Parse flags
	flag.Parse()

	if *installCommand == "bash" {
		fmt.Println(bashInstallScript)
		os.Exit(0)
	}

	// Check for correct number of arguments
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}
	yamlPath := flag.Arg(0)

	// Read the YAML file content once at the beginning
	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		fmt.Printf("Error reading YAML file '%s': %v\n", yamlPath, err)
		os.Exit(1)
	}

	var baseExpression string
	if flag.NArg() < 2 {
		baseExpression = "."
	} else {
		baseExpression = flag.Arg(1)
	}

	// Loop until user selects empty or cancels
	currentExpression := baseExpression
	for {
		// Test the current expression to make sure it's valid
		_, err := testYQExpression(yamlPath, currentExpression)
		if err != nil {
			if currentExpression == baseExpression {
				fmt.Println("Base expression is invalid")
				os.Exit(1)
			} else {
				fmt.Printf("Expression '%s' is invalid, stopping\n", currentExpression)
				break
			}
		}

		// Generate potential continuations based on the current expression
		continuations := suggestContinuations(currentExpression, string(yamlContent))

		// Collect valid continuations
		var validContinuations []string
		for _, cont := range continuations {
			output, err := testYQExpression(yamlPath, cont)
			if err == nil && output != "" && output != "null" {
				validContinuations = append(validContinuations, cont)
			}
		}

		// If we have valid continuations, feed them to fzf for selection
		if len(validContinuations) > 0 {
			selected := runFzfSelection(validContinuations)
			if selected == "" {
				// User cancelled or selected empty, break and print final command
				break
			}

			// Update current expression and continue
			currentExpression = selected
		} else {
			// No valid continuations available
			break
		}
	}

	// Print the final command
	fmt.Printf("yq '%s' %s\n", currentExpression, yamlPath)
}
