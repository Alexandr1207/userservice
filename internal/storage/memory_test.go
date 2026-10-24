package storage

import (
	"errors"
	"testing"

	"github.com/Alexandr1207/userservice/internal/user"
)

func TestMemoryRepo_SaveAssignsIDs(t *testing.T) {
	repo := NewMemoryRepo()

	first, err := repo.Save(user.User{Name: "Anna", Email: "anna@x.com"})
	if err != nil {
		t.Fatalf("first Save: unexpected error: %v", err)
	}

	second, err := repo.Save(user.User{Name: "Boris", Email: "boris@x.com"})
	if err != nil {
		t.Fatalf("second Save: unexpected error: %v", err)
	}

	if first.ID != 1 {
		t.Errorf("first ID = %d, want 1", first.ID)
	}

	if second.ID != 2 {
		t.Errorf("second ID = %d, want 2", second.ID)
	}
}

func TestMemoryRepo_SaveDuplicateEmail(t *testing.T) {
	repo := NewMemoryRepo()
	_, err := repo.Save(user.User{Name: "Anna", Email: "anna@com"})
	if err != nil {
		t.Fatalf("first Save: unexpected error: %v", err)
	}
	_, err = repo.Save(user.User{Name: "Boris", Email: "anna@com"})
	if !errors.Is(err, user.ErrMailTaken) {
		t.Errorf("second Save: got error %v, want ErrMailTaken", err)
	}

	list, err := repo.List()
	if err != nil {
		t.Fatalf("List: unexpected error: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("List has %d users, want 1", len(list))
	}
}

func TestMemoryRepo_FindByEmail(t *testing.T) {
	repo := NewMemoryRepo()
	_, err := repo.Save(user.User{Name: "Anna", Email: "anna@x.com"})
	if err != nil {
		t.Fatalf("Save: unexpected error: %v", err)
	}

	tests := []struct {
		name  string
		email string
		want  string
		err   error
	}{
		{name: "found", email: "anna@x.com", want: "Anna", err: nil},
		{name: "not found", email: "nobody@x.com", want: "", err: user.ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByEmail(tt.email)
			if !errors.Is(err, tt.err) {
				t.Fatalf("FindByEmail(%q): got error %v, want %v", tt.email, err, tt.err)
			}
			if got.Name != tt.want {
				t.Errorf("FindByEmail(%q): got name %q, want %q", tt.email, got.Name, tt.want)
			}
		})
	}
}

func TestMemoryRepo_ListSorted(t *testing.T) {
	repo := NewMemoryRepo()

	emails := []string{
		"anna@x.com",
		"boris@x.com",
		"clara@x.com",
		"dmitry@x.com",
		"elena@x.com",
	}

	for _, email := range emails {
		if _, err := repo.Save(user.User{Name: email, Email: email}); err != nil {
			t.Fatalf("Save(%q): unexpected error: %v", email, err)
		}
	}

	for i := 0; i < 20; i++ {
		list, err := repo.List()
		if err != nil {
			t.Fatalf("List (iteration %d): unexpected error: %v", i, err)
		}
		if len(list) != len(emails) {
			t.Fatalf("List (iteration %d): got %d users, want %d", i, len(list), len(emails))
		}
		for j := 1; j < len(list); j++ {
			if list[j-1].ID >= list[j].ID {
				t.Fatalf("List (iteration %d): IDs not strictly increasing at index %d: %d >= %d", i, j, list[j-1].ID, list[j].ID)
			}
		}
	}
}
