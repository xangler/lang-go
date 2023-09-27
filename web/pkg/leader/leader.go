package leader

import (
	"context"
	"errors"
	"os"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	notLeader           = 0
	leader              = 1
	defaultConnTimeout  = 10 * time.Second
	minElectStableTime  = 60 * time.Second
	defaultElectTimeout = 300 * time.Second
)

type PodLeader struct {
	id             string
	leader         int32
	conf           *PLFConfig
	startedLeading chan struct{}
	stoppedLeading chan struct{}
	forwarder      *Forwarder
	k8sClient      *ClientSet
}

func NewPodLeader(conf *PLFConfig) (*PodLeader, error) {
	pod := PodLeader{
		conf:           conf,
		id:             getHostID(),
		startedLeading: make(chan struct{}),
		stoppedLeading: make(chan struct{}),
		forwarder:      NewForwarder(conf.ForWardConfig, grpc.WithTransportCredentials(insecure.NewCredentials())),
	}
	err := pod.startElectionIfNeeded(context.Background())
	return &pod, err
}

func (s *PodLeader) startElectionIfNeeded(ctx context.Context) error {
	if s.conf.LeaderElect {
		var err error
		s.k8sClient, err = NewClientSet(ctx, s.conf.ClientConfig)
		if err != nil {
			return err
		}
		go func() {
			backoff := []int{1, 2, 4, 8, 16, 32, 64}
			err := leaderElectWithRetries(ctx, minElectStableTime, defaultElectTimeout, backoff, func() {
				s.k8sClient.LeaderElect(ctx, s)
			})
			if err != nil {
				log.Fatal(err)
			}
		}()
		return nil
	}
	// disable leader election in single node mode
	s.OnStartedLeading(ctx)
	return nil
}

func (s *PodLeader) GetForwarder() *Forwarder {
	if s.conf.ServerForward {
		return s.forwarder
	}
	return nil
}

func (s *PodLeader) StartedLeading() <-chan struct{} {
	return s.startedLeading
}

func (s *PodLeader) StoppedLeading() <-chan struct{} {
	return s.stoppedLeading
}

// IsLeader ...
func (s *PodLeader) IsLeader() bool {
	val := atomic.LoadInt32(&s.leader)
	return val == leader
}

// BecomeLeader ...
func (s *PodLeader) BecomeLeader() bool {
	flg := atomic.CompareAndSwapInt32(&s.leader, notLeader, leader)
	if flg && s.conf.WithNotifyChan {
		s.startedLeading <- struct{}{}
	}
	return flg
}

// Resign ...
func (s *PodLeader) Resign() {
	if s.conf.WithNotifyChan {
		s.stoppedLeading <- struct{}{}
	}
	atomic.StoreInt32(&s.leader, notLeader)
}

// OnStartedLeading ...
func (s *PodLeader) OnStartedLeading(ctx context.Context) {
	if s.BecomeLeader() {
		log.Info("I'm in leading now, id: ", s.id)
	}
}

// OnStoppedLeading ...
func (s *PodLeader) OnStoppedLeading() {
	// we can do cleanup here, or after the RunOrDie method
	// returns
	s.Resign()
	log.Infof("%s: lost", s.id)
}

// OnNewLeader ...
func (s *PodLeader) OnNewLeader(identity string) {
	// we're notified when new leader elected
	if identity == s.id {
		// I just got the lock
		log.Info("I just got the lead elect lock, ", s.id)
		return
	}
	log.Info("identify: ", identity, " my id: ", s.id)
	if s.conf.ServerForward {
		_, err := s.forwarder.ConnForAddr(identity)
		if err != nil {
			log.Error("connect to leader failed, ", err)
			return
		}
	}
	log.Info("new leader elected: ", identity)
}

// GetID ...
func (s *PodLeader) GetID() string {
	return s.id
}

func getHostID() string {
	host, ok := os.LookupEnv("POD_IP")
	if !ok {
		host = "localhost"
	}
	return host
}

func leaderElectWithRetries(ctx context.Context, minElectStableTime, electTimeout time.Duration, backoff []int, elect func()) error {
	i := 0
	timer := time.NewTimer(electTimeout)
	defer timer.Stop()
	// if election exception occurred, start the election again
	for {
		select {
		case <-ctx.Done():
			log.Info("context done, election exited")
			return nil
		case <-timer.C:
			return errors.New("the leader election timeout")
		default:
			beforeElect := time.Now()
			elect()
			// step here if election exception occurred or ctx done
			// too short, seen as an unstable election, and will retry again
			if time.Since(beforeElect) < minElectStableTime {
				log.Warn("restart election after ", backoff[i], " seconds")
				time.Sleep(time.Duration(backoff[i]) * time.Second)
				if i < len(backoff)-1 {
					i++
				}
				continue
			}
			// or else reset retry status
			i = 0
			select {
			case <-timer.C:
			default:
			}
			timer.Reset(electTimeout)
		}
	}
}
