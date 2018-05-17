PROJECT_NAME = ExpertSystem

#FLAGS = -gcflags -m
FLAGS = 

all: $(PROJECT_NAME)

$(PROJECT_NAME):
	@echo "Building $(PROJECT_NAME).."
	@go build $(FLAGS) -o $(PROJECT_NAME) $(wildcard */*.go) && echo "Building done !" || echo "Building failed :("

deps:
	go get github.com/buger/jsonparser

install: deps

%.go:

re: fclean all

fclean:
	@rm -rf $(PROJECT_NAME)

.PHONY: re fclean deps $(PROJECT_NAME)
