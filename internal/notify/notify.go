package notify

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type Console struct{}

func (Console) Send(to, msg string) error {
	fmt.Println("To:", to, ", msg:", msg)
	return nil
}

type Email struct {
	From string
}

func (e Email) Send(to, msg string) error {
	if !strings.Contains(to, "@") {
		return errors.New("bad address")
	}
	fmt.Println("[email] from=" + e.From + " to:" + to + ", msg: " + msg)
	return nil
}

type Fake struct {
	mu   sync.Mutex
	Sent []string
	Fail bool
}

func (f *Fake) Send(to, msg string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.Fail {
		return errors.New("fake failure")
	}
	f.Sent = append(f.Sent, to+": "+msg)
	return nil
}

type Telegram struct{}

func (t Telegram) Send(to, msg string) error {
	fmt.Println("[telegram] chat=" + to + ": " + msg)
	return nil
}
