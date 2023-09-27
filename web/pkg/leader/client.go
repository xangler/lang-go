package leader

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/transport"

	log "github.com/sirupsen/logrus"
)

type LeaderCallbacks interface {
	// OnStartedLeading is called when a LeaderElector client starts leading
	OnStartedLeading(context.Context)
	// OnStoppedLeading is called when a LeaderElector client stops leading
	OnStoppedLeading()
	// OnNewLeader is called when the client observes a leader that is
	// not the previously observed leader. This includes the first observed
	// leader when the client starts.
	OnNewLeader(identity string)
	// GetID returns the identifier of the candidate
	GetID() string
}

type ClientFactory struct {
	config     *ClientConfig
	restConfig *rest.Config
	// TODO: use generated Clientset and Informer to List-Watch
}

type ClientSetOption func(*ClientFactory)

type ClientSet struct {
	config    *ClientConfig
	clientset *kubernetes.Clientset
}

func NewClientFactory(ctx context.Context, cfg *ClientConfig, opts ...ClientSetOption) (*ClientFactory, error) {
	var config *rest.Config
	var err error
	if cfg.KubeConfig == "" {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		log.Info("load kubeconfig in cluster, ", config)
	} else {
		config, err = clientcmd.BuildConfigFromFlags(cfg.MasterUrl, cfg.KubeConfig)
		if err != nil {
			return nil, err
		}
		log.Info("load kubeconfig from configured kubeconfig, ", config)
	}

	config.Wrap(transport.ContextCanceller(ctx, fmt.Errorf("the client is shutting down")))

	log.Debugf("k8s configs: %+#v", *cfg)
	fac := &ClientFactory{
		config:     cfg,
		restConfig: config,
	}
	for _, o := range opts {
		o(fac)
	}
	return fac, nil
}

func (c *ClientFactory) CreateClientSet() (*ClientSet, error) {
	clientset, err := kubernetes.NewForConfig(c.restConfig)
	if err != nil {
		return nil, err
	}
	log.Info("k8s clientset created")
	clientSet := &ClientSet{
		config:    c.config,
		clientset: clientset,
	}
	return clientSet, nil
}

func NewClientSet(ctx context.Context, config *ClientConfig, opts ...ClientSetOption) (*ClientSet, error) {
	factory, err := NewClientFactory(ctx, config, opts...)
	if err != nil {
		log.Error("failed to create k8s client factory, ", err)
		return nil, err
	}

	cliset, err := factory.CreateClientSet()
	if err != nil {
		log.Error("create client set failed, ", err)
		return nil, err
	}
	return cliset, nil
}

func (c *ClientSet) Clientset() *kubernetes.Clientset {
	return c.clientset
}

func (c *ClientSet) LeaderElect(ctx context.Context, cbs LeaderCallbacks) {
	// leaderElect(ctx, c.config.Namespace, c.clientset, cbs)
	// use client-go leaderelect tool
	leaseLockName := c.config.LockName
	leaseLockNamespace := c.config.Namespace

	// we use the Lease lock type since edits to Leases are less common
	// and fewer objects in the cluster watch "all Leases".
	lock := &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      leaseLockName,
			Namespace: leaseLockNamespace,
		},
		Client: c.clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: cbs.GetID(),
		},
	}

	// start the leader election code loop
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   12 * time.Second,
		RenewDeadline:   5 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: cbs.OnStartedLeading,
			OnStoppedLeading: cbs.OnStoppedLeading,
			OnNewLeader:      cbs.OnNewLeader,
		},
	})

	// we no longer hold the lease, so perform any cleanup and then
	// exit
	log.Infof("lease of %s: done", cbs.GetID())
}
