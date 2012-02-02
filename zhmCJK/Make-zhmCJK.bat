if exist zhmCJK.sty del zhmCJK.sty
if exist zhmCJK.pdf del zhmCJK.pdf
tex zhmCJK.ins
latex zhmCJK.dtx
makeindex -s gind zhmCJK.idx
makeindex -s gglo -o zhmCJK.gls zhmCJK.glo
latex zhmCJK.dtx
latex zhmCJK.dtx
dvipdfmx zhmCJK.dvi
if exist zhmCJK.zip del zhmCJK.zip
zip zhmCJK zhmCJK.dtx zhmCJK.ins zhmCJK.sty zhmCJK.pdf zhmetrics.tfm texfonts.map zhmCJK.map test-zhmCJK.tex
del zhmCJK.dvi zhmCJK.aux zhmCJK.log zhmCJK.glo zhmCJK.gls zhmCJK.idx zhmCJK.ind zhmCJK.ilg zhmCJK.out zhmCJK.tmp zhmCJK.hd zhmCJK.*~
