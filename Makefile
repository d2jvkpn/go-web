build:
	bash scripts/go_build2.sh

run:
	bash scripts/go_build.sh
	./target/main

go-web:
	bash scripts/go_build2.sh
	./target/go-web serve
