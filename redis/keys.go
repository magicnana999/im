package redis

import "fmt"

const (
	userSig = "im:%s:user:sig:%s"
)

func KeyUserSig(appId, sig string) string {
	return fmt.Sprintf(userSig, appId, sig)
}
