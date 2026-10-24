package storage

import (
	"sort"
	"sync"

	"github.com/Alexandr1207/userservice/internal/user"
)

type MemoryRepo struct {
	mu      sync.RWMutex
	nextID  int
	byEmail map[string]user.User
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{byEmail: make(map[string]user.User)}
}

func (r *MemoryRepo) Save(u user.User) (user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.byEmail[u.Email]; ok {
		return user.User{}, user.ErrMailTaken
	}
	r.nextID++
	u.ID = r.nextID
	r.byEmail[u.Email] = u
	return u, nil
}

func (r *MemoryRepo) FindByEmail(email string) (user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, ok := r.byEmail[email]
	if !ok {
		return user.User{}, user.ErrNotFound
	}
	return val, nil
}

func (r *MemoryRepo) List() ([]user.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]user.User, 0, len(r.byEmail))
	for _, val := range r.byEmail {
		list = append(list, val)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list, nil
}

var _ user.Repository = (*MemoryRepo)(nil)
