#
# cross-compile helper.
#
# Install Go with cross-compile enabled:
#    brew install go --cross-cc-all
#
# or, re-install Go with cross-compile enabled:
#    brew reinstall go --cross-cc-all
#
# Then, install gox.
#    go get github.com/mitchellh/gox

.PHONY : clean osx linux

all: clean osx linux

clean:
	find . -type f -perm +111 -name 'indefatigable*' | xargs rm

osx:
	go build

linux:
	gox -osarch="linux/amd64"
