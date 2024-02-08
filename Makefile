GO_CMD = GOARCH=amd64 go
GO_CMD_ARM = GOARCH=arm64 go
GO_CMD_LINUX = GOOS=linux $(GO_CMD)
GO_CMD_DARWIN = GOOS=darwin $(GO_CMD)
GO_CMD_DARWIN_ARM = GOOS=darwin $(GO_CMD_ARM)
GO_CMD_WINDOWS = GOOS=windows $(GO_CMD)

GO_FILES = $(shell find ./ -type f -name '*.go')

.PHONY:all
all: ldap-utils-linux ldap-utils-macos ldap-utils-macos-arm ldap-utils-windows.exe

.PHONY:clean
clean:
	rm -f ldap-utils-linux
	rm -f ldap-utils-macos
	rm -f ldap-utils-macos-arm
	rm -f ldap-utils-windows.exe

ldap-utils-linux:$(GO_FILES)
	$(GO_CMD_LINUX) build -o $@ $^

ldap-utils-macos:$(GO_FILES)
	$(GO_CMD_DARWIN) build -o $@ $^

ldap-utils-macos-arm:$(GO_FILES)
	$(GO_CMD_DARWIN_ARM) build -o $@ $^

ldap-utils-windows.exe:$(GO_FILES)
	$(GO_CMD_WINDOWS) build -o $@ $^
