@echo off
REM This is file 'makefonts.cmd', a very simple batch script to generate TFM
REM files and Type1 fonts from TrueType fonts.
REM
REM This script is for personal use only. It provides a minimal installation
REM of Chinese TeX fonts for zhtex.tex. Only Unicode encoding is supported. If
REM more features (e.g. LaTeX NFSS) are needed, it is preferred to use a more
REM complicated script such as 'CTeXFonts.lua'.
REM
REM Requirements:
REM  * Microsoft Windows
REM  * ttf2tfm (available in TeX Live / MiKTeX)
REM  * otfinfo (available in TeX Live)
REM  * ttf2pt1 (available in GnuWin32)
REM
REM Usage:
REM  1. Put the TrueType file "foo.ttf" in current directory
REM  2. Run
REM       makefonts foo
REM  3. Copy the generated files to TDS tree
REM  4. Run updmap with proper options, for TeX Live:
REM       updmap --enable Map=foo-pt1.map
REM  5. Happy TeXing.
REM
REM
REM Copyright (C) 2013 by Leo Liu <leoliu.pku@gmail.com>

setlocal

if not exist tfm mkdir tfm
if not exist map mkdir map
if not exist type1 mkdir type1

:loop

if "%1"=="" goto endloop
set fontname=%1
shift

for /f %%i in ('otfinfo -p %fontname%.ttf') do set psname=%%i
for /f %%i in ('otfinfo -a %fontname%.ttf') do set familyname=%%i

if not exist tfm\%fontname% mkdir tfm\%fontname%
if not exist map\%fontname% mkdir map\%fontname%
if not exist type1\%fontname% mkdir type1\%fontname%

echo %% Map file for TrueType font %fontname%.ttf > map\%fontname%\%fontname%-ttf.map
echo %% Family name: %familyname% >> map\%fontname%\%fontname%-ttf.map
ttf2tfm %fontname%.ttf -q tfm\%fontname%\%fontname%@Unicode@.tfm >> map\%fontname%\%fontname%-ttf.map

echo %% Map file for Type1 fonts from %fontname%.ttf > map\%fontname%\%fontname%-pt1.map
echo %% PostScript font name: %psname% >> map\%fontname%\%fontname%-pt1.map
echo %% Family name: %familyname% >> map\%fontname%\%fontname%-pt1.map

for %%i in (0,1,2,3,4,5,6,7,8,9,a,b,c,d,e,f) do (
  for %%j in (0,1,2,3,4,5,6,7,8,9,a,b,c,d,e,f) do (
    if exist tfm\%fontname%\%fontname%%%i%%j.tfm (
      ttf2pt1 -b -GFae -l plane+0x%%i%%j -p ttf %fontname%.ttf type1\%fontname%\%fontname%%%i%%j
      echo %fontname%%%i%%j %psname%-%%i%%j ^<%fontname%%%i%%j.pfb >> map\%fontname%\%fontname%-pt1.map
    )
  )
)

goto loop
:endloop

endlocal
