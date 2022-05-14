version:
	APP_GitForce=true bash build.sh
	./target/goapp version

serve:
	APP_GitForce=true bash build.sh
	./target/goapp serve
