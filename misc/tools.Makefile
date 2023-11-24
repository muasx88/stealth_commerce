# This makefile should be used to hold functions/variables

define github_url
    https://github.com/$(GITHUB)/releases/download/v$(VERSION)/$(ARCHIVE)
endef

# creates a directory bin.
bin:
	@ mkdir -p $@


# ~~ [ air ] ~~~ https://github.com/cosmtrek/air ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

AIR := $(shell command -v air || echo "bin/air")
air: bin/air ## Installs air (go file watcher)

bin/air: VERSION := 1.15.1
bin/air: GITHUB  := cosmtrek/air
bin/air: ARCHIVE := air_$(VERSION)_$(OSTYPE)_amd64.tar.gz
bin/air: bin
	@ printf "Install air... "
	@ curl -Ls $(shell echo $(call github_url) | tr A-Z a-z) | tar -zOxf - air > $@ && chmod +x $@
	@ echo "done."