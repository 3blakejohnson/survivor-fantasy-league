.PHONY: run templ-generate

run: templ-generate
	go run main.go

templ-generate:
	templ generate