package concourse_test

import (
	"fmt"
	"github.com/durmaze/gobank"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"testing"
)

var mbClient *gobank.Client

func getMBHost() string {
	if "" == os.Getenv("MB_HOST") {
		return "localhost"
	}
	return os.Getenv("MB_HOST")
}

func getMBPort() int {
	if "" == os.Getenv("MB_PORT") {
		return 2525
	}
	i, err := strconv.Atoi(os.Getenv("MB_PORT"))
	if err != nil {
		panic(errors.Wrap(err, "Could not start MounteBank"))
	}
	return i
}

func mbURI() string {
	return fmt.Sprintf("http://%s:%d", getMBHost(), getMBPort())
}

func TestMain(m *testing.M) {
	mbClient = gobank.NewClient(mbURI())
	runTests := m.Run()
	mbClient.DeleteAllImposters()
	os.Exit(runTests)
}
