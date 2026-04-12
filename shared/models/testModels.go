package models

type LoggerInterfaceMock struct{}

func (l *LoggerInterfaceMock) Infof(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LoggerInterfaceMock) Errorf(format string, args ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l *LoggerInterfaceMock) Debugf(format string, args ...interface{}) {}
func (l *LoggerInterfaceMock) Warnf(format string, args ...interface{})  {}
