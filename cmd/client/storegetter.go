package client

import (
	"github.com/harrybrwn/apizza/cmd/cli"
	"github.com/harrybrwn/apizza/cmd/internal/obj"
	"github.com/harrybrwn/apizza/dawg"
	"github.com/harrybrwn/apizza/pkg/errs"
)

// StoreFinder is a mixin that allows for efficient caching and retrival of
// store structs.
type StoreFinder interface {
	Store() *dawg.Store
}

// storegetter is meant to be a mixin for any struct that needs to be able to
// get a store.
type storegetter struct {
	getaddr   func() dawg.Address
	getmethod func() string
	dstore    *dawg.Store
}

// NewStoreGetter will create a new storefinder.
func NewStoreGetter(builder cli.Builder) StoreFinder {
	return &storegetter{
		getmethod: func() string {
			return builder.Config().Service
		},
		getaddr: builder.Address,
		dstore:  nil,
	}
}

// NewStoreGetterFunc creates a new store getter from two funcs
func NewStoreGetterFunc(service func() string, addr func() dawg.Address) StoreFinder {
	return &storegetter{
		getmethod: service,
		getaddr:   addr,
		dstore:    nil,
	}
}

func (s *storegetter) Store() *dawg.Store {
	if s.dstore == nil {
		var err error
		var address = s.getaddr()
		if obj.AddrIsEmpty(address) {
			errs.Handle(errs.New("no address given in config file or as flag"), "Error", 1)
		}
		s.dstore, err = dawg.NearestStore(address, s.getmethod())
		if err != nil {
			errs.Handle(err, "Store Find Error", 1) // will exit
		}
	}
	return s.dstore
}
