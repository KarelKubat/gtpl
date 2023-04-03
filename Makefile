.phony: all
all:
	make README
	make install

# Run `make install` to install it.
.phony: install
install:
	go install gtpl.go

# Just for the generation of up-to-date docs.
.phony: README
README:
	cat fragments/README.md > README.md
	tools/genexampledocs.pl >> README.md
	tools/genbuiltinlist.pl >> README.md
