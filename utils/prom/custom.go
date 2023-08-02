package prom

import (
	"fmt"
	"io"
	"sort"
)

type Custom struct {
	Desc string
	Type string                // summary
	Func func() map[string]any // key is the full metric name (including _suffix and {labels})
}

func (this *Custom) Print(name string, w io.Writer) error {
	m := this.Func()
	if len(m) == 0 {
		return nil
	}
	_, err := fmt.Fprintf(w, "# HELP %s %s\n", name, this.Desc)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "# TYPE %s %s\n", name, this.Type)
	if err != nil {
		return err
	}
	ks := make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	for _, k := range ks {
		_, err = fmt.Fprintf(w, "%s %v\n", k, m[k])
		if err != nil {
			return err
		}
	}
	return nil
}
