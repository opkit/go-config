package url

import (
	"github.com/opkit/go-config/source"
	"time"
)

type urlWatcher struct {
	u    *urlSource
	exit chan bool
	ch   chan *source.ChangeSet
}

func newWatcher(u *urlSource) (*urlWatcher, error) {
	uw := &urlWatcher{
		u:    u,
		ch:   make(chan *source.ChangeSet),
		exit: make(chan bool),
	}
	go func() {
		for {
			if cs, err := u.Read(); err == nil {
				uw.ch <- cs
			}
			time.Sleep(3 * time.Second)
		}
	}()

	return uw, nil
}

func (u *urlWatcher) Next() (*source.ChangeSet, error) {

	select {
	case cs := <-u.ch:
		return cs, nil
	case <-u.exit:
		return nil, source.ErrWatcherStopped
	}
}

func (u *urlWatcher) Stop() error {
	select {
	case <-u.exit:
	default:
	}
	return nil
}
