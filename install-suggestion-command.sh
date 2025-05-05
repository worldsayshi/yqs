#!/bin/bash

# If yqs is not availible as a command, assume that we're
# in the dev folder so yqs should refer to `go run .`.
if ! command -v yqs &> /dev/null
then
    # Check if we're in the dev folder
    if [ ! -f "go.mod" ]; then
        echo "yqs not found and not in dev folder."
        exit 1
    fi
    yqs="go run ."
else
    yqs="yqs"
fi

suggest_command() {

    all_yaml_files=$(find . -type f -name "*.yaml" -o -name "*.yml")
    if [ -z "$all_yaml_files" ]; then
        echo "No YAML files found in the current directory."
        return
    fi
    yaml_path=$(fzf --header "Select a YAML file" --height 40% --preview 'cat {}' <<< "$all_yaml_files")
    cmd=$($yqs $yaml_path)
    if [ -z "$cmd" ]; then
        echo "No command generated."
        return
    fi
    READLINE_LINE="$cmd"
    READLINE_POINT=${#READLINE_LINE}
}

bind -x '"\C-g": suggest_command'