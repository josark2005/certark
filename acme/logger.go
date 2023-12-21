package acme

import (
	"fmt"

	"github.com/go-acme/lego/v4/log"
	"github.com/josark2005/certark/ark"
)

type ArkLogger struct {
}

func init() {
	log.Logger = &ArkLogger{}
}

func (a *ArkLogger) Fatal(args ...interface{}) {
	msg := fmt.Sprint(args...)
	ark.Fatal().Msg(msg)
}

func (a *ArkLogger) Fatalln(args ...interface{}) {
	msg := fmt.Sprintln(args...)
	ark.Fatal().Msg(msg)
}

func (a *ArkLogger) Fatalf(f string, args ...interface{}) {
	msg := fmt.Sprintf(f, args...)
	ark.Fatal().Msg(msg)
}

func (a *ArkLogger) Print(args ...interface{}) {
	msg := fmt.Sprint(args...)
	ark.Info().Msg(msg)
}

func (a *ArkLogger) Println(args ...interface{}) {
	msg := fmt.Sprintln(args...)
	ark.Info().Msg(msg)
}

func (a *ArkLogger) Printf(f string, args ...interface{}) {
	msg := fmt.Sprintf(f, args...)
	ark.Info().Msg(msg)
}
