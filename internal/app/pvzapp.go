package app

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"pvZ/init/migrations"
	"pvZ/internal/config"
)

func RunPVZ() {

	cfg := config.Load()

	fmt.Println("DSN:", cfg.DB.DSN())

	dbx := sqlx.MustConnect("postgres", cfg.DB.DSN())
	defer dbx.Close()

	db := dbx.DB

	//if err := migrations.RollbackMigrations(db); err != nil {
	//	log.Println("Пиздец4")
	//	log.Println(err)
	//	log.Println("Пиздец5")
	//}

	if err := migrations.RunMigrations(db); err != nil {
		log.Println("Пиздец2")
		log.Println(err)
		log.Println("Пиздец3")
	}

	secretKey := []byte(cfg.App.JWTSecret)

	r := SetupRouter(dbx, secretKey)

	fmt.Printf("Сервер на порту :%s\n", cfg.App.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.App.Port, r))

}
