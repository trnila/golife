# golife
Simple parallel approach for Game Of Life in Golang.

Each generation splits world to parts that are handled by lightweight goroutines concurrently.
Current state of world is read from read-only matrix and updated state is written to write-only matrix.
Matrices are swapped before each generation.

[![asciicast](https://asciinema.org/a/iPwGXcExEz9YnFBEtZ7FFpabb.png)](https://asciinema.org/a/iPwGXcExEz9YnFBEtZ7FFpabb)

## Install
Assuming you have correctly set your **$GOPATH** and **$PATH**.

```sh
go get github.com/trnila/golife
curl https://raw.githubusercontent.com/trnila/golife/master/boards/gosper_gun | golife
```

## Local build
```
git clone https://github.com/trnila/golife
cd golife
go build .
./golife boards/gosper_gun
```

## Limit of cores
Number of used cores for parallel computation can be changed with environment variable **GOMAXPROCS**, eg:
```
GOMAXPROCS=2 golife
```

## TUI Controls
You can move in world with arrows keys and zoom in (**+**) or zoom out (**-**).
