`yqs` is a complement to [`yq`](https://github.com/mikefarah/yq) that gives you clever yq query completion in the terminal when you press ctrl-g.

Install:
```bash
go install github.com/worldsayshi/yqs@v0.0.3
# And put the following in your ~/.bashrc or equivalent:
source <(yqs --command-installation bash)
```


Dev:
```bash
GO_YQ_SUGGEST_DIR=$(pwd)
alias yqs="cd $GO_YQ_SUGGEST_DIR && go run ."
source <(yqs --command-installation bash)
```
Then press ctrl-g.

## Demo

[![asciicast](https://asciinema.org/a/K9IyNRSgoBm7FsurAw35V7utu.svg)](https://asciinema.org/a/K9IyNRSgoBm7FsurAw35V7utu)

## TODO's

- [ ] Make key mapping customizable

## Future work:

Figure out how how to turn it into yq tab autocompletion. This is probably not worth it.
