if exist diagbox.sty del diagbox.sty
if exist diagbox.pdf del diagbox.pdf
pdftex diagbox.ins
xelatex diagbox.dtx
zhmakeindex -s gind.ist diagbox.idx
zhmakeindex -s gglo.ist -o diagbox.gls diagbox.glo
xelatex diagbox.dtx
xelatex diagbox.dtx

mkdir tds\source\latex\diagbox
mkdir tds\tex\latex\diagbox
mkdir tds\doc\latex\diagbox
copy diagbox.dtx tds\source\latex\diagbox
copy diagbox.ins tds\source\latex\diagbox
copy diagbox.sty tds\tex\latex\diagbox
copy diagbox.pdf tds\doc\latex\diagbox
copy README tds\doc\latex\diagbox
if exist diagbox.tds.zip del diagbox.tds.zip
cd tds
zip -r ..\diagbox.tds.zip .
cd ..
rd /s /q tds

mkdir diagbox
for %%i in (diagbox.dtx diagbox.ins diagbox.pdf README) do (
	copy %%i diagbox\
)
if exist diagbox.zip del diagbox.zip
zip -r diagbox.zip diagbox\
rd /s /q diagbox

del diagbox.sty diagbox.pdf
del diagbox.aux diagbox.log diagbox.glo diagbox.gls diagbox.idx diagbox.ind diagbox.ilg diagbox.out diagbox.tmp diagbox.hd diagbox.*~
