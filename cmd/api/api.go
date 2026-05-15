package api

import (
	"log"
	"net/http"

	"github.com/Mikhail-tal63/viltrum_empier/internal/auth"
	"github.com/Mikhail-tal63/viltrum_empier/middleware"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct {
	addr string
	db   *mongo.Database
}

func NewAPIServer(addr string, db *mongo.Database) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	userRepo := auth.NewAuthRepository(s.db)
	userService := auth.NewAuthService(userRepo)
	userHandler := auth.NewAuthHandler(userService)
	userHandler.AuthRouter(apiRouter.PathPrefix("/auth").Subrouter())

	log.Println("listening on ", s.addr)

	return http.ListenAndServe(s.addr, middleware.CORS(router))

}
