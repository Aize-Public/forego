package caches

import (
	"github.com/Aize-Public/forego/ctx"
)

type Cache[K comparable, V any] interface {
	Get(ctx.C, K) (V, HitMiss, error)
	Invalidate(ctx.C, K)
}

type HitMiss string

const (
	HIT  HitMiss = "hit"
	MISS         = "miss"
	ERR          = "err"
	NO           = "no" // no cache
)
