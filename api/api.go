package api

import "github.com/Aize-Public/forego/ctx"

type Op interface {
	Do(c ctx.C) error
}
