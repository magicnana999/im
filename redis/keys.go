package redis

import "fmt"

const (
	broker           = "im:broker:%s"
	userSig          = "im:%s:user:sig:%s"
	user             = "im:%s:user:%d"
	userLock         = "im:%s:user:%d:lock"
	sequence         = "im:%s:sequence:%s"
	sequenceLock     = "im:%s:sequence:%s:lock"
	userConn         = "im:%s:user:connect:%s"
	userClients      = "im:%s:user:clients:%d"
	userConnLock     = "im:%s:user:connect:%s:lock"
	groupMembers     = "im:%s:group:members:%d"
	groupMembersLock = "im:%s:group:members:%d:lock"
)

func KeyUserSig(appId, sig string) string {
	return fmt.Sprintf(userSig, appId, sig)
}

func KeySequence(appId, sequenceId string) string {
	return fmt.Sprintf(sequence, appId, sequenceId)
}

func KeySequenceLock(appId, seqId string) string {
	return fmt.Sprintf(sequenceLock, appId, seqId)
}

func KeyBroker(addr string) string {
	return fmt.Sprintf(broker, addr)
}

func KeyUserConn(appId, ucLabel string) string {
	return fmt.Sprintf(userConn, appId, ucLabel)
}

func KeyUserClients(appId string, userId int64) string {
	return fmt.Sprintf(userClients, appId, userId)
}

func KeyUserConnLock(appId, ucLabel string) string {
	return fmt.Sprintf(userConnLock, appId, ucLabel)
}

func KeyGroupMembers(appId string, groupId int64) string {
	return fmt.Sprintf(groupMembers, appId, groupId)
}

func KeyGroupMembersLock(appId string, groupId int64) string {
	return fmt.Sprintf(groupMembersLock, appId, groupId)
}

func KeyUser(appId string, userId int64) string {
	return fmt.Sprintf(user, appId, userId)
}

func KeyUserLock(appId string, userId int64) string {
	return fmt.Sprintf(userLock, appId, userId)
}
