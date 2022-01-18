hello:
	echo "hello"

build:
	go build -race -o dongel ./

clean:
	rm -rf dongel
	echo "" > logs/app.log