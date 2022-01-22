hello:
	echo "hello"

build:
	go mod vendor
	go build -race -o dongel ./

clean:
	rm -rf dongel
	echo "" > logs/app.log

dep:
	go mod vendor