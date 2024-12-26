say_hello:
	@echo "Hello. I'm a make file. I'm not sure why I'm here."

build:
	go build -o dist/TimecodeTool ./cmd/TimecodeTool/main.go

build_wasm:
	GOOS=js GOARCH=wasm go build -o dist/timecodetool.wasm ./cmd/wasm/main.go
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" dist/
	cp cmd/wasm/index.html dist/

build_wasm_tinygo:
	tinygo build -o dist/timecodetool_tiny.wasm -target wasm ./cmd/wasm/main.go
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" dist/
	cp cmd/wasm/index.html dist/

test:
	@go test ./internal -v

build_docs:
	# this is for github actions
	ifeq ($(wildcard dist/TimecodeTool-linux-amd64),)
		$(MAKE) build
	else
		# Proceed with your commands if file exists
		echo "File exists"
	endif

	@./dist/TimecodeTool gendocs web/docs/

	cd web/docs/ && for file in *.md; do \
		pandoc "$$file" -o "$${file%.md}.html" "--template=/Users/marcleonard/Projects/TimecodeTool/web/templates/_template.html"; \
	done

	cd web/docs/ && for file in *.html; do \
        sed -i '' 's/<h2 id="\([^"]*\)">/<h2 id="\1" class="major">/g' "$$file"; \
    	sed -i '' 's/\(<a href="\)\([a-zA-Z0-9_-]*\)\.md">/\1\2.html">/g' "$$file"; \
    done

	# This creates a simple entry point at /docs
	cp web/docs/TimecodeTool.html web/docs/index.html
