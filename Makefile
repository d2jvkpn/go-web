# APP_GitForce=true
version:
	bash build.sh && mv ./target/goapp main
	./main version

serve:
	bash build.sh && mv ./target/goapp main
	./main serve
