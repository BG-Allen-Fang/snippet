package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	greetp "se09.com/greetpb_ticket"
	greetpb "se09.com/greetpb_user"
	"se09.com/pkg/forms"
	"se09.com/pkg/models"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
}

func (app *application) Logout(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:99", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	ctx := context.Background()
	req := &greetpb.GreetLogout{}
	stream, err := c.GreetLogOut(ctx, req)
	if stream.Result == 0 {
		app.session.Put(r, "flash", "You have not logged in yet")
		app.home(w, r)
	} else {
		app.session.Put(r, "flash", "Logged out")
		app.home(w, r)
	}
}

func (app *application) showFilm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})
}
func (app *application) createFilmForm(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:99", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	ctx := context.Background()
	req := &greetpb.GreetRequestGet{}
	stream, err := c.GetGreet(ctx, req)
	user := &models.Kino_user{ID: int(stream.Greeting.Id), Login: stream.Greeting.Login, Balans: 0, Pass: ""}
	if user.ID != 1 {
		app.session.Put(r, "flash", "You have not access")
		app.home(w, r)
	} else {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}
}
func (app *application) createFilm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("Name", "Description", "Time")
	form.MaxLength("Name", 100)
	form.PermittedValues("Time", "365", "7", "1")
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}
	id, err := app.snippets.Insert(form.Get("Name"), form.Get("Description"), form.Get("Time"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Film successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/film/%d", id), http.StatusSeeOther)
}
func (app *application) BuyTicket(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	conn, err := grpc.Dial("localhost:99", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	ctx := context.Background()
	req := &greetpb.GreetRequestGet{}
	stream, err := c.GetGreet(ctx, req)
	user := &models.Kino_user{ID: int(stream.Greeting.Id), Login: stream.Greeting.Login, Balans: 0, Pass: ""}
	if user.ID == 0 {
		app.session.Put(r, "flash", "Firstly login plz for buy ticket")
		app.home(w, r)
	} else {
		film, _ := app.snippets.Get(id)
		check := app.snippets.Update(film.ID)
		if check == "You successfully bought ticket" {
			req1 := &greetp.GreetRequestInsert{
				Greeting: &greetp.Greeting{
					Name:  film.Name,
					Time:  film.Time.String(),
					UId:   int64(user.ID),
					Price: 0,
				},
			}

			con, err := grpc.Dial("localhost:50041", grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
			defer con.Close()

			cc := greetp.NewGreetServiceClient(con)
			_, err = cc.InsertGreet(ctx, req1)
		}
		app.session.Put(r, "flash", check)
		app.home(w, r)
	}
}

func (app *application) SignForm(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:99", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	ctx := context.Background()
	req := &greetpb.GreetRequestGet{}
	stream, err := c.GetGreet(ctx, req)
	user := &models.Kino_user{ID: int(stream.Greeting.Id), Login: stream.Greeting.Login, Balans: 0, Pass: ""}
	if user.ID != 0 {
		app.session.Put(r, "flash", "Logout firstly")
		app.home(w, r)
	} else {
		app.render(w, r, "sign.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}
}

func (app *application) Sign(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("Login", "Pass")
	form.MaxLength("Login", 100)
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}
	conn, err := grpc.Dial("localhost:99", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	ctx := context.Background()
	req := &greetpb.GreetRequest{Greeting: &greetpb.Greeting{
		Id:       int64(0),
		Login:    form.Get("Login"),
		Password: form.Get("Pass"),
		Balans:   string("0"),
	}}
	stream, err := c.Greet(ctx, req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC %v", err)
	}
	if stream.Result != 0 {
		app.session.Put(r, "flash", "User successfully created!")
	} else {
		app.session.Put(r, "flash", "User not created")
	}
	app.home(w, r)
}

func (app *application) LoginForm(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:99", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	ctx := context.Background()
	req := &greetpb.GreetRequestGet{}
	stream, err := c.GetGreet(ctx, req)
	user := &models.Kino_user{ID: int(stream.Greeting.Id), Login: stream.Greeting.Login, Balans: 0, Pass: ""}
	if user.ID != 0 {
		app.session.Put(r, "flash", "Logout firstly")
		app.home(w, r)
	} else {
		app.render(w, r, "login.page.tmpl", &templateData{
			Form: forms.New(nil),
		})
	}
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("Login", "Pass")
	form.MaxLength("Login", 100)
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}
	conn, err := grpc.Dial("localhost:99", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	ctx := context.Background()
	req := &greetpb.GreetRequestCheck{Greeting: &greetpb.Greeting{
		Id:       int64(0),
		Login:    form.Get("Login"),
		Password: form.Get("Pass"),
		Balans:   string("0"),
	}}
	stream, err := c.GreetCheck(ctx, req)

	if err == nil {
		app.session.Put(r, "flash", "User successfully login!")
	} else {
		app.session.Put(r, "flash", "No such user")
	}
	fmt.Println(stream.Result)
	app.home(w, r)

}

func (app *application) showProfile(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial("localhost:99", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	ctx := context.Background()
	req := &greetpb.GreetRequestGet{}
	stream, err := c.GetGreet(ctx, req)
	user := &models.Kino_user{ID: int(stream.Greeting.Id), Login: stream.Greeting.Login, Balans: 0, Pass: ""}
	if user.ID == 0 {
		app.session.Put(r, "flash", "Firstly login plz")
		app.home(w, r)
	} else {
		req1 := &greetp.GreetRequestLatest{Result: stream.Greeting.Id}

		con, err := grpc.Dial("localhost:50041", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
		defer con.Close()

		cc := greetp.NewGreetServiceClient(con)
		stream1, err := cc.LatestGreet(ctx, req1)
		tickets := []*models.Ticket{}

	LOOP:
		for {
			res, err := stream1.Recv()
			if err != nil {
				if err == io.EOF {
					// we've reached the end of the stream
					break LOOP
				}
				log.Fatalf("error while reciving from GreetManyTimes RPC %v", err)
			}
			ticket := &models.Ticket{
				ID:    int(res.Greeting.Id),
				U_id:  int(res.Greeting.UId),
				Name:  res.Greeting.Name,
				Time:  res.Greeting.Time[:19],
				Price: int(res.Greeting.Price),
			}
			tickets = append(tickets, ticket)
		}
		app.render(w, r, "profile.page.tmpl", &templateData{
			Kino_user: user,
			Ticket:    tickets,
		})
	}
}
