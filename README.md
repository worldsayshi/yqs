

Install:
```bash
go install github.com/worldsayshi/yqs@v0.0.1
source <(yqs --command-installation bash)
```


Dev:
```bash
GO_YQ_SUGGEST_DIR=$(pwd)
alias yqs="cd $GO_YQ_SUGGEST_DIR && go run ."
source <(yqs --command-installation bash)
```
Then press ctrl-g.

## TODO's

- [ ] Make key mapping customizable

## Future work:

Figure out how how to turn it into yq tab autocompletion. This is probably not worth it.
