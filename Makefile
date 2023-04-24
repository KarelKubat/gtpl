what:
	@echo
	@echo "Make what?"
	@echo "  make install   to install it (or just run `go install gtpl.go`)"
	@echo "  make README    to refresh docs (for maintainers)"
	@echo "  make newmod    to refresh go.mod and go.sum (for maintainers)"
	@echo "  make all       for all of the above"
	@exit 1

all:
	make newmod
	tools/checks/gotests
	make README
	make install

# Run `make install` to install it.
install:
	go install gtpl.go

# Just for the generation of up-to-date docs.
# fragments/hosts.md how has copy/pasted sources from examples/hosts/* and the output of the
# corresponding `gtpl` commands. I might change that into generated info.
README:
	cat fragments/README.md    >  /tmp/gtpl.README.md
	echo                       >> /tmp/gtpl.README.md
	tools/genexampledocs.pl    >> /tmp/gtpl.README.md
	echo                       >> /tmp/gtpl.README.md
	cat fragments/hosts.md     >> /tmp/gtpl.README.md
	echo                       >> /tmp/gtpl.README.md
	tools/genbuiltinlist.pl    >> /tmp/gtpl.README.md
	echo                       >> /tmp/gtpl.README.md
	cat fragments/embedding.md >> /tmp/gtpl.README.md
	echo                       >> /tmp/gtpl.README.md
	cp /tmp/gtpl.README.md README.md
	rm /tmp/gtpl.README.md
	-mdtoc --inplace README.md

newmod:
	rm -f go.mod go.sum
	go mod init
	go mod tidy