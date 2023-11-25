.PHONY: run generate deploy

run:
	cd src && go run main.go

generate:
	cd src && go get github.com/99designs/gqlgen@v0.17.40 && go run github.com/99designs/gqlgen generate

deploy:
	cd aws && cdk deploy