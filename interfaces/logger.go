package interfaces

type logger interface {
	Printf(format string, v ...interface{})
}
