package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	pb "champslol_grpc/proto"

	_ "github.com/denisenkom/go-mssqldb" // Importar driver de SQL Server
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

var db *sql.DB // Puntero a la base de datos

type server struct {
	pb.UnimplementedChampServiceServer
}

// GetChampInfo implementa el método para obtener información de un campeón.
func (s *server) GetChampInfo(ctx context.Context, req *pb.ChampRequest) (*pb.ChampResponse, error) {
	var Name, Rol, Type string

	query := "SELECT Name, Rol, Type FROM champs.leyend WHERE Name LIKE @Name"
	row := db.QueryRowContext(ctx, query, sql.Named("Name", "%"+req.Name+"%"))

	err := row.Scan(&Name, &Rol, &Type)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.ChampResponse{
				Name: "Not found",
				Rol:  "Not found",
				Type: "Not found",
			}, nil
		}
		return nil, err
	}

	return &pb.ChampResponse{
		Name: Name,
		Rol:  Rol,
		Type: Type,
	}, nil
}

// GetChampList implementa el método para obtener la lista de campeones.
func (s *server) GetChampList(req *pb.Empty, stream pb.ChampService_GetChampListServer) error {
	query := "SELECT Name, Rol, Type FROM champs.leyend"
	rows, err := db.Query(query)
	if err != nil {
		log.Panic(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var Name, Rol, Type string
		if err := rows.Scan(&Name, &Rol, &Type); err != nil {
			log.Panic(err)
			return err
		}

		if err := stream.Send(&pb.ChampResponse{
			Name: Name,
			Rol:  Rol,
			Type: Type,
		}); err != nil {
			log.Panic(err)
			return err
		}
	}
	return nil
}

// AddChamps implementa el método para agregar campeones de manera continua.
func (s *server) AddChamps(stream pb.ChampService_AddChampsServer) error {
	var count int32
	for {
		req, err := stream.Recv() // Obtiene la información del cliente
		if err == io.EOF {        // Cuando se termina la transmisión
			return stream.SendAndClose(&pb.AddChampResponse{Count: count})
		}
		if err != nil {
			log.Panic(err)
			return err
		}

		query := "INSERT INTO champs.leyend (Name, Rol, Type) VALUES (@Name, @Rol, @Type)"
		_, err = db.Exec(query,
			sql.Named("Name", req.Name),
			sql.Named("Rol", req.Rol),
			sql.Named("Type", req.Type))
		if err != nil {
			log.Panic(err)
			return err
		}

		count++ // Incrementa el contador
		log.Printf("Inserted %s", req.Name)
	}
}

// GetChampsbyType implementa un método bidireccional para filtrar campeones por tipo.
func (s *server) GetChampsbyType(stream pb.ChampService_GetChampsbyTypeServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("End of stream")
			return nil
		}
		if err != nil {
			log.Panic(err)
			return err
		}

		query := "SELECT Name, Rol, Type FROM champs.leyend WHERE lower(Type) = lower(@Type)"
		rows, err := db.Query(query, sql.Named("Type", req.Type))
		if err != nil {
			log.Panic(err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var Name, Rol, Type string

			if err := rows.Scan(&Name, &Rol, &Type); err != nil {
				log.Panic(err)
				return err
			}
			if err := stream.Send(&pb.ChampResponse{
				Name: Name,
				Rol:  Rol,
				Type: Type,
			}); err != nil {
				log.Panic(err)
				return err
			}
		}
	}
}

// main inicializa el servidor y la conexión con la base de datos.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s := os.Getenv("DB_SERVER")
	puerto := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		user, password, s, puerto, database)

	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	log.Println("Connected to database")

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChampServiceServer(grpcServer, &server{})

	//Iniciar health check
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		log.Println("Starting health check server on port 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	log.Println("Starting server on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
