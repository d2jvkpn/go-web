build:
	bash scripts/build_go-web.sh

run:
	bash scripts/build_main.sh
	./target/main

go-web:
	bash scripts/build_go-web.sh
	./target/go-web serve
