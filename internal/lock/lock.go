package lock

import (
	"context"
	"github.com/coreos/airlock/internal/config"
)

type Manager interface {
	RecursiveLock(ctx context.Context, key string) (*Semaphore, error)
	UnlockIfHeld(ctx context.Context, id string) (*Semaphore, error)
	FetchSemaphore(ctx context.Context) (*Semaphore, error)
	Close()
}

func NewManager(ctx context.Context, a *config.Settings, group string, maxSlots uint64) (Manager, error) {

	// Create a new lock manager.
	var err error
	var manager Manager
	if len(a.EtcdEndpoints) > 0 {
		manager, err = NewEtcdManager(ctx, a.EtcdEndpoints, a.ClientCertPubPath, a.ClientCertKeyPath, a.EtcdTxnTimeout, group, maxSlots)
	}
	if a.MemDBEnabled {
		manager, err = NewMemDBManager(ctx, group, maxSlots)
	}
	return manager, err
}
