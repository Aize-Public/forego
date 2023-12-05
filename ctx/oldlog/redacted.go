package oldlog

// simple string that will log itself as "***"
type RedactedString string

var _ Loggable = RedactedString("")

func (this RedactedString) LogAs(*Tags) any {
	return "***"
}
