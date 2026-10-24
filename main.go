package main

import (
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/Alexandr1207/userservice/internal/notify"
	"github.com/Alexandr1207/userservice/internal/storage"
	"github.com/Alexandr1207/userservice/internal/user"
)

func main() {
	// Simple use
	repo := storage.NewMemoryRepo()
	console := notify.Console{}
	svc := user.NewService(repo, console)

	anna, err := svc.Register("Anna", "anna@x.com")
	if err != nil {
		log.Fatalf("register Anna: %v", err)
	}
	fmt.Printf("registered: %+v\n", anna)

	found, err := svc.Find("anna@x.com")
	if err != nil {
		log.Fatalf("find Anna: %v", err)
	}
	fmt.Printf("found: %+v\n", found)

	// Goroutines check
	concurrentRepo := storage.NewMemoryRepo()
	fake := &notify.Fake{}
	concurrentSvc := user.NewService(concurrentRepo, fake)

	const n = 100
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			email := fmt.Sprintf("user%d@x.com", i)
			if _, err := concurrentSvc.Register(fmt.Sprintf("User%d", i), email); err != nil {
				log.Printf("register %s: %v", email, err)
			}
		}(i)
	}
	wg.Wait()

	users, err := concurrentSvc.Users()
	if err != nil {
		log.Fatalf("list users: %v", err)
	}

	if len(users) != n {
		log.Fatalf("concurrency check: got %d users, want %d", len(users), n)
	}

	ids := make([]int, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}
	sort.Ints(ids)

	for i, id := range ids {
		want := i + 1
		if id != want {
			log.Fatalf("concurrency check: at position %d got ID %d, want %d (ids=%v)", i, id, want, ids)
		}
	}

	fmt.Printf("concurrency check passed: %d users, IDs 1..%d without gaps or duplicates\n", len(users), n)
}
