# Go guilds and moves the executable to the /usr/bin directory
# so that it can be run from anywhere in the terminal.
# Linux only currently, if you want to use this on Windows then please look at yourself in the mirror and ask yourself why you're using Windows.

EXECUTABLE=modcheck
INSTALL_DIR=/usr/bin

build:
	go build -o $(EXECUTABLE)

install:
	go build -o $(EXECUTABLE)
	sudo mv $(EXECUTABLE) $(INSTALL_DIR)

uninstall:
	sudo rm $(INSTALL_DIR)/$(EXECUTABLE)