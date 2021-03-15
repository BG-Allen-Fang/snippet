package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	dynamicMiddleware := alice.New(app.session.Enable)
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/profile/", dynamicMiddleware.ThenFunc(app.showProfile))
	mux.Get("/sign", dynamicMiddleware.ThenFunc(app.SignForm))
	mux.Post("/sign", dynamicMiddleware.ThenFunc(app.Sign))
	mux.Get("/login", dynamicMiddleware.ThenFunc(app.LoginForm))
	mux.Get("/logout", dynamicMiddleware.ThenFunc(app.Logout))
	mux.Post("/login", dynamicMiddleware.ThenFunc(app.Login))
	mux.Get("/film/create", dynamicMiddleware.ThenFunc(app.createFilmForm))
	mux.Post("/film/create", dynamicMiddleware.ThenFunc(app.createFilm))
	mux.Get("/film/:id", dynamicMiddleware.ThenFunc(app.showFilm))
	mux.Post("/film/:id", dynamicMiddleware.ThenFunc(app.BuyTicket))
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))
	return standardMiddleware.Then(mux)
}
