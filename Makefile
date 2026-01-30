BINARY_NAME=gemini-wm-remove
INSTALL_DIR=$(HOME)/.local/bin

.PHONY: all build install clean

all: build

build:
	go build -o $(BINARY_NAME)

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)
	chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed to $(INSTALL_DIR)/$(BINARY_NAME)"

clean:
	rm -f $(BINARY_NAME)
