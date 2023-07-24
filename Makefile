up:
	docker compose up -d

b:
	rm out/main && go build -o out/main && ./out/main

r:
	go run .