AUXFILES = *.aux *.log *.glo *.gls *.idx *.ind *.ilg *.out *.tmp *.hd *~
SOURCE = seealso.dtx seealso.ins README

.PHONY: all
all: seealso.zip

seealso.zip: $(SOURCE) seealso.pdf
	rm -rf seealso
	mkdir seealso
	cp $(SOURCE) seealso.pdf seealso/
	zip -r seealso.zip seealso/
	rm -r seealso

seealso.sty: $(SOURCE)
	pdftex seealso.ins

.PHONY: doc
doc: seealso.pdf

seealso.pdf: $(SOURCE)
	xelatex seealso.dtx
	makeindex -s gind.ist seealso.idx
	makeindex -s gglo.ist -o seealso.gls seealso.glo
	xelatex seealso.dtx
	xelatex seealso.dtx

.PHONY: clean
clean:
	rm -f $(AUXFILES)
	rm -rf seealso
	rm -f seealso.sty seealso.pdf seealso.zip
