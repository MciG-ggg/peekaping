.PHONY: dev-watch doc-watch gen-m

dev-watch:
	nodemon --watch './**/*' -e go,html --signal SIGTERM --exec "go run ./src"

doc-watch:
	nodemon --watch './**/*' --ignore docs -e go --exec "swag init -d src"

gen-m:
	go run scripts/generate_module.go
