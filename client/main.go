package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	pb "moviesapp.com/grpc/protos"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(),
		grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewMoviesappClient(conn)

	runGetMovies(client)
	// runGetMovie(client, "1")
	// runCreateMovie(client, "24325645", "Spiderman Spiderverse",
	// 	"Stan", "Lee")
	// runUpdateMovie(client, "98498081", "24325645", "Spiderman Spiderverse",
	// 	"Peter", "Parker")
	// runDeleteMovie(client, "98498081")
}

func runGetMovies(client pb.MoviesappClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Empty{}
	stream, err := client.GetMovies(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetMovies(_) = _, %v", client, err)
	}
	for {
		row, err := stream.Recv()
		if err != io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetMovies(_) = _, %v", client, err)
		}
		log.Fatalf("MovieInfo: %v", row)

	}
}

func runGetMovie(client pb.MoviesappClient, movieid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: movieid}
	res, err := client.GetMovie(ctx, req)
	if err != nil {
		log.Fatalf("%v.GetMovie(_) = _, %v", client, err)
	}
	log.Printf("MovieInfo: %v", res)
}

func runCreateMovie(client pb.MoviesappClient, isbn string,
	title string, firstname string, lastname string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.MovieInfo{Isbn: isbn, Title: title,
		Director: &pb.Director{Firstname: firstname,
			Lastname: lastname}}
	res, err := client.CreateMovie(ctx, req)
	if err != nil {
		log.Fatalf("%v.CreateMovie(_) = _, %v", client, err)
	}
	if res.GetValue() != "" {
		log.Printf("CreateMovie Id: %v", res)
	} else {
		log.Printf("CreateMovie Failed")
	}

}

func runUpdateMovie(client pb.MoviesappClient, movieid string,
	isbn string, title string, firstname string, lastname string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.MovieInfo{Id: movieid, Isbn: isbn,
		Title: title, Director: &pb.Director{
			Firstname: firstname, Lastname: lastname}}
	res, err := client.UpdateMovie(ctx, req)
	if err != nil {
		log.Fatalf("%v.UpdateMovie(_) = _, %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("UpdateMovie Success")
	} else {
		log.Printf("UpdateMovie Failed")
	}
}

func runDeleteMovie(client pb.MoviesappClient, movieid string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	req := &pb.Id{Value: movieid}
	res, err := client.DeleteMovie(ctx, req)
	if err != nil {
		log.Fatalf("%v.DeleteMovie(_) = _, %v", client, err)
	}
	if int(res.GetValue()) == 1 {
		log.Printf("DeleteMovie Success")
	} else {
		log.Printf("DeleteMovie Failed")
	}
}
