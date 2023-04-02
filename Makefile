.phony: README
README:
	cat fragments/README.md > README.md
	tools/genexampledocs.pl >> README.md
	tools/genbuiltinlist.pl >> README.md
