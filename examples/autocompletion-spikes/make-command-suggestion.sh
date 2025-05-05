#!/bin/bash

# Function to suggest a command
suggest_command() {
    local cmd="echo hello"
    READLINE_LINE="$cmd"
    READLINE_POINT=${#READLINE_LINE}
}

# Add this to your ~/.bashrc
setup_command_suggestion() {
    # Bind Ctrl+M to the suggest_command function
    # Note: Ctrl+M is the same as Enter in many terminals, so we'll use Ctrl+G instead
    if [[ -n "$BASH_VERSION" ]]; then
        bind -x '"\C-g": suggest_command'
        echo "Command suggestion enabled. Press Ctrl+G to suggest 'echo hello'."
    else
        echo "This script requires Bash to work properly."
    fi
}

# If this script is being sourced (not executed), set up the binding
if [[ "${BASH_SOURCE[0]}" != "${0}" ]]; then
    setup_command_suggestion
else
    # If executed directly, provide instructions
    echo "To enable command suggestions, add the following to your ~/.bashrc:"
    echo "source $(realpath "${BASH_SOURCE[0]}")"
    echo ""
    echo "This will bind Ctrl+G to suggest 'echo hello'. You can edit the suggestion before pressing Enter."
fi