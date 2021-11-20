package api

import (
	"net/http"
	"ordercontext/internal/domain"
	"ordercontext/pkg/aggregate"
	"ordercontext/pkg/httperr"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) useMiddlewares() {
	s.useLogger()
	s.useRecover()
	s.useRequestID()
	s.useTimeout()

	s.useErrorHandler(httperr.NewHandler(
		httperr.DefaultHandler.WithMap(http.StatusBadRequest,
			aggregate.ErrNotFound,
			domain.ErrNotPaid,
			domain.ErrInvalidValue,
		),
		httperr.DefaultHandler.WithMapFunc(
			func(err error) (int, bool) {
				_, ok := err.(validator.ValidationErrors)
				return http.StatusBadRequest, ok
			}),
	))
}

func (s *Server) useTimeout() {
	s.echo.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: time.Second * time.Duration(s.cfg.Server.Timeout),
	}))
}

func (s *Server) useLogger() {
	s.echo.Use(middleware.Logger())
}

func (s *Server) useRecover() {
	s.echo.Use(middleware.Recover())
}

func (s *Server) useRequestID() {
	s.echo.Use(middleware.RequestID())
}

func (s *Server) useErrorHandler(httpErrHandler *httperr.Handler) {
	s.echo.HTTPErrorHandler = httpErrHandler.Handle()
}
