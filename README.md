

Install:
```bash
go install github.com/worldsayshi/yq-suggest@main
source <(yq-completion completion bash)
```


Dev:
```bash
GO_YQ_SUGGEST_DIR=$(pwd)
echo "alias yqc=\"cd "$GO_YQ_SUGGEST_DIR" && go run .\"" >> test.txt
```

## Future work:

Figure out how how to turn it into yq autocompletion. This might not be worth it.