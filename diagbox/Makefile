AUXFILES = *.aux *.log *.glo *.gls *.idx *.ind *.ilg *.out *.tmp *.hd *~
SOURCE = diagbox.dtx diagbox.ins README

.PHONY: all
all: diagbox.zip

diagbox.zip: $(SOURCE) diagbox.pdf
	rm -rf diagbox
	mkdir diagbox
	cp $(SOURCE) diagbox.pdf diagbox/
	zip -r diagbox.zip diagbox/
	rm -r diagbox

diagbox.sty: $(SOURCE)
	pdftex diagbox.ins

.PHONY: doc
doc: diagbox.pdf

diagbox.pdf: $(SOURCE)
	xelatex diagbox.dtx
	zhmakeindex -s gind.ist diagbox.idx
	zhmakeindex -s gglo.ist -o diagbox.gls diagbox.glo
	xelatex diagbox.dtx
	xelatex diagbox.dtx

.PHONY: clean
clean:
	rm -f $(AUXFILES)
	rm -rf diagbox
	rm -f diagbox.sty diagbox.pdf diagbox.zip
