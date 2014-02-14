// $Id$

package main

import (
	"log"
	"sort"
	"strconv"
	"unicode"
)

// 对应不同的分类排序方式
type IndexCollator interface {
	InitGroups(style *OutputStyle) []IndexGroup
	Group(entry *IndexEntry) int
	RuneCmp(a, b rune) int
}

// 排序器
type IndexSorter struct {
	IndexCollator
}

func NewIndexSorter(method string) *IndexSorter {
	switch method {
	case "bihua", "stroke":
		return &IndexSorter{
			IndexCollator: StrokeIndexCollator{},
		}
	case "pinyin", "reading":
		return &IndexSorter{
			IndexCollator: ReadingIndexCollator{},
		}
	case "bushou", "radical":
		return &IndexSorter{
			IndexCollator: RadicalIndexCollator{},
		}
	default:
		log.Fatalln("未知排序方式")
	}
	return nil
}

func (sorter *IndexSorter) SortIndex(input *InputIndex, style *OutputStyle, option *OutputOptions) *OutputIndex {
	out := new(OutputIndex)
	// 分组
	out.groups = sorter.InitGroups(style)

	// 先整体排序
	sort.Sort(IndexEntrySlice{
		entries:  *input,
		colattor: sorter.IndexCollator,
	})

	// 再依次对页码排序，并分组添加
	pagesorter := NewPageSorter(style, option)
	for _, entry := range *input {
		pageranges := pagesorter.Sort(entry.pagelist)
		pageranges = pagesorter.Merge(pageranges)
		item := IndexItem{
			level: len(entry.level) - 1,
			text:  entry.level[len(entry.level)-1].text,
			page:  pageranges,
		}
		group := sorter.Group(&entry)
		out.groups[group].items = append(out.groups[group].items, item)
	}

	return out
}

type IndexEntrySlice struct {
	entries  []IndexEntry
	colattor IndexCollator
}

func (s IndexEntrySlice) Len() int {
	return len(s.entries)
}

func (s IndexEntrySlice) Swap(i, j int) {
	s.entries[i], s.entries[j] = s.entries[j], s.entries[i]
}

// 比较两个串的大小
func (s IndexEntrySlice) Strcmp(a, b string) int {
	// 先尝试按数字比较
	if cmp := DecimalStrcmp(a, b); cmp != 0 {
		return cmp
	}
	a_rune, b_rune := []rune(a), []rune(b)
	// 特例：在符号开头的串 < 数字开头的串，后面是字母汉字开头的串
	// Unicode 中字母汉字总在数字之后，只特别处理串首是符号或是数字的情形
//	var a0, b0 rune
//	if len(a_rune) > 0 {
//		a0 = a_rune[0]
//	}
//	if len(b_rune) > 0 {
//		b0 = b_rune[0]
//	}
//	if !IsNumRune(a0) && !unicode.IsLetter(a0) && IsNumRune(b0) {
//		return -1
//	} else if IsNumRune(a0) && !IsNumRune(b0) && !unicode.IsLetter(b0) {
//		return 1
//	}
	// 忽略大小写，按字典序比较
	for i := range a_rune {
		if i >= len(b_rune) {
			return 1
		}
		cmp := s.colattor.RuneCmp(a_rune[i], b_rune[i])
		if cmp != 0 {
			return cmp
		}
	}
	if len(a_rune) < len(b_rune) {
		return -1
	}
	// 不忽略大小写重新比较串
	if a < b {
		return -1
	} else if a > b {
		return 1
	} else {
		return 0
	}
}

func (s IndexEntrySlice) Less(i, j int) bool {
	a, b := s.entries[i], s.entries[j]
	for i := range a.level {
		if i >= len(b.level) {
			return false
		}
		keycmp := s.Strcmp(a.level[i].key, b.level[i].key)
		if keycmp < 0 {
			return true
		} else if keycmp > 0 {
			return false
		}
		textcmp := s.Strcmp(a.level[i].text, b.level[i].text)
		if textcmp < 0 {
			return true
		} else if textcmp > 0 {
			return false
		}
	}
	if len(a.level) < len(b.level) {
		return true
	}
	return false
}

// 页码排序器
type PageSorter struct {
	precedence    map[NumFormat]int
	strict        bool
	disable_range bool
}

func NewPageSorter(style *OutputStyle, option *OutputOptions) *PageSorter {
	var sorter PageSorter
	sorter.precedence = make(map[NumFormat]int)
	for i, r := range style.page_precedence {
		switch r {
		case 'r':
			sorter.precedence[NUM_ROMAN_LOWER] = i
		case 'n':
			sorter.precedence[NUM_ARABIC] = i
		case 'a':
			sorter.precedence[NUM_ALPH_LOWER] = i
		case 'R':
			sorter.precedence[NUM_ROMAN_UPPER] = i
		case 'A':
			sorter.precedence[NUM_ALPH_UPPER] = i
		default:
			log.Println("page_precedence 语法错误，采用默认值")
			sorter.precedence = map[NumFormat]int{
				NUM_ROMAN_LOWER: 0,
				NUM_ARABIC:      1,
				NUM_ALPH_LOWER:  2,
				NUM_ROMAN_UPPER: 3,
				NUM_ALPH_UPPER:  4,
			}
		}
	}
	sorter.strict = option.strict
	sorter.disable_range = option.disable_range
	return &sorter
}

// 处理输入的页码，生成页码区间组
func (sorter *PageSorter) Sort(pages []PageInput) []PageRange {
	//	debug.Println(pages)
	var out []PageRange
	// 合并前排序。传统 Makeindex 按原始输入的次序，在处理多个文件时可能不大好
	if sorter.strict {
		sort.Sort(PageInputSliceStrict{
			PageInputSlice{pages: pages, sorter: sorter}})
	} else {
		sort.Sort(PageInputSliceLoose{
			PageInputSlice{pages: pages, sorter: sorter}})
	}
	//debug.Println(pages)
	// 使用一个栈来合并页码区间
	// 这里的合并只将 1( 2 3 3) 合并为 1--3，不处理相邻区间，后者需要再做 Merge 操作
	var stack []PageInput
	for i := 0; i < len(pages); i++ {
		p := pages[i]
		//debug.Printf("处理页码 %s{%s} %s\n", p.encap, p.NumString(), p.rangetype)
		if len(stack) == 0 {
			switch p.rangetype {
			case PAGE_NORMAL:
				// 输出独立页
				out = append(out, PageRange{begin: p, end: p})
			case PAGE_OPEN:
				// 压栈
				stack = append(stack, p)
			case PAGE_CLOSE:
				log.Printf("页码区间有误，区间末尾 %s{%s} 没有匹配的区间头。\n", p.encap, p)
				// 输出从空白到当前页的伪区间
				out = append(out, PageRange{begin: p.Empty(), end: p})
			}
		} else {
			front := stack[0]
			top := stack[len(stack)-1]
			if p.format != top.format {
				// 标准 Makeindex 会尝试把区间断开，这里只给出警告
				log.Printf("页码区间可能有误，页码 %s{%s -- %s} 跨过多种数字格式\n", top.encap, top, p)
			}
			if p.encap != front.encap {
				if sorter.strict {
					log.Printf("页码区间可能有误，区间头 %s 没有对应的区间尾\n", front)
					// 输出从区间头到空白的伪区间，并清空栈
					out = append(out, PageRange{begin: front, end: front.Empty()})
					stack = nil
					// 退回重新处理此项
					i--
					continue
				} else {
					// 只输出独立页面，与 Makeindex 行为类似
					if p.rangetype == PAGE_NORMAL {
						out = append(out, PageRange{begin: p, end: p})
					} else {
						log.Printf("页码区间 %s{%s--} 内 %s%s{%s} 命令格式不同，可能丢失信息",
							front.encap, front, p.rangetype, p.encap, p)
					}
				}
			}
			switch p.rangetype {
			case PAGE_NORMAL:
				// 什么也不做
			case PAGE_OPEN:
				// 压栈
				stack = append(stack, p)
			case PAGE_CLOSE:
				// 栈中只有一个元素时输出正常区间，弹栈
				if len(stack) == 1 {
					out = append(out, PageRange{begin: front, end: p})
				}
				stack = stack[:len(stack)-1]
			}
		}
	}
	if len(stack) > 0 {
		log.Printf("页码区间有误，未找到与 %s{%s} 匹配的区间尾。\n", stack[0].encap, stack[0])
		// 输出从当前页到空白的伪区间
		out = append(out, PageRange{begin: stack[0], end: stack[0].Empty()})
	}
	//	debug.Println(out)
	return out
}

// 合并相邻的页码区间
// 输入是 1 2--3 4--6 7，输出 1--7
func (sorter *PageSorter) Merge(pages []PageRange) []PageRange {
	var out []PageRange
	for i, r := range pages {
		// 跳过首项；按设置跳过单页页码
		if i == 0 {
			out = append(out, r)
			continue
		}
		// 合并重复页和区间
		prev := out[len(out)-1]
		if sorter.disable_range &&
			(r.begin.rangetype == PAGE_NORMAL || prev.begin.rangetype == PAGE_NORMAL) {
			// 合并（跳过）重复页
			if prev.begin == r.begin {
				continue
			} else {
				out = append(out, r)
			}
		} else if prev.begin.encap == r.begin.encap &&
			r.begin.format == prev.begin.format &&
			r.begin.page-prev.end.page <= 1 {
			// 合并区间，只用后一区间尾替换前一区间尾
			out[len(out)-1].end = r.end
		} else {
			out = append(out, r)
		}
	}
	// 修正区间类型（似乎无用）
	for i := range out {
		if out[i].begin.encap == out[i].end.encap {
			if out[i].begin.format == out[i].end.format &&
				out[i].begin.page == out[i].end.page {
				out[i].begin.rangetype = PAGE_NORMAL
				out[i].end.rangetype = PAGE_NORMAL
			}
			// 保留首尾区间类型，可以输出时判断是否是合并得到的区间
		}
		// encap 不同是不匹配区间或不完全区间，不修正
	}
	return out
}

type PageInputSlice struct {
	pages  []PageInput
	sorter *PageSorter
}

func (p PageInputSlice) Len() int {
	return len(p.pages)
}

func (p PageInputSlice) Swap(i, j int) {
	p.pages[i], p.pages[j] = p.pages[j], p.pages[i]
}

type PageInputSliceStrict struct {
	PageInputSlice
}

// 先按 encap 类型比较，然后按页码类型，然后页码数值，最后是 rangetype，方便以后合并
// 不同 encap 严格分离
func (p PageInputSliceStrict) Less(i, j int) bool {
	a, b := p.pages[i], p.pages[j]
	if a.encap < b.encap {
		return true
	} else if a.encap > b.encap {
		return false
	}
	if p.sorter.precedence[a.format] < p.sorter.precedence[b.format] {
		return true
	} else if p.sorter.precedence[a.format] > p.sorter.precedence[b.format] {
		return false
	}
	if a.page < b.page {
		return true
	} else if a.page > b.page {
		return false
	}
	if a.rangetype < b.rangetype {
		return true
	} else {
		return false
	}
}

type PageInputSliceLoose struct {
	PageInputSlice
}

// 先按页码类型比较，然后按页码数值，然后 rangetype，最后是 encap 类型，方便以后合并
// 允许不同 encap 合并，接近传统的 Makeindex 行为
func (p PageInputSliceLoose) Less(i, j int) bool {
	a, b := p.pages[i], p.pages[j]
	if p.sorter.precedence[a.format] < p.sorter.precedence[b.format] {
		return true
	} else if p.sorter.precedence[a.format] > p.sorter.precedence[b.format] {
		return false
	}
	if a.page < b.page {
		return true
	} else if a.page > b.page {
		return false
	}
	if a.rangetype < b.rangetype {
		return true
	} else if a.rangetype > b.rangetype {
		return false
	}
	if a.encap < b.encap {
		return true
	} else {
		return false
	}
}

// 忽略大小写，按内码比较两个字符
// 此过程被其他 collator 的 RuneCmp 调用
func RuneCmpIgnoreCases(a, b rune) int {
	la, lb := unicode.ToLower(a), unicode.ToLower(b)
	return int(la - lb)
}

// 测试是否是数字，但把“〇”单独算做汉字
func IsNumRune(r rune) bool {
	return unicode.IsNumber(r) && r != '〇'
}

// 测试是否为数字串
// 此过程被其他 collator 的 RuneCmp 调用
func IsNumString(s string) bool {
	for _, r := range s {
		if !IsNumRune(r) {
			return false
		}
	}
	return true
}

// 按数字大小比较自然数串，如果不是自然数串视为相等
func DecimalStrcmp(a, b string) int {
	aint, err := strconv.ParseUint(a, 10, 64)
	if err != nil {
		return 0
	}
	bint, err := strconv.ParseUint(b, 10, 64)
	if err != nil {
		return 0
	}
	switch {
	case aint < bint:
		return -1
	case aint > bint:
		return 1
	default:
		return 0
	}
}
