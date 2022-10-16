package logger

import "github.com/sirupsen/logrus"

type Logger interface {
	SetField(key string, value any) Logger
	SetFields(fields map[string]any) Logger
	SetError(err error) Logger

	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
	WithError(err error) Logger

	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Printf(format string, args ...any)
	Warnf(format string, args ...any)
	Warningf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	Panicf(format string, args ...any)

	Debug(args ...any)
	Info(args ...any)
	Print(args ...any)
	Warn(args ...any)
	Warning(args ...any)
	Error(args ...any)
	Fatal(args ...any)
	Panic(args ...any)

	Debugln(args ...any)
	Infoln(args ...any)
	Println(args ...any)
	Warnln(args ...any)
	Warningln(args ...any)
	Errorln(args ...any)
	Fatalln(args ...any)
	Panicln(args ...any)
}

type logger struct {
	log logrus.FieldLogger
}

func WrapLogrus(log logrus.FieldLogger) *logger {
	return &logger{
		log,
	}
}

func (l *logger) SetField(key string, value any) Logger {
	l.log = l.log.WithField(key, value)
	return l
}

func (l *logger) SetFields(fields map[string]any) Logger {
	l.log = l.log.WithFields(fields)
	return l
}

func (l *logger) SetError(err error) Logger {
	l.log = l.log.WithError(err)
	return l
}

func (l *logger) WithField(key string, value any) Logger {
	return WrapLogrus(l.log.WithField(key, value))
}

func (l *logger) WithFields(fields map[string]any) Logger {
	return WrapLogrus(l.log.WithFields(fields))
}

func (l *logger) WithError(err error) Logger {
	return WrapLogrus(l.log.WithError(err))
}

func (l *logger) Debugf(format string, args ...any) {
	l.log.Debugf(format, args...)
}

func (l *logger) Infof(format string, args ...any) {
	l.log.Infof(format, args...)
}

func (l *logger) Printf(format string, args ...any) {
	l.log.Printf(format, args...)
}

func (l *logger) Warnf(format string, args ...any) {
	l.log.Warnf(format, args...)
}

func (l *logger) Warningf(format string, args ...any) {
	l.log.Warningf(format, args...)
}

func (l *logger) Errorf(format string, args ...any) {
	l.log.Errorf(format, args...)
}

func (l *logger) Fatalf(format string, args ...any) {
	l.log.Fatalf(format, args...)
}

func (l *logger) Panicf(format string, args ...any) {
	l.log.Panicf(format, args...)
}

func (l *logger) Debug(args ...any) {
	l.log.Debug(args...)
}

func (l *logger) Info(args ...any) {
	l.log.Info(args...)
}

func (l *logger) Print(args ...any) {
	l.log.Print(args...)
}

func (l *logger) Warn(args ...any) {
	l.log.Warn(args...)
}

func (l *logger) Warning(args ...any) {
	l.log.Warning(args...)
}

func (l *logger) Error(args ...any) {
	l.log.Error(args...)
}

func (l *logger) Fatal(args ...any) {
	l.log.Fatal(args...)
}

func (l *logger) Panic(args ...any) {
	l.log.Panic(args...)
}

func (l *logger) Debugln(args ...any) {
	l.log.Debugln(args...)
}

func (l *logger) Infoln(args ...any) {
	l.log.Infoln(args...)
}

func (l *logger) Println(args ...any) {
	l.log.Println(args...)
}

func (l *logger) Warnln(args ...any) {
	l.log.Warnln(args...)
}

func (l *logger) Warningln(args ...any) {
	l.log.Warningln(args...)
}

func (l *logger) Errorln(args ...any) {
	l.log.Errorln(args...)
}

func (l *logger) Fatalln(args ...any) {
	l.log.Fatalln(args...)
}

func (l *logger) Panicln(args ...any) {
	l.log.Panicln(args...)
}
