package daos

import (
	"context"
	"fmt"
	"testing"

	"github.com/learn-go/web/pkg/dbutils"
)

var lockerDao Locker

func TestListLockers(t *testing.T) {
	sin := &LockerSearchArgs{}
	objs, err := lockerDao.ListLockers(context.Background(), sin, 0, 10, dbutils.NoLock)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("find lockers count %v\n", len(objs))
	for _, obj := range objs {
		fmt.Printf("%+v\n", obj)
	}
}
