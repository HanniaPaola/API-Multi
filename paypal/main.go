package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"api/paypal/src/archivo/application"
	"api/paypal/src/archivo/infrastructure/api"
	"api/paypal/src/archivo/infrastructure/persistence/mysql"
)

func main() {
	// 1. Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Obtener credenciales de las variables de entorno
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Valor por defecto
	}

	// 3. Configurar conexión a MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al conectar a MySQL: %v", err)
	}
	defer db.Close()

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error al hacer ping a MySQL: %v", err)
	}
	log.Println("Conexión a MySQL establecida correctamente")

	// 4. Configurar CORS (opcional, para desarrollo)
	frontendURL := os.Getenv("FRONTEND_URL")

	// 5. Inicializar repositorio y servicio
	planRepo := mysql.NewMySQLPlanRepository(db)
	planService := application.NewPlanService(planRepo)

	// 6. Configurar router
	r := mux.NewRouter()
	
	// Configurar CORS (solo para desarrollo)
	if frontendURL != "" {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", frontendURL)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
				if r.Method == "OPTIONS" {
					return
				}
				next.ServeHTTP(w, r)
			})
		})
	}

	api.SetupRoutes(r, planService)

	// 7. Iniciar servidor
	serverAddr := fmt.Sprintf(":%s", port)
	log.Printf("Servidor iniciado en %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}