package toon

import (
	"strconv"
	"strings"
)

type Field struct {
	Key   string
	Value string
}

func EscapeValue(v string) string {
	v = strings.ReplaceAll(v, "\n", " ")
	v = strings.ReplaceAll(v, "\r", " ")
	return strings.TrimSpace(v)
}

func BuildObject(name string, fields []Field) string {
	var b strings.Builder
	b.WriteString(name)
	b.WriteString(":\n")

	for _, f := range fields {
		if f.Key == "" {
			continue
		}
		b.WriteString("  ")
		b.WriteString(f.Key)
		b.WriteString(": ")
		b.WriteString(EscapeValue(f.Value))
		b.WriteString("\n")
	}

	return b.String()
}

func BuildTable(name string, columns []string, rows [][]string) string {
	var b strings.Builder

	b.WriteString(name)
	b.WriteString("[")
	b.WriteString(strconv.Itoa(len(rows)))
	b.WriteString("]{")
	b.WriteString(strings.Join(columns, ","))
	b.WriteString("}:\n")

	for _, row := range rows {
		if len(row) == 0 {
			continue
		}
		b.WriteString("  ")
		for i, cell := range row {
			if i > 0 {
				b.WriteString(",")
			}
			b.WriteString(EscapeValue(cell))
		}
		b.WriteString("\n")
	}

	return b.String()
}

