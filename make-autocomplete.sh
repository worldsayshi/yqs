#!/bin/bash

# Custom completion function for 'echo'
_echo_completion() {
    local cur prev
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    echo "Current word: $cur"

    # If the previous word is 'echo', suggest 'hello'
    if [[ "$prev" == "echo" ]]; then
        COMPREPLY=( $(compgen -W "hello" -- "$cur") )
        return 0
    fi
}

# Set up the completion
setup_echo_completion() {
    if [[ -n "$BASH_VERSION" ]]; then
        # Register the completion function for the 'echo' command
        complete -F _echo_completion echo
        echo "Custom echo completion enabled. Type 'echo ' and press Tab to get 'hello' suggestion."
    else
        echo "This script requires Bash to work properly."
    fi
}

# If this script is being sourced (not executed), set up the completion
if [[ "${BASH_SOURCE[0]}" != "${0}" ]]; then
    echo "Sourcing..."
    setup_echo_completion
else
    # If executed directly, provide instructions
    echo "To enable custom echo completion, add the following to your ~/.bashrc:"
    echo "source $(realpath "${BASH_SOURCE[0]}")"
    echo ""
    echo "After sourcing, when you type 'echo ' and press Tab, 'hello' will be suggested."
fi