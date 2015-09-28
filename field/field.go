package field

import (
	"sort"
	"strconv"
	"strings"
)

type Field struct {
	Name    string
	Amount  float64
	IsShare bool
}

func (f *Field) ToString(length int) (s string) {
	name := f.formatName()
	if nameLength := len(name); nameLength < length {
		s = name + ":" + strings.Repeat(" ", length-nameLength)
	} else {
		s = name + ":"
	}
	return s + " $" + strconv.FormatFloat(f.Amount, 'f', 2, 64) + "\n"
}

func (f *Field) formatName() (s string) {
	s = strings.Title(f.Name)
	if f.IsShare {
		arr := strings.Split(s, " ")
		first, rest := arr[0]+"'s", arr[1:]
		s = strings.Join(append([]string{first}, rest...), " ")
	}
	return strings.Replace(s, "-", " ", -1)
}

type Fields []Field

func (fields Fields) LongestTitle() (l int) {
	for _, b := range fields {
		if length := len(b.formatName()); length > l {
			l = length
		}
	}
	return
}

func (fields Fields) ToString(l int) (s string, longest int) {
	sort.Sort(fields)
	for _, f := range fields {
		str := f.ToString(l)
		s += str
		if length := len(str); length > longest {
			longest = length
		}
	}
	return s, longest
}

func (fields Fields) Total() (total float64) {
	for _, f := range fields {
		total += f.Amount
	}
	return
}

func (fields Fields) Len() int {
	return len(fields)
}

func (fields Fields) Less(i, j int) bool {
	return fields[i].Name < fields[j].Name
}

func (fields Fields) Swap(i, j int) {
	fields[i], fields[j] = fields[j], fields[i]
}
