package repository

var DefaultGroupMemberRepository = &GroupMemberRepository{}

type GroupMemberRepository struct {
}

func InitGroupMemberRepository() *GroupMemberRepository {
	InitGorm()
	return DefaultGroupMemberRepository
}
