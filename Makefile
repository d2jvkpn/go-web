build:
	bash scripts/build_main.sh

main:
	bash scripts/build_main.sh
	./target/main


build_go-web:
	bash scripts/build_go-web.sh

go-web:
	bash scripts/build_go-web.sh
	./target/go-web serve
