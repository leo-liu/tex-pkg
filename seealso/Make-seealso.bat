if exist seealso.sty del seealso.sty
if exist seealso.pdf del seealso.pdf
pdftex seealso.ins
xelatex seealso.dtx
makeindex -s gind seealso.idx
makeindex -s gglo -o seealso.gls seealso.glo
xelatex seealso.dtx
xelatex seealso.dtx
if exist seealso.zip del seealso.zip
zip seealso seealso.dtx seealso.ins seealso.sty seealso.pdf
del seealso.aux seealso.log seealso.glo seealso.gls seealso.idx seealso.ind seealso.ilg seealso.out seealso.tmp seealso.hd  seealso.*~
