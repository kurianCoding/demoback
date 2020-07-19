run: bin/herokuappkurian
	@PATH="$(PWD)/bin:$(PATH)" heroku local

bin/herokuappkurian: main.go
	go build -o bin/herokuappkurian main.go

clean:
	rm -rf bin
