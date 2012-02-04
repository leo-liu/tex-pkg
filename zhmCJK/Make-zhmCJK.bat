latex zhmCJK.ins
latex zhmCJK.dtx
makeindex -s gind zhmCJK.idx
makeindex -s gglo -o zhmCJK.gls zhmCJK.glo
latex zhmCJK.dtx
latex zhmCJK.dtx
dvipdfmx zhmCJK.dvi
for %%i in (zhmCJK.dvi zhmCJK.aux zhmCJK.log zhmCJK.glo zhmCJK.gls zhmCJK.idx zhmCJK.ind zhmCJK.ilg zhmCJK.out zhmCJK.tmp zhmCJK.hd zhmCJK.*~ Make-zhmCJK.bat~) do del %%i

if "%1"=="zip" call :zip
if "%1"=="tds" call :tds
shift
if "%1"=="zip" call :zip
if "%1"=="tds" call :tds
exit /b

:zip
if exist zhmCJK.zip del zhmCJK.zip
zip zhmCJK zhmCJK.dtx zhmCJK.ins README.txt zhmCJK.sty zhmCJK.pdf zhmetrics.tfm texfonts.map zhmCJK.map zhmCJK-test.tex
exit /b

:tds
if exist zhmCJK-tds.zip del zhmCJK-tds.zip
for %%i in (doc fonts source tex) do if exist %%i rmdir /s %%i
mkdir source\latex\zhmCJK
copy zhmCJK.dtx source\latex\zhmCJK\
copy zhmCJK.ins source\latex\zhmCJK\
mkdir tex\latex\zhmCJK
copy zhmCJK.sty tex\latex\zhmCJK\
mkdir doc\latex\zhmCJK
copy zhmCJK.pdf doc\latex\zhmCJK\
copy zhmCJK-test.tex doc\latex\zhmCJK\
copy README.txt doc\latex\zhmCJK\
mkdir fonts\map\fontname
copy zhmCJK.map fonts\map\fontname\
copy texfonts.map fonts\map\fontname\
mkdir fonts\tfm\zhmetrics
copy zhmetrics.tfm fonts\tfm\zhmetrics\
zip -r zhmCJK-tds source tex doc fonts
exit /b
