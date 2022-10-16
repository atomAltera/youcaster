package logger

type Opts struct {
	Level  string `long:"level" env:"LEVEL" default:"info"`
	Format string `long:"format" env:"FORMAT" default:"logfmt" choice:"logfmt" choice:"json"`
}
