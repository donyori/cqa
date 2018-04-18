package web

import (
	"context"
	"errors"
	"net/http"
	"sync"
)

type Server struct {
	server           *http.Server
	sLock            sync.RWMutex
	onShutdownBuffer []func()
	osbLock          sync.RWMutex
}

var (
	ErrNilServer     error = errors.New("server is nil")
	ErrAlreadyLaunch error = errors.New("server has already launched")
	ErrNotLaunch     error = errors.New("server is NOT launched")

	server Server
)

func IsLaunched() bool {
	return server.IsLaunched()
}

func Launch() error {
	return server.Launch(GlobalSettings.GetAddr(), nil)
}

func LaunchBackground() <-chan error {
	return server.LaunchBackground(GlobalSettings.GetAddr(), nil)
}

func Shutdown() error {
	return server.Shutdown()
}

func RegisterOnShutdown(f func()) {
	server.RegisterOnShutdown(f)
}

func NewServer() *Server {
	return new(Server)
}

func (srv *Server) IsLaunched() bool {
	if srv == nil {
		return false
	}
	srv.sLock.RLock()
	defer srv.sLock.RUnlock()
	return srv.server != nil
}

func (srv *Server) Launch(addr string, handler http.Handler) error {
	if srv == nil {
		return ErrNilServer
	}
	Init()
	err := func() error {
		srv.sLock.Lock()
		defer srv.sLock.Unlock()
		if srv.server != nil {
			return ErrAlreadyLaunch
		}
		srv.server = &http.Server{Addr: addr, Handler: handler}
		return nil
	}()
	if err != nil {
		return err
	}
	srv.sLock.RLock()
	defer srv.sLock.RUnlock()
	if srv.server == nil {
		return errors.New("assign server fail")
	}
	func() {
		srv.osbLock.RLock()
		defer srv.osbLock.RUnlock()
		for _, f := range srv.onShutdownBuffer {
			srv.server.RegisterOnShutdown(f)
		}
	}()
	return srv.server.ListenAndServe()
}

func (srv *Server) LaunchBackground(addr string, handler http.Handler) (
	resC <-chan error) {
	errC := make(chan error, 1)
	go func() {
		defer close(errC)
		err := srv.Launch(addr, handler)
		errC <- err
	}()
	return errC
}

func (srv *Server) Shutdown() error {
	if srv == nil {
		return ErrNilServer
	}
	err := func() error {
		srv.sLock.RLock()
		defer srv.sLock.RUnlock()
		if srv.server == nil {
			return ErrNotLaunch
		}
		return srv.server.Shutdown(context.Background())
	}()
	if err != nil {
		return err
	}
	srv.sLock.Lock()
	defer srv.sLock.Unlock()
	srv.server = nil
	return nil
}

func (srv *Server) RegisterOnShutdown(f func()) error {
	if srv == nil {
		return ErrNilServer
	}
	srv.sLock.RLock()
	defer srv.sLock.RUnlock()
	if srv.server != nil {
		srv.server.RegisterOnShutdown(f)
	}
	srv.osbLock.Lock()
	defer srv.osbLock.Unlock()
	srv.onShutdownBuffer = append(srv.onShutdownBuffer, f)
	return nil
}
