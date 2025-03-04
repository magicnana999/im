package training

import (
	"github.com/magicnana999/im/conf"
	"testing"
)

func Test_foo(t *testing.T) {

	conf.LoadConfig("/Users/jinsong/source/github/im/conf/im-broker.yaml")

	main()
}
