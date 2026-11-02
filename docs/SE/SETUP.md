create a new golang project with two modules, one named loomclient, the other loomserver. Both have a main function.
+Thought · 1.3s
$ pwd && ls -la && go version
/home/rotscher/Documents/git/github.com/rotscher/httploom
total 8
drwxrwxr-x 2 rotscher rotscher 4096 Nov  2 22:28 .
drwxrwxr-x 8 rotscher rotscher 4096 Nov  2 22:28 ..
go version go1.27.3 linux/amd64
+Thought · 1.3s
$ mkdir -p loomclient loomserver && go mod init github.com/rotscher/httploom/loomclient && mv go.mod loomclient/go.mod && cd loomserver && go mod init github.com/rotscher/httploom/loomserver
go: creating new go.mod: module github.com/rotscher/httploom/loomclient
go: to add module requirements and sums:
	go mod tidy
go: creating new go.mod: module github.com/rotscher/httploom/loomserver
# Wroteloomclient/main.go
package main

import "fmt"

func main() {
	fmt.Println("loomclient")
}

# Wroteloomserver/main.go
package main

import "fmt"

func main() {
	fmt.Println("loomserver")
}

$ go work init ./loomclient ./loomserver && gofmt -l . && go build ./... && cat go.work && find . -maxdepth 2 -type f | sort
pattern ./...: directory prefix . does not contain modules listed in go.work or their selected dependencies

Exited with code 1
+Thought · 1.1s
$ gofmt -l . ; (cd loomclient && go build ./...) && (cd loomserver && go build ./...) && cat go.work && find . -maxdepth 2 -type f | sort && ./loomclient/loomclient 2>/dev/null || true
(8 earlier lines)
./loomclient/loomclient
./loomclient/main.go
./loomserver/go.mod
./loomserver/loomserver
./loomserver/main.go
loomclient
+Thought · 1.3s
$ rm -f loomclient/loomclient loomserver/loomserver && printf '/loomclient/loomclient\n/loomserver/loomserver\n' > .gitignore && cat > README.md <<'EOF'
# httploom…
./.gitignore
./go.work
./loomclient/go.mod
./loomclient/main.go
./loomserver/go.mod
./loomserver/main.go
./README.md
Done. Created a Go workspace with two independent modules:
httploom/
├── go.work
├── .gitignore
├── README.md
├── loomclient/
│   ├── go.mod   (module github.com/rotscher/httploom/loomclient)
│   └── main.go  (func main)
└── loomserver/
    ├── go.mod   (module github.com/rotscher/httploom/loomserver)
    └── main.go  (func main)
Both modules are joined by go.work, each has its own go.mod, and each main.go has a main function printing its name. Verified both compile.
Run them with:
go run ./loomclient
go run ./loomserver
Note: go build ./... from the repo root fails because the root isn't a module (standard for multi-module workspaces); build per module instead, e.g. cd loomclient && go build ./....
Build · Big Pickle · 20.7s · 64.4 tok/s