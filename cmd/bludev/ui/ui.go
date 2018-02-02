package ui

type UI interface {
	SetStatus(status ...interface{})
	Log(msg ...interface{})
	Logf(format string, args ...interface{})
	HRule()

	AddSlot(label string)
	SetSlot(label, buf string)
}
