package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"pvZ/init/migrations"
	"pvZ/internal/adapters/api/rest"
	"pvZ/internal/adapters/api/rest/middleware"
	"pvZ/internal/adapters/db/postgreSQL"
	"pvZ/internal/config"
	"pvZ/internal/domain/usecases/usecase_impl"
)

func RunPVZ() {

	cfg := config.Load()

	fmt.Println("DSN:", cfg.DB.DSN())

	dbx := sqlx.MustConnect("postgres", cfg.DB.DSN())
	defer dbx.Close()

	db := dbx.DB
	if err := migrations.RunMigrations(db); err != nil {
		log.Println("Пиздец2")
		log.Println(err)
		log.Println("Пиздец3")

	}

	secretKey := []byte(cfg.App.JWTSecret)

	userRepo := postgreSQL.NewUserRepository(dbx)
	userUsecase := usecase_impl.NewUserUsecase(userRepo, secretKey)
	userController := rest.NewUserController(userUsecase)
	auth := middleware.NewAuthMiddleware(secretKey)

	pvzRepo := postgreSQL.NewPVZRepository(dbx)
	receptionRepo := postgreSQL.NewReceptionRepository(dbx)
	productRepo := postgreSQL.NewProductRepository(dbx)

	pvzUC := usecase_impl.NewPVZUsecase(pvzRepo)
	receptionUC := usecase_impl.NewReceptionUsecase(receptionRepo)
	productUC := usecase_impl.NewProductUsecase(productRepo, receptionRepo)

	pvzController := rest.NewPVZController(pvzUC)
	receptionController := rest.NewReceptionController(receptionUC)
	productController := rest.NewProductController(productUC)

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

	fmt.Printf("Сервер на порту :%s\n", cfg.App.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.App.Port, r))

}
