package main

import (
	_ "context"
	"encoding/json"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"pvZ/dataproviders"
	_ "pvZ/dataproviders/models"
	myhttp "pvZ/entrypoints/http"
	"pvZ/usecases"
	//_ "github.com/lib/pq"
)

type Response struct {
	Message string `json:"message"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(Response{Message: "Привет, JSON!"})
	if err != nil {
		return
	}
}

func main() {
	dsn := "postgres://postgres:password@localhost:55555/pvz?sslmode=disable"
	dbx := sqlx.MustConnect("postgres", dsn)
	defer dbx.Close()

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
