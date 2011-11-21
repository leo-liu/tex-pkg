if exist diagbox.sty del diagbox.sty
if exist diagbox.pdf del diagbox.pdf
tex diagbox.ins
xelatex diagbox.dtx
makeindex -s gind diagbox.idx
makeindex -s gglo -o diagbox.gls diagbox.glo
xelatex diagbox.dtx
xelatex diagbox.dtx
if exist diagbox.zip del diagbox.zip
zip diagbox diagbox.dtx diagbox.ins diagbox.sty diagbox.pdf
del diagbox.aux diagbox.log diagbox.glo diagbox.gls diagbox.idx diagbox.ind diagbox.ilg diagbox.out diagbox.tmp  diagbox.*~