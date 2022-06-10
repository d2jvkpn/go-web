build:
	bash scripts/go_build.sh
	ls -l target/

run:
	# go build -o main main.go
	bash scripts/go_build.sh
	./target/goapp serve
