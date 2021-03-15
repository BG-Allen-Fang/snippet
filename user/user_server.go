package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"se09.com/greetpb_user"
	"se09.com/user/models/postgres"
)

var idG = 0
var loginG = ""

type application struct {
	snippets *postgres.Usermodel
}

type Server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *Server) GreetLogOut(ctx context.Context, req *greetpb.GreetLogout) (*greetpb.GreetLogOutResponse, error) {

	if idG == 0 {
		return &greetpb.GreetLogOutResponse{
			Result: 0,
		}, nil
	} else {
		idG = 0
		loginG = ""
		return &greetpb.GreetLogOutResponse{
			Result: 1,
		}, nil
	}
}

func (s *Server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {

	login := req.GetGreeting().GetLogin()
	pass := req.GetGreeting().GetPassword()

	id := app.snippets.Insert(login, pass)

	res := &greetpb.GreetResponse{
		Result: int64(id),
	}

	return res, nil
}

func (s *Server) GreetCheck(ctx context.Context, req *greetpb.GreetRequestCheck) (*greetpb.GreetResponseCheck, error) {

	login := req.GetGreeting().GetLogin()
	pass := req.GetGreeting().GetPassword()

	user, err := app.snippets.Check(login, pass)
	if user == nil {
		return nil, err
	} else {
		res := &greetpb.GreetResponseCheck{
			Result: int64(user.ID),
		}
		idG = user.ID
		loginG = user.Login
		return res, err
	}
}

func (s *Server) GetGreet(ctx context.Context, req *greetpb.GreetRequestGet) (*greetpb.GreetResponseGet, error) {

	res := &greetpb.GreetResponseGet{
		Greeting: &greetpb.Greeting{
			Id:    int64(idG),
			Login: loginG,
		},
	}

	return res, nil
}

var app *application

func main() {
	dsn := flag.String("dsn", "postgres://postgres:147urqafjmvz@localhost:5432/User", "PostgreSQL data source name")

	flag.Parse()

	errLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()
	app1 := &application{
		snippets: &postgres.Usermodel{DB: db},
	}
	app = app1
	l, err := net.Listen("tcp", "0.0.0.0:99")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &Server{})
	log.Println("Server is running on port:99")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}

func openDB(dsn string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return db, err
}
