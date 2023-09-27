package transaction

import (
	"context"
	"fmt"
	"sync"

	"github.com/learn-go/web/module/models"
)

type daoActions struct {
	Add    []interface{}
	Update []interface{}
	Delete []interface{}
}

type JobTransaction interface {
	Add(obj interface{})
	Update(obj interface{})
	Delete(obj interface{})
	CancelAll()
	ProcessDaoAction(ctx context.Context) error
}

type jobTransaction struct {
	daoLock sync.Mutex
	cancel  bool
	actions *daoActions
	dClient Transaction
}

func NewJobTransaction(dClient Transaction) JobTransaction {
	return &jobTransaction{
		dClient: dClient,
		cancel:  false,
		daoLock: sync.Mutex{},
		actions: &daoActions{
			Add:    []interface{}{},
			Update: []interface{}{},
			Delete: []interface{}{},
		},
	}
}

func (s *jobTransaction) CancelAll() {
	s.daoLock.Lock()
	defer s.daoLock.Unlock()
	s.cancel = true
}

func (s *jobTransaction) Add(obj interface{}) {
	s.daoLock.Lock()
	defer s.daoLock.Unlock()
	s.actions.Add = append(s.actions.Add, obj)
}

func (s *jobTransaction) Update(obj interface{}) {
	s.daoLock.Lock()
	defer s.daoLock.Unlock()
	s.actions.Update = append(s.actions.Update, obj)
}

func (s *jobTransaction) Delete(obj interface{}) {
	s.daoLock.Lock()
	defer s.daoLock.Unlock()
	s.actions.Delete = append(s.actions.Delete, obj)
}

func (s *jobTransaction) ProcessDaoAction(ctx context.Context) error {
	if s.cancel {
		return nil
	}
	err := s.dClient.WithTransaction(
		func(d Transaction) error {
			if s.actions.Add != nil {
				for _, action := range s.actions.Add {
					switch actionObj := action.(type) {
					case *models.Locker:
						if err := d.AddLocker(ctx, actionObj); err != nil {
							return fmt.Errorf("JobTransactionAction AddLocker error %v", err)
						}
					default:
						return fmt.Errorf("job transaction add not support %v", action)
					}
				}
			}
			if s.actions.Update != nil {
				for _, action := range s.actions.Update {
					switch actionObj := action.(type) {
					case *models.Locker:
						if err := d.UpdateLocker(ctx, actionObj); err != nil {
							return fmt.Errorf("JobTransactionAction UpdateLocker error %v", err)
						}
					default:
						return fmt.Errorf("job transaction update not support %v", action)
					}
				}
			}
			if s.actions.Delete != nil {
				for _, action := range s.actions.Delete {
					switch actionObj := action.(type) {
					case *models.Locker:
						if err := d.DeleteLocker(ctx, actionObj.Id); err != nil {
							return fmt.Errorf("JobTransactionAction DeleteLocker error %v", err)
						}
					default:
						return fmt.Errorf("job transaction delete not support %v", action)
					}
				}
			}
			return nil
		},
	)
	if err != nil {
		return err
	}
	return nil
}
