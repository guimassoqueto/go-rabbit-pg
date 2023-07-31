i:
	go get

or:
	open https://github.com/guimassoqueto/scrap-colly

env:
	cat .env.sample 1> .env
	 
a:
	rm out/main && go build -o out/main && ./out/main
