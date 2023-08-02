package storage

import "github.com/Aize-Public/forego/enc"

type Filter struct {
	Field string
	Cmp   Cmp
	Val   enc.Node
}

type Cmp string

const (
	Equal     Cmp = "EQ"
	NotEqual      = "NEQ"
	Greater       = "G"
	Lesser        = "L"
	GreaterEq     = "GEQ"
	LesserEq      = "LEQ"
)
