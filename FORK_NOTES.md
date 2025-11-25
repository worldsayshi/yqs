# Fork notes

EDIT: This idea has kinda evolved into:
https://github.com/worldsayshi/monorepo-private-experiments/tree/main/2025/duck-query

Maybe there should be a fork of yqs that allow constructing pipe commands that involve `yq`, `kubectl` and `cat`?

Let's call this project `kqs` for now. As a working name.

The intended UX flow should be something like this:

- Run `kqs` or `<command> | kqs` or press ctrl-g, optionally with a existing command in waiting in the bash prompt to start with.
- If there's no command
- If there is an existing command, `kqs` asks for confirmation to run the command.
- The initial input is parsed as yaml and used as the initial input to use for completions.

## Fork refactorization notes

Some notes on how to achieve good structure:
- I think that `suggestContinuations` can be slightly generalized so that we can have something like this:
```go
func (sg YQSuggester) suggestContinuations(baseExpression, yamlContent string) []string {
    // ...
}
```
(The exact function signature might need to be adjusted somewhat. And initially we might not have any information to put into the sg object.)

- And then we should have a `CombinedSuggester`. Below is a rough sketch of what it should look like:
```go
func (sg CombinedSuggester) suggestContinuations(baseExpression, yamlContent string) []string {
    // If the baseExpression indicate that we have started with a specific command, engage the suggester of that command
    // ...
    if yamlContent == "" {
        // Return suggestion for `cat` or `kubectl`
    } else {
        // Return suggestion `yq`
    }
}
```

# TODO

- [ ] Rename the project to kqs
- [ ] Move `suggestContinuations` implementation and related tests to an internal package
- [ ] Create an initial very simple interface for `Suggester`
- [] Try using an expression like: `compgen -W "$(kubectl --help 2>/dev/null | grep -E '^\s+[a-z]' | awk '{print $1}' | grep -v '^$')" -- ""` (it uses command list from kubectl help)