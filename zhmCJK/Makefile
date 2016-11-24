AUXFILES = *.dvi *.aux *.log *.glo *.gls *.idx *.ind *.ilg *.out *.hd *~
MAP = zhmCJK.map texfonts.map.template
SOURCE = zhmCJK.dtx zhmCJK.ins zhmCJK.lua zhmCJK-addmap.pl Makefile

.PHONY: all
all: zhmcjk.zip

zhmcjk.zip: $(SOURCE) zhmCJK.sty $(MAP) zhmCJK.tfm doc zhmCJK-addmap.pl miktex-tfm.tar.bz2 zhmcjk.tds.zip
	rm -rf zhmcjk
	mkdir zhmcjk
	cp README zhmcjk/
	mkdir zhmcjk/source
	cp $(SOURCE) zhmcjk/source/
	mkdir zhmcjk/tex
	cp zhmCJK.sty zhmcjk/tex/
	mkdir zhmcjk/fontname
	cp $(MAP) zhmcjk/fontname/
	mkdir zhmcjk/tfm
	cp zhmcjk.tfm zhmcjk/tfm/
	mkdir zhmcjk/doc
	cp zhmCJK.pdf zhmCJK-test.tex zhmcjk/doc/
	cp miktex-tfm.tar.bz2 zhmcjk/
	rm -f zhmcjk.zip
	zip -r zhmcjk.zip zhmcjk/
	zip zhmcjk.zip zhmcjk.tds.zip
	rm -r zhmcjk

zhmcjk.tds.zip: $(SOURCE) zhmCJK.sty $(MAP) zhmCJK.tfm doc
	rm -rf tds
	mkdir -p tds/source/latex/zhmcjk
	cp $(SOURCE) tds/source/latex/zhmcjk/
	mkdir -p tds/tex/latex/zhmcjk
	cp zhmCJK.sty tds/tex/latex/zhmcjk/
	mkdir -p tds/doc/latex/zhmcjk
	cp zhmCJK.pdf zhmCJK-test.tex README tds/doc/latex/zhmcjk/
	mkdir -p tds/fonts/map/fontname
	cp $(MAP) tds/fonts/map/fontname/
	mkdir -p tds/fonts/tfm/zhmcjk
	cp zhmCJK.tfm tds/fonts/tfm/zhmcjk/
	rm -f zhmcjk.tds.zip
	cd tds && zip -r ../zhmcjk.tds.zip source tex doc fonts
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
	rm -f zhmCJK.sty zhmCJK-test.tex README zhmCJK.pdf zhmCJK.tfm miktex-tfm.tar.bz2 $(MAP) zhmcjk.zip zhmcjk.tds.zip
