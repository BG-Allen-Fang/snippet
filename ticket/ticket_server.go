package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	greetpb "se09.com/greetpb_ticket"
	"se09.com/ticket/models/postgresql"
)

var idG = 0
var app *application

type application struct {
	snippets *postgresql.TicketModel
}

type Server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *Server) InsertGreet(ctx context.Context, req *greetpb.GreetRequestInsert) (*greetpb.GreetResponseInsert, error) {

	name := req.GetGreeting().GetName()
	u_id := req.GetGreeting().GetUId()
	time := req.GetGreeting().GetTime()
	price := req.GetGreeting().GetPrice()

	id, err := app.snippets.Insert(name, time, int(u_id), int(price))

	if err != nil {
		return nil, err
	} else {
		res := &greetpb.GreetResponseInsert{
			Result: int64(id),
		}
		return res, nil
	}

}

func (s *Server) LatestGreet(req *greetpb.GreetRequestLatest, stream greetpb.GreetService_LatestGreetServer) error {

	id := req.GetResult()

	models, err := app.snippets.Latest(int(id))

	if err != nil {
		return err

	} else {
		for _, value := range models {
			res := &greetpb.GreetResponseLatest{
				Greeting: &greetpb.Greeting{
					Id:    int64(value.ID),
					Name:  value.Name,
					UId:   int64(value.U_id),
					Time:  value.Time,
					Price: int64(value.Price),
				},
			}
			if err := stream.Send(res); err != nil {
				log.Fatalf("error while sending greet many times responses: %v", err.Error())
			}
		}
		return nil
	}
}

func main() {

	dsn := flag.String("dsn", "postgres://postgres:147urqafjmvz@localhost:5432/Tickets", "PostgreSQL data source name")

	flag.Parse()

	errLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()
	app1 := &application{
		snippets: &postgresql.TicketModel{DB: db},
	}
	app = app1
	l, err := net.Listen("tcp", "0.0.0.0:50041")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &Server{})
	log.Println("Server is running on port:50041")
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
