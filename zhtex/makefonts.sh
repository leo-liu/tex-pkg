#!/usr/bin/env bash

# This is file 'makefonts.sh', a very simple bash script to generate TFM files
# and Type1 fonts from TrueType fonts.
#
# This script is for personal use only. It provides a minimal installation
# of Chinese TeX fonts for zhtex.tex. Only Unicode encoding is supported. If
# more features (e.g. LaTeX NFSS) are needed, it is preferred to use a more
# complicated script such as 'CTeXFonts.lua'.
#
# Requirements:
#  * Bash
#  * ttf2tfm (available in TeX Live)
#  * otfinfo (available in TeX Live)
#  * ttf2pt1 (available in sourceforge.net)
#
# Usage:
#  1. Put the TrueType file "foo.ttf" in current directory
#  2. Run
#       ./makefonts.sh foo
#  3. Copy the generated files to TDS tree
#  4. Run updmap with proper options, for TeX Live:
#       updmap --enable Map=foo-pt1.map
#  5. Happy TeXing.
#
#
# Copyright (C) 2016 by Leo Liu <leoliu.pku@gmail.com>

for i in tfm map type1
do
	if [ ! -e $i ]
	then
		mkdir $i
	fi
done

function convert_ttf {
	local fontname=$1
	local psname=`otfinfo -p $fontname.ttf`
	local familyname=`otfinfo -a $fontname.ttf`

	for i in tfm map type1
	do
		if [ ! -e $i/$fontname ]
		then
			mkdir $i/$fontname
		fi
	done

	echo %% Map file for TrueType font ${fontname}.ttf > map/${fontname}/${fontname}-ttf.map
	echo %% Family name: ${fontname} >> map/${fontname}/${fontname}-ttf.map
	ttf2tfm ${fontname}.ttf -q tfm/${fontname}/${fontname}@Unicode@.tfm >> map/${fontname}/${fontname}-ttf.map

	echo %% Map file for Type1 fonts from ${fontname}.ttf > map/${fontname}/${fontname}-pt1.map
	echo %% PostScript font name: $psname >> map/${fontname}/${fontname}-pt1.map
	echo %% Family name: $familyname >> map/${fontname}/${fontname}-pt1.map

	for i in 0 1 2 3 4 5 6 7 8 9 a b c d e f
	do
  		for j in 0 1 2 3 4 5 6 7 8 9 a b c d e f
		do
			if [ -e tfm/$fontname/$fontname$i$j.tfm ]
			then
				ttf2pt1 -b -GFae -l plane+0x$i$j -p ttf $fontname.ttf type1/$fontname/$fontname$i$j
				echo $fontname$i$j $psname-$i$j \<$fontname$i$j.pfb >> map/$fontname/$fontname-pt1.map
			fi
		done
	done
}

for f
do
	convert_ttf $f
done
