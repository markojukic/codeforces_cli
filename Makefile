all:
	gcc -lm -o executor cmd/executor/main.c
	go build -o codeforces_cli ./cmd/cli/