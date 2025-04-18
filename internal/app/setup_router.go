package app

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"pvZ/internal/adapters/api/rest"
	"pvZ/internal/adapters/api/rest/middleware"
	"pvZ/internal/adapters/db/postgreSQL"
	"pvZ/internal/domain/usecases/usecase_impl"
)

func SetupRouter(dbx *sqlx.DB, secretKey []byte) *mux.Router {
	// Репозитории
	userRepo := postgreSQL.NewUserRepository(dbx)
	pvzRepo := postgreSQL.NewPVZRepository(dbx)
	receptionRepo := postgreSQL.NewReceptionRepository(dbx)
	productRepo := postgreSQL.NewProductRepository(dbx)

	// Usecase
	userUsecase := usecase_impl.NewUserUsecase(userRepo, secretKey)
	pvzUC := usecase_impl.NewPVZUsecase(pvzRepo)
	receptionUC := usecase_impl.NewReceptionUsecase(receptionRepo)
	productUC := usecase_impl.NewProductUsecase(productRepo, receptionRepo)

	// Контроллеры
	userController := rest.NewUserController(userUsecase)
	pvzController := rest.NewPVZController(pvzUC)
	receptionController := rest.NewReceptionController(receptionUC)
	productController := rest.NewProductController(productUC)

	auth := middleware.NewAuthMiddleware(secretKey)

	// Маршруты
	r := mux.NewRouter()

	r.HandleFunc("/dummyLogin", userController.DummyLoginHandler).Methods("POST")
	r.HandleFunc("/register", userController.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", userController.LoginHandler).Methods("POST")

	r.Handle("/pvz", auth.RequireRole("moderator")(http.HandlerFunc(pvzController.CreatePVZHandler))).Methods("POST")
	r.Handle("/pvz", auth.RequireRole("employee")(http.HandlerFunc(pvzController.ListPVZHandler))).Methods("GET")

	r.Handle("/receptions", auth.RequireRole("employee")(http.HandlerFunc(receptionController.CreateReceptionHandler))).Methods("POST")
	r.Handle("/pvz/{pvzId}/close_last_reception", auth.RequireRole("employee")(http.HandlerFunc(receptionController.CloseLastReceptionHandler))).Methods("POST")

	r.Handle("/products", auth.RequireRole("employee")(http.HandlerFunc(productController.AddProductHandler))).Methods("POST")
	r.Handle("/pvz/{pvzId}/delete_last_product", auth.RequireRole("employee")(http.HandlerFunc(productController.DeleteLastProductHandler))).Methods("POST")

	return r
}
