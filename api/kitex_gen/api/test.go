package api

import (
	"fmt"
	"github.com/magicnana999/im/pkg/jsonext"
)

func foo() {
	fmt.Println(string(jsonext.PbMarshalNoErr(NewHeartbeat(100))))
	fmt.Println(string(jsonext.PbMarshalNoErr(NewHeartbeat(100).Wrap())))

	m := NewMessage(1, 1, 1, 1, "2", "1", nil)
	fmt.Println(string(jsonext.PbMarshalNoErr(m)))
}
