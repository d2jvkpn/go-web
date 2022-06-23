build:
	bash scripts/go_build.sh

run:
	bash scripts/go_build.sh
	./target/go-web serve
