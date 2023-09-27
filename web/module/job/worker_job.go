package job

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/learn-go/web/config"
	"github.com/learn-go/web/module/models"
	"github.com/learn-go/web/module/transaction"
	"github.com/learn-go/web/pkg/dbutils"
	log "github.com/sirupsen/logrus"
)

type WorkerJob struct {
	conf       *config.WebConfig
	dClient    transaction.Transaction
	workerLock sync.Mutex
	isLeader   int32
}

func NewWorkerJob(conf *config.WebConfig, dClient transaction.Transaction) *WorkerJob {
	s := &WorkerJob{
		conf:       conf,
		dClient:    dClient,
		isLeader:   0,
		workerLock: sync.Mutex{},
	}
	return s
}

func (s *WorkerJob) Start(ctx context.Context) {
	s.electionAction(ctx)
	s.workAction(ctx)
	go s.electionJob(ctx)
	go s.workJob(ctx)
}

func (s *WorkerJob) electionJob(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.conf.Locker.Period) * time.Second)
	defer ticker.Stop()
Loop:
	for {
		select {
		case <-ctx.Done():
			log.Info("Worker quit election job")
			break Loop
		case t := <-ticker.C:
			log.Infof("Worker election timestamp:%v\n", t)
			err := s.electionAction(ctx)
			if err != nil {
				log.Errorf("electionAction crash %v", err)
			} else {
				log.Info("electionAction succ")
			}
		}
	}
}

func (s *WorkerJob) electionAction(ctx context.Context) error {
	current := time.Now().Unix()
	lockerObj, err := s.dClient.GetLockerByName(ctx, s.conf.Locker.Name, dbutils.NoLock)
	if err != nil {
		return fmt.Errorf("GetLockerByName first crash %v", err)
	}
	if lockerObj == nil {
		lockerObj = &models.Locker{
			Name:    s.conf.Locker.Name,
			Version: current,
			Master:  s.conf.Locker.WorkMark,
		}
		err = s.dClient.AddLocker(ctx, lockerObj)
		if err != nil {
			return fmt.Errorf("AddLock crash %v", err)
		}
	} else if lockerObj.Master == s.conf.Locker.WorkMark || current-lockerObj.Version > s.conf.Locker.Dead {
		oldVersion := lockerObj.Version
		lockerObj.Version = current
		lockerObj.Master = s.conf.Locker.WorkMark
		err = s.dClient.UpdateLockerWithVersion(ctx, lockerObj, oldVersion)
		if err != nil {
			return fmt.Errorf("UpdateLockWithVersion leader crash %v", err)
		}
	} else {
		goto RESET
	}

	lockerObj, err = s.dClient.GetLockerByName(ctx, s.conf.Locker.Name, dbutils.NoLock)
	if err != nil {
		return fmt.Errorf("GetLockerByName second crash %v", err)
	}
RESET:
	if lockerObj.Master == s.conf.Locker.WorkMark {
		if atomic.LoadInt32(&s.isLeader) == 0 {
			log.Info("woker elect as leader")
		}
		atomic.StoreInt32(&s.isLeader, 1)
	} else {
		if atomic.LoadInt32(&s.isLeader) == 1 {
			log.Info("woker elect as follower")
		}
		atomic.StoreInt32(&s.isLeader, 0)
	}
	return nil
}

func (s *WorkerJob) workJob(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(s.conf.CheckPeriod) * time.Second)
	defer ticker.Stop()
Loop:
	for {
		select {
		case <-ctx.Done():
			log.Info("Worker quit work job")
			break Loop
		case t := <-ticker.C:
			log.Infof("Worker work timestamp:%v\n", t)
			err := s.workAction(ctx)
			if err != nil {
				log.Errorf("workAction crash %v", err)
			} else {
				log.Info("workAction succ")
			}
		}
	}
}

func (s *WorkerJob) workAction(ctx context.Context) error {
	if atomic.LoadInt32(&s.isLeader) == 0 {
		log.Info("worker skip this work action")
		return nil
	}
	log.Info("worker action succ")
	return nil
}
