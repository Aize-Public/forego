package storage

import "github.com/Aize-Public/forego/enc"

type Filter struct {
	Field string // the field we are checking
	Cmp   Cmp
	Val   enc.Node // the literal we are testing against
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
