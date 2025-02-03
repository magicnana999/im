package storage

import (
	"context"
	"github.com/magicnana999/im/service/repository"
	"testing"
)

func TestSetUserByUserSig(t *testing.T) {
	{
		u, _ := repository.SelectUserByUserId("19860220", 1200120)
		SetUserByUserSig(context.Background(), "19860220", u)
	}

	{
		u, _ := repository.SelectUserByUserId("19860220", 1200121)
		SetUserByUserSig(context.Background(), "19860220", u)
	}

	{
		u, _ := repository.SelectUserByUserId("19860220", 1200122)
		SetUserByUserSig(context.Background(), "19860220", u)
	}

}
