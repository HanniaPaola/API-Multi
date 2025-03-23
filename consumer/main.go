package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"api/consumer/src/application"
	"api/consumer/src/infraestructure"
	"api/consumer/src/infraestructure/controller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func initDB() *sql.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando el archivo .env")
	}
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbSchema := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbUser, dbPass, dbHost, dbSchema)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error conectando a la BD: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a MySQL: %v", err)
	}
	fmt.Println("Conectado a MySQL")
	return db
}

func main() {
	db := initDB()
	defer db.Close()
    RABBIT := "amqp://Eduard034:Hacker20@34.196.146.229/"
    QUEUENAME := "temperature_messages"

	// Configuración de dependencias
	notificationService := infraestructure.NewNotificationServiceImpl()
	tokenRepo := infraestructure.NewTokenRepositoryImpl(db)
	notificationAppService := application.NewNotificationAppService(notificationService, tokenRepo)

	// Configuración del consumidor de RabbitMQ
	rabbitRepo, err := infraestructure.NewRabbitMQRepositoryImpl(RABBIT)
	if err != nil {
		log.Fatalf("Error inicializando RabbitMQ: %v", err)
	}
	rabbitMQAppService := application.NewRabbitMQAppService(rabbitRepo, notificationService)

	// Iniciar el consumidor de RabbitMQ
	go func() {
		if err := rabbitMQAppService.StartConsuming(QUEUENAME); err != nil {
			log.Fatalf("Error iniciando el consumidor de RabbitMQ: %v", err)
		}
	}()

	// Configuración del servidor HTTP
	subscriptionHandler := controller.NewSubscriptionHandler(notificationAppService)
	mux := http.NewServeMux()
	mux.HandleFunc("/suscribe-topic", subscriptionHandler.SubscribeHandler)

	// Usar middleware CORS
	handler := corsMiddleware(mux)

	log.Println("Servidor escuchando en el puerto 8081")
	log.Fatal(http.ListenAndServe(":8081", handler))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Manejar preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}