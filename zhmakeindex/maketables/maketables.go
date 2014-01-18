// $Id$

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"unicode"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	outdir := flag.String("d", ".", "输出目录")
	flag.Parse()
	make_stroke_table(*outdir)
	make_reading_table(*outdir)
}

func make_stroke_table(outdir string) {
	const MAX_CODEPOINT = 0x40000 // 覆盖 Unicode 第 0、1、2、3 平面
	var CJKstrokes [MAX_CODEPOINT][]byte
	var maxStroke int = 0
	var unicodeVersion string
	// 使用海峰五笔码表数据，生成笔顺表
	sunwb_file, err := os.Open("sunwb_strokeorder.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer sunwb_file.Close()
	scanner := bufio.NewScanner(sunwb_file)
	for i := 1; scanner.Scan(); i++ {
		if scanner.Err() != nil {
			log.Fatalln(scanner.Err())
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, "\t")
		if len(fields) != 2 ||
			len([]rune(fields[0])) != 1 ||
			strings.IndexFunc(fields[1], isNotDigit) != -1 {
			log.Printf("笔顺文件第 %d 行语法错误，忽略。\n", i)
			continue
		}
		var r rune = []rune(fields[0])[0]
		var order []byte
		for _, rdigit := range fields[1] {
			digit, _ := strconv.ParseInt(string(rdigit), 10, 8)
			order = append(order, byte(digit))
		}
		CJKstrokes[r] = order
		if len(order) > maxStroke {
			maxStroke = len(order)
		}
	}
	// 使用 Unihan 数据库，读取笔画数补全其他字符
	unihan_file, err := os.Open("Unihan_DictionaryLikeData.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer unihan_file.Close()
	scanner = bufio.NewScanner(unihan_file)
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatalln(scanner.Err())
		}
		line := scanner.Text()
		if strings.Contains(line, "Unicode version:") {
			unicodeVersion = strings.TrimPrefix(line, "# ")
		}
		if strings.HasPrefix(line, "U+") && strings.Contains(line, "kTotalStrokes") {
			fields := strings.Split(line, "\t")
			var r rune
			fmt.Sscanf(fields[0], "U+%X", &r)
			var stroke int
			fmt.Sscanf(fields[2], "%d", &stroke)
			if CJKstrokes[r] != nil { // 笔顺数据已有，检查一致性
				if stroke != len(CJKstrokes[r]) {
					log.Printf("U+%04X (%c) 的笔顺数据（%d 画）与 unihan 笔画数（%d 画）不一致，跳过 unihan 数据\n",
						r, r, len(CJKstrokes[r]), stroke)
				}
			} else { // 无笔顺数据，假定每个笔画都是 6 号（未知）
				var order = make([]byte, stroke)
				for i := range order {
					order[i] = 6
				}
				CJKstrokes[r] = order
				if stroke > maxStroke {
					maxStroke = stroke
				}
			}
		}
	}
	// 输出笔顺表
	outfile, err := os.Create(path.Join(outdir, "strokes.go"))
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()
	fmt.Fprintln(outfile, `// 这是由程序自动生成的文件，请不要直接编辑此文件`)
	fmt.Fprintln(outfile, `// 笔顺来源：sunwb_strokeorder.txt`)
	fmt.Fprintln(outfile, `// 笔画数来源：Unihan_DictionaryLikeData.txt`)
	fmt.Fprintf(outfile, "// Unicode 版本：%s\n", unicodeVersion)
	fmt.Fprintln(outfile, `package main`)
	fmt.Fprintln(outfile, `var CJKstrokes = map[rune]string{`)
	for r, order := range CJKstrokes {
		if order == nil {
			continue
		}
		fmt.Fprintf(outfile, "\t%#x: \"", r)
		for _, s := range order {
			fmt.Fprintf(outfile, "\\x%02x", s)
		}
		fmt.Fprintf(outfile, "\", // %c\n", r)
	}
	fmt.Fprintln(outfile, `}`)
	fmt.Fprintf(outfile, "\nconst MAX_STROKE = %d\n", maxStroke)
}

func isNotDigit(r rune) bool {
	return !unicode.IsDigit(r)
}

func make_reading_table(outdir string) {
	// 读取 Unihan 读音表
	reading_table := make(map[rune]*ReadingEntry)
	reading_file, err := os.Open("Unihan_Readings.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer reading_file.Close()
	scanner := bufio.NewScanner(reading_file)
	largest := rune(0)
	var version string
	for scanner.Scan() {
		if scanner.Err() != nil {
			log.Fatalln(scanner.Err())
		}
		line := scanner.Text()
		if strings.Contains(line, "Unicode version:") {
			version = strings.TrimPrefix(line, "# ")
		}
		if strings.HasPrefix(line, "U+") {
			fields := strings.Split(line, "\t")
			var r rune
			fmt.Sscanf(fields[0], "U+%X", &r)
			if reading_table[r] == nil {
				reading_table[r] = &ReadingEntry{}
			}
			switch fields[1] {
			case "kHanyuPinlu":
				reading_table[r].HanyuPinlu = fields[2]
			case "kHanyuPinyin":
				reading_table[r].HanyuPinyin = fields[2]
			case "kMandarin":
				reading_table[r].Mandarin = fields[2]
			case "kXHC1983":
				reading_table[r].XHC1983 = fields[2]
			}
			if r > largest {
				largest = r
			}
		}
	}
	// 整理所有汉字的拼音表
	out_reading_table := make([]string, largest+1)
	for k, v := range reading_table {
		pinyin := v.regular()
		numbered := NumberedPinyin(pinyin)
		out_reading_table[k] = numbered
	}
	// 单独增加数字“〇”的读音
	if out_reading_table['〇'] == "" {
		out_reading_table['〇'] = "ling2"
	}
	// 输出
	outfile, err := os.Create(path.Join(outdir, "readings.go"))
	if err != nil {
		log.Fatalln(err)
	}
	defer outfile.Close()
	fmt.Fprintln(outfile, `// 这是由程序自动生成的文件，请不要直接编辑此文件`)
	fmt.Fprintln(outfile, `// 来源：Unihan_Readings.txt`)
	fmt.Fprintln(outfile, `//`, version)
	fmt.Fprintln(outfile, `package main`)
	fmt.Fprintln(outfile, `var CJKreadings = map[rune]string{`)
	for k, v := range out_reading_table {
		if v != "" {
			fmt.Fprintf(outfile, "\t%#x: %s, // %c\n", k, strconv.Quote(v), k)
		}
	}
	fmt.Fprintln(outfile, `}`)
}

type ReadingEntry struct {
	HanyuPinlu  string
	HanyuPinyin string
	Mandarin    string
	XHC1983     string
}

// 取出最常用的一个拼音
// 按如下优先次序：HanyuPinlu -> Mandarin -> XHC1983 -> HanyuPinyin
func (entry *ReadingEntry) regular() string {
	// xHanyuPinlu Syntax: [a-z\x{300}-\x{302}\x{304}\x{308}\x{30C}]+\([0-9]+\)
	// 如 cān(525) shēn(25)
	if entry.HanyuPinlu != "" {
		// 第一个括号之前的部分即可
		return strings.Split(entry.HanyuPinlu, "(")[0]
	}
	// kMandarin Syntax: [a-z\x{300}-\x{302}\x{304}\x{308}\x{30C}]+
	// 如 lüè
	if entry.Mandarin != "" {
		// 目前文件中没有多值情况，不过按 UAX #38 允许多值
		return strings.Split(entry.Mandarin, " ")[0]
	}
	// kXHC1983 Syntax: [0-9]{4}\.[0-9]{3}\*?(,[0-9]{4}\.[0-9]{3}\*?)*:[a-z\x{300}\x{301}\x{304}\x{308}\x{30C}]+
	// 如 1327.041:yán 1333.051:yàn
	if entry.XHC1983 != "" {
		// 第一项中第一个引号后的部分
		b := strings.Index(entry.XHC1983, ":")
		e := strings.Index(entry.XHC1983, " ")
		if e > 0 {
			return entry.XHC1983[b+1 : e]
		} else {
			return entry.XHC1983[b+1:]
		}
	}
	// kHanyuPinyin Syntax: (\d{5}\.\d{2}0,)*\d{5}\.\d{2}0:([a-z\x{300}-\x{302}\x{304}\x{308}\x{30C}]+,)*[a-z\x{300}-\x{302}\x{304}\x{308}\x{30C}]+
	// 如 10093.130:xī,lǔ 74609.020:lǔ,xī
	if entry.HanyuPinyin != "" {
		// 第一个冒号后，逗号/空格或词尾前的部分
		b := strings.Index(entry.HanyuPinyin, ":")
		e := strings.IndexAny(entry.HanyuPinyin[b:], " ,")
		if e > 0 {
			return entry.HanyuPinyin[b+1 : b+e]
		} else {
			return entry.HanyuPinyin[b+1:]
		}
	}
	// 没有汉语读音
	return ""
}

// 把拼音转换为无声调的拼音加数字声调
// 其中 ü 变为 v，轻声调号为 5，如 lǎo 转换为 lao3，lǘ 转换为 lv2
func NumberedPinyin(pinyin string) string {
	if pinyin == "" {
		return ""
	}
	numbered := []rune{}
	tone := 5
	for _, r := range pinyin {
		if Vowel[r] == 0 {
			numbered = append(numbered, r)
		} else {
			numbered = append(numbered, Vowel[r])
		}
		if Tones[r] != 0 {
			tone = Tones[r]
		}
	}
	numbered = append(numbered, []rune(strconv.Itoa(tone))...)
	return string(numbered)
}

var Vowel = map[rune]rune{
	'ā': 'a', 'á': 'a', 'ǎ': 'a', 'à': 'a',
	'ō': 'o', 'ó': 'o', 'ǒ': 'o', 'ò': 'o',
	'ē': 'e', 'é': 'e', 'ě': 'e', 'è': 'e',
	'ī': 'i', 'í': 'i', 'ǐ': 'i', 'ì': 'i',
	'ū': 'u', 'ú': 'u', 'ǔ': 'u', 'ù': 'u',
	'ǖ': 'v', 'ǘ': 'v', 'ǚ': 'v', 'ǜ': 'v', 'ü': 'v',
	'ń': 'n', 'ň': 'n', 'ǹ': 'n',
}

var Tones = map[rune]int{
	'ā': 1, 'ō': 1, 'ē': 1, 'ī': 1, 'ū': 1, 'ǖ': 1,
	'á': 2, 'ó': 2, 'é': 2, 'í': 2, 'ú': 2, 'ǘ': 2, 'ń': 2,
	'ǎ': 3, 'ǒ': 3, 'ě': 3, 'ǐ': 3, 'ǔ': 3, 'ǚ': 3, 'ň': 3,
	'à': 4, 'ò': 4, 'è': 4, 'ì': 4, 'ù': 4, 'ǜ': 4, 'ǹ': 4,
	'ü': 5,
}
