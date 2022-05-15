# APP_GitForce=true
version:
	bash scripts/build.sh && mv ./target/goapp main
	./main version

serve:
	bash scripts/build.sh && mv ./target/goapp main
	./main serve
