package aws

//go:generate mockgen -source=log.go -destination=mocks/log.go -package=mocks

type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}
