package field

import (
	"sort"
	"strconv"
	"strings"
)

type Field struct {
	Name    string
	Amount  int
	IsShare bool
}

func (f *Field) ToString(length int) (s string) {
	name := f.formatName()
	if nameLength := len(name); nameLength < length {
		s = name + ":" + strings.Repeat(" ", length-nameLength)
	} else {
		s = name + ":"
	}
	return s + " $" + f.stringifyAmount() + "\n"
}

func (f *Field) stringifyAmount() string {
	dollars := strconv.Itoa(f.Amount / 100)
	cents := strconv.Itoa(f.Amount % 100)
	if len(cents) == 1 {
		cents = "0" + cents
	}
	return dollars + "." + cents
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

func (fields Fields) LongestName() (l int) {
	for _, b := range fields {
		if length := len(b.formatName()); length > l {
			l = length
		}
	}
	return
}

func (fields Fields) ToString(l int) (s string) {
	sort.Sort(fields)
	for _, f := range fields {
		s += f.ToString(l)
	}
	return
}

func (fields Fields) Total() (total int) {
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
