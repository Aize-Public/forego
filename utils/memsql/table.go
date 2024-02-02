package memsql

import "github.com/Aize-Public/forego/ctx"

type table struct {
	cols []col
	rows [][]any
}

/*
func (this *table) insert(c ctx.C, row row) error {
	err := this.validRow(c, row)
	if err != nil {
		return err
	}
	// TODO unique indexes
	this.rows = append(this.rows, row)
	return nil
}
*/

func (this *table) validRow(c ctx.C, row row) error {
	return nil
}
