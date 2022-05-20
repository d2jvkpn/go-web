# APP_GitForce=true
version:
	bash scripts/build.sh
	./target/goapp version

serve:
	bash scripts/build.sh
	./target/goapp serve

run:
	go build -o main main.go
	./main serve
