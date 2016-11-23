AUXFILES = *.dvi *.aux *.log *.glo *.gls *.idx *.ind *.ilg *.out *.hd *~
MAP = zhmCJK.map texfonts.map.template
SOURCE = zhmCJK.dtx zhmCJK.ins zhmCJK.lua zhmCJK-addmap.pl Makefile

.PHONY: all
all: zhmCJK.zip

zhmCJK.zip: $(SOURCE) zhmCJK.sty $(MAP) zhmCJK.tfm doc zhmCJK-addmap.pl miktex-tfm.tar.bz2 zhmCJK.tds.zip
	rm -rf zhmCJK
	mkdir zhmCJK
	cp README zhmCJK/
	mkdir zhmCJK/source
	cp $(SOURCE) zhmCJK/source/
	mkdir zhmCJK/tex
	cp zhmCJK.sty zhmCJK/tex/
	mkdir zhmCJK/fontname
	cp $(MAP) zhmCJK/fontname/
	mkdir zhmCJK/tfm
	cp zhmCJK.tfm zhmCJK/tfm/
	mkdir zhmCJK/doc
	cp zhmCJK.pdf zhmCJK-test.tex zhmCJK/doc/
	cp miktex-tfm.tar.bz2 zhmCJK/
	rm -f zhmCJK.zip
	zip -r zhmCJK.zip zhmCJK/
	zip zhmCJK.zip zhmCJK.tds.zip
	rm -r zhmCJK

zhmCJK.tds.zip: $(SOURCE) zhmCJK.sty $(MAP) zhmCJK.tfm doc
	rm -rf tds
	mkdir -p tds/source/latex/zhmCJK
	cp $(SOURCE) tds/source/latex/zhmCJK/
	mkdir -p tds/tex/latex/zhmCJK
	cp zhmCJK.sty tds/tex/latex/zhmCJK/
	mkdir -p tds/doc/latex/zhmCJK
	cp zhmCJK.pdf zhmCJK-test.tex README tds/doc/latex/zhmCJK/
	mkdir -p tds/fonts/map/fontname
	cp $(MAP) tds/fonts/map/fontname/
	mkdir -p tds/fonts/tfm/zhmCJK
	cp zhmCJK.tfm tds/fonts/tfm/zhmCJK/
	rm -f zhmCJK.tds.zip
	cd tds && zip -r ../zhmCJK.tds.zip source tex doc fonts
	rm -r tds

.PHONY: doc
doc: zhmCJK.pdf zhmCJK-test.tex README

zhmCJK.sty zhmCJK-test.tex README: zhmCJK.dtx zhmCJK.ins
	pdftex zhmCJK.ins
	mv README.txt README

$(MAP) zhmCJK.tfm: zhmCJK.lua
	texlua zhmCJK.lua map

miktex-tfm.tar.bz2: zhmCJK.lua
	texlua zhmCJK.lua nomap
	rm -f miktex-tfm.tar.bz2
	tar cfj miktex-tfm.tar.bz2 miktex-tfm/
	rm -r miktex-tfm/

zhmCJK.pdf: zhmCJK.sty $(MAP)
	latex zhmCJK.dtx
	makeindex -s gind zhmCJK.idx
	makeindex -s gglo -o zhmCJK.gls zhmCJK.glo
	latex zhmCJK.dtx
	latex zhmCJK.dtx
	dvipdfmx zhmCJK.dvi
	rm -f $(AUXFILES)

.PHONY: clean
clean:
	rm -f $(AUXFILES)
	rm -f zhmCJK.sty zhmCJK-test.tex README zhmCJK.pdf zhmCJK.tfm miktex-tfm.tar.bz2 $(MAP) zhmCJK.zip zhmCJK.tds.zip
