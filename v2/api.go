package dnssd

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type LookupRequest struct {
	Name string
	Timeout time.Duration
	MaxResults int
}

func BrowseAsync(name string) (c <-chan *BrowseEntry, cancel func()) {
	ch := make(chan *BrowseEntry, 1)
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	addFn := func(e BrowseEntry) {
		ch <- &e
	}
	rmvFn := func(e BrowseEntry) {
	}
	service := fmt.Sprintf("_%s._tcp.local.", name)
	go func() {
		err := LookupType(ctx, service, addFn, rmvFn)
		if errors.Is(err, context.Canceled) {
			close(ch)
		}
	}()
	c = ch
	return c, cancel
}

func Browse(req LookupRequest) []*BrowseEntry {
	var timer *time.Timer
	if req.Timeout != 0 {
		timer = time.NewTimer(req.Timeout)
	} else {
		timer = time.NewTimer(time.Second)
	}
	n := 0
	ch, cancel := BrowseAsync(req.Name)
	res := []*BrowseEntry{}
	for {
		select {
		case svc, ok := <-ch:
			if !ok {
				timer.Stop()
				return res
			}
			res = append(res, svc)
			n += 1
			if req.MaxResults > 0 && n >= req.MaxResults {
				timer.Stop()
				cancel()
				return res
			}
		case <-timer.C:
			cancel()
			return res
		}
	}
	return res
}

func Register(name string, port int) (shutdown func(), err error) {
	cfg := Config{
		Name: fmt.Sprintf("%s server", name),
		Type: fmt.Sprintf("_%s._tcp", name),
		Port: port,
	}
	svc, err := NewService(cfg)
	if err != nil {
		return nil, err
	}
	rp, err := NewResponder()
	if err != nil {
		return nil, err
	}
	handle, err := rp.Add(svc)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	go rp.Respond(ctx)
	return func() {
		rp.Remove(handle)
		cancel()
	}, nil
}
