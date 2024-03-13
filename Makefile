init:
	rm -f go.mod
	go mod init github.com/gdelrieuadevinta/layerzero-takehome
	go mod tidy

build: 
	go build -o price_checker

run: build
	./price_checker

docker_build:
	docker build -t currency-price-checker .

kind_load_image:
	kind load docker-image currency-price-checker

install_nginx:
	kind load docker-image currency-price-checker