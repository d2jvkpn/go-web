build:
	bash scripts/go-web.sh

run:
	bash scripts/go_build.sh
	./target/main

go-web:
	bash scripts/go-web.sh
	./target/go-web serve
