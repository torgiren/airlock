package lock

import (
	"context"
	"errors"
)

type MemDBManager struct {
	lock *Semaphore
}

var (
	ErrEmptyGroupName  = errors.New("empty group name")
	ErrZeroSlots       = errors.New("zero slots specified")
	ErrNilMemDBManager = errors.New("nil MemDBManager")

	memoryDB map[string]*Semaphore = make(map[string]*Semaphore, 0)
)

func NewMemDBManager(ctx context.Context, group string, maxSlots uint64) (*MemDBManager, error) {
	if group == "" {
		return nil, ErrEmptyGroupName
	}

	if maxSlots == 0 {
		return nil, ErrZeroSlots
	}

	manager := &MemDBManager{
		lock: memoryDB[group],
	}
	if manager.lock == nil {
		manager.lock = NewSemaphore(maxSlots)
		memoryDB[group] = manager.lock
	}

	return manager, nil
}

func (m *MemDBManager) RecursiveLock(ctx context.Context, key string) (*Semaphore, error) {
	if m.lock == nil {
		return nil, ErrNilMemDBManager
	}

	held, err := m.lock.RecursiveLock(key)
	if err != nil {
		return nil, err
	}
	if held {
		return m.lock, nil
	}

	return m.lock, nil
}

func (m *MemDBManager) UnlockIfHeld(ctx context.Context, id string) (*Semaphore, error) {
	if m.lock == nil {
		return nil, ErrNilMemDBManager
	}

	if err := m.lock.UnlockIfHeld(id); err != nil {
		return nil, err
	}

	return m.lock, nil
}

func (m *MemDBManager) FetchSemaphore(ctx context.Context) (*Semaphore, error) {
	if m.lock != nil {
		return m.lock, nil
	}
	return nil, errors.New("no semaphore found")
}

func (m *MemDBManager) Close() {
}
