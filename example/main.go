package main

import (
	"github.com/Aize-Public/forego/ctx"
	"github.com/Aize-Public/forego/ctx/log"
)

func main() {
	c := ctx.TODO()
	log.Debugf(c, "init")
}
