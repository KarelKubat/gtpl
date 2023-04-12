.phony: all
all:
	tools/checks/gotests
	make README
	make install

# Run `make install` to install it.
.phony: install
install:
	go install gtpl.go

# Just for the generation of up-to-date docs.
.phony: README
README:
	cat fragments/README.md    >  /tmp/gtpl.README.md
	echo                       >> /tmp/gtpl.README.md
	tools/genexampledocs.pl    >> /tmp/gtpl.README.md
	echo                       >> /tmp/gtpl.README.md
	tools/genbuiltinlist.pl    >> /tmp/gtpl.README.md
	echo                       >> /tmp/gtpl.README.md
	cat fragments/embedding.md >> /tmp/gtpl.README.md
	echo                       >> /tmp/gtpl.README.md
	cp /tmp/gtpl.README.md README.md
	rm /tmp/gtpl.README.md
	-mdtoc --inplace README.md
