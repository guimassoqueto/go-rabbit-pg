i:
	go get

or:
	open https://github.com/guimassoqueto/scrap-colly

env:
	cp .env.sample .env
	 
a:
	echo 123 1> out/main && rm out/main && go build -o out/main && ./out/main
