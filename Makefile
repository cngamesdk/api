PROJECT_PATH 		=	$(PWD)
BINARY_DIR			=	$(PROJECT_PATH)/build/
BINARY_NAME			=	http_server
BINARY_FULL_PATH	=	$(BINARY_DIR)$(BINARY_NAME)
CONFIG_PATH			=	$(PROJECT_PATH)/config/config.yaml

runNum				=	$(shell ps aux | grep $(BINARY_NAME) | grep -v grep | wc -l | sed -e 's/^[ \t]*//g')
runProcess			=	$(shell ps aux | grep $(BINARY_NAME) | grep -v grep | awk '{print $$2}')

.PHONY: build

clean:
	@if [ -f $(BINARY_FULL_PATH) ]; then rm $(BINARY_FULL_PATH); fi
	@echo "clean successful"

build:	clean
	@go build -ldflags '-w -s' -o $(BINARY_FULL_PATH) ./cmd/http/main.go
	@echo "build successful"

stop:
    ifeq ($(runNum), 0)
		@echo "no process is running."
    else
		@for x in $(runProcess); do kill $$x; done
		@echo "stop successful"
    endif

start:
     ifeq ($(runNum), 0)
		@if [ -f $(BINARY_FULL_PATH) ]; then $(BINARY_FULL_PATH) -c=$(CONFIG_PATH); else echo "execute file is not exists"; fi
	else
		@echo "server is running"
    endif
