package ctx

// TEMP
func Span(c C, name string) (C, CancelFunc) {
	c, cf := WithCancel(c)
	// TODO add opentelemetry support
	return c, cf
}
