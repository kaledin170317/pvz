package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"pvZ/dataproviders"
	myhttp "pvZ/entrypoints/http"
	"pvZ/usecases"
)

func main() {
	fmt.Println("hello мир1")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	fmt.Println("DSN:", dsn) // можно удалить в релизе
	dbx := sqlx.MustConnect("postgres", dsn)
	defer dbx.Close()
	fmt.Println("hello мир3")

	db := dbx.DB
	if err := dataproviders.RunMigrations(db); err != nil {
		log.Fatal(err)
	}

	secretKey := []byte("your-secret")

	userRepo := dataproviders.NewUserRepository(dbx)
	userUsecase := usecases.NewUserUsecase(userRepo, secretKey)
	userController := myhttp.NewUserController(userUsecase)
	auth := myhttp.NewAuthMiddleware(secretKey)

	pvzRepo := dataproviders.NewPVZRepository(dbx)
	receptionRepo := dataproviders.NewReceptionRepository(dbx)
	productRepo := dataproviders.NewProductRepository(dbx)

	pvzUC := usecases.NewPVZUsecase(pvzRepo)
	receptionUC := usecases.NewReceptionUsecase(receptionRepo)
	productUC := usecases.NewProductUsecase(productRepo, receptionRepo)

	pvzController := myhttp.NewPVZController(pvzUC)
	receptionController := myhttp.NewReceptionController(receptionUC)
	productController := myhttp.NewProductController(productUC)

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

	fmt.Println("Сервер на порту :1488")
	log.Fatal(http.ListenAndServe(":1488", r))
}
