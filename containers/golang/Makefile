# This makefile compiles the server for testing purposes
# Please use docker-compose in production

EXEC_NAME=server
SRC=./src

all: build

build:
	go build -o $(EXEC_NAME) $(SRC) CGO_ENABLED=0

clean:

fclean: clean
	rm -rf $(EXEC_NAME)

re: fclean build
