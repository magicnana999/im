package redis

import "fmt"

const (
	broker       = "im:broker:%s"
	userSig      = "im:%s:user:sig:%s"
	sequence     = "im:%s:sequence:%s"
	sequenceLock = "im:%s:sequence:%s:lock"
	userConn     = "im:%s:user:connect:%s"
	userClients  = "im:%s:user:clients:%s"
	userLock     = "im:%s:user:lock:%s"
)

func KeyUserSig(appId, sig string) string {
	return fmt.Sprintf(userSig, appId, sig)
}

func KeySequence(appId, cId string) string {
	return fmt.Sprintf(sequence, appId, cId)
}

func KeySequenceLock(appId, cId string) string {
	return fmt.Sprintf(sequenceLock, appId, cId)
}

func KeyBroker(addr string) string {
	return fmt.Sprintf(broker, addr)
}

func KeyUserConn(appId, ucLabel string) string {
	return fmt.Sprintf(userConn, appId, ucLabel)
}

func KeyUserClients(appId, ucLabel string) string {
	return fmt.Sprintf(userClients, appId, ucLabel)
}

func KeyUserLock(appId, ucLabel string) string {
	return fmt.Sprintf(userLock, appId, ucLabel)
}
