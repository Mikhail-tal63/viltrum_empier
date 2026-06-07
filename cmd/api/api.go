package api

import (
	"log"
	"net/http"

	"github.com/Mikhail-tal63/viltrum_empier/internal/auth"
	"github.com/Mikhail-tal63/viltrum_empier/internal/board"
	"github.com/Mikhail-tal63/viltrum_empier/internal/task"
	"github.com/Mikhail-tal63/viltrum_empier/internal/websocket"
	"github.com/Mikhail-tal63/viltrum_empier/internal/workspace"
	"github.com/Mikhail-tal63/viltrum_empier/middleware"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

type APIServer struct {
	addr string
	db   *mongo.Database
	hub  *websocket.Hub
}

func NewAPIServer(addr string, db *mongo.Database, hub *websocket.Hub) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
		hub:  hub,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	userRepo := auth.NewAuthRepository(s.db)
	userService := auth.NewAuthService(userRepo)
	userHandler := auth.NewAuthHandler(userService)
	userHandler.AuthRouter(apiRouter.PathPrefix("/auth").Subrouter())

	/*protected router*************************************************/

	protectedRouter := apiRouter.PathPrefix("").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware)

	userHandler.ProtectedRouter(protectedRouter)

	/*workspace*****************************/

	workspacerepo := workspace.NewWorkspaceRepo(s.db)
	workspaceService := workspace.NewWorkspaceService(workspacerepo)
	workspaceHandler := workspace.NewWorkspaceHandler(workspaceService, userService)
	workspaceHandler.WorkspaceRouter(protectedRouter)

	/*boards & colmuns *******************************/

	boardRepo := board.NewBoardRepository(s.db)
	taskrepo := task.NewTaskRepository(s.db)
	boardService := board.NewBoardService(boardRepo, s.hub, taskrepo)
	boardHandler := board.NewBoardHandler(boardService)
	boardHandler.BoardRouter(protectedRouter)

	/*tasks**********************************************/
	taskrepo = task.NewTaskRepository(s.db)
	taskservice := task.NewTaskService(taskrepo, s.hub, boardRepo)
	taskhandler := task.NewTaskHandler(taskservice)
	taskhandler.TaskRouter(protectedRouter)

	log.Println("listening on ", s.addr)

	/*websokcet*/
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(s.hub, w, r)
	})

	return http.ListenAndServe(s.addr, middleware.CORS(router))

}
