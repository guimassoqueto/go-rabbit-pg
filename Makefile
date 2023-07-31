up:
	docker compose up -d

b:
	rm out/main && go build -o out/main && ./out/main

r:
	go run .

or:
	open https://github.com/guimassoqueto/go-rabbit-pig

a:
	rm out/main && go build -o out/main && ./out/main
