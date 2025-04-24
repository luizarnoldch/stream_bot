package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/luizarnoldch/stream_bot/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Espera al menos un argumento (el comando)
	args := os.Args[1:]
	if len(args) < 1 {
		help()
		return
	}

	sub := args[0]

	switch sub {
	case "reset":
		resetDB(cfg)
	case "up", "down":
		runGoose(cfg, sub)
	case "gen":
		runGen()
	default:
		help()
	}
}

func resetDB(cfg *config.CONFIG) {
	dbName := cfg.MICRO.DB.PSQL.DB
	log.Printf("Reiniciando base de datos '%s'...", dbName)

	// Construir el comando
	cmd := exec.Command("psql",
		"-h", cfg.MICRO.DB.PSQL.HOST,
		"-p", cfg.MICRO.DB.PSQL.PORT,
		"-d", "postgres", // Conectamos a la DB default
		"-U", cfg.MICRO.DB.PSQL.USER,
		"-f", "db/sql/reset_master.sql",
	)

	// Configurar variables de entorno para la conexión
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PGPASSWORD=%s", cfg.MICRO.DB.PSQL.PASS),
	)

	// Capturar salida
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Ejecutar comando
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error ejecutando reset: %v", err)
	}

	log.Printf("Base de datos '%s' reiniciada exitosamente", dbName)
}

func runGoose(cfg *config.CONFIG, action string) {
	// Cambiar al directorio de schemas
	dir := "db/sql/schemas"
	if err := os.Chdir(dir); err != nil {
		log.Fatalf("No se pudo cambiar al directorio %s: %v", dir, err)
	}

	port, err := strconv.Atoi(cfg.MICRO.DB.PSQL.PORT)
	if err != nil {
		log.Fatalf("Invalid port value: %v", err)
	}

	// Construir la cadena de conexión al estilo PostgreSQL URL
	dbUrl := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		cfg.MICRO.DB.PSQL.USER,
		cfg.MICRO.DB.PSQL.PASS,
		cfg.MICRO.DB.PSQL.HOST,
		port,
		cfg.MICRO.DB.PSQL.DB,
	)

	log.Printf("     ** GOOSE %s **", strings.ToUpper(action))

	// Crear el comando con el formato específico
	cmdArgs := []string{
		"postgres", // Driver de PostgreSQL
		dbUrl,      // URL de conexión
		action,     // Acción (up/down)
	}

	cmd := exec.Command("goose", cmdArgs...)

	// Configurar salidas y variables de entorno
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	log.Printf("Ejecutando goose %s en '%s'...",
		action,
		cfg.MICRO.DB.PSQL.DB)

	// Ejecutar y manejar errores
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error ejecutando goose %s: %v", action, err)
	}

	log.Printf("goose %s completado exitosamente", action)
}

func runGen() {
	// Change to the db/sql directory
	dir := "db/sql"
	if err := os.Chdir(dir); err != nil {
		log.Fatalf("No se pudo cambiar al directorio %s: %v", dir, err)
	}

	// Create the sqlc generate command
	cmd := exec.Command("sqlc", "generate")

	// Configure outputs and environment
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	log.Printf("Generating SQL code with sqlc...")

	// Execute and handle errors
	if err := cmd.Run(); err != nil {
		log.Fatalf("Error ejecutando sqlc generate: %v", err)
	}

	log.Printf("SQL code generation completed successfully")
}

func help() {
	fmt.Println("Uso: go run main.go <comando>")
	fmt.Println("Comandos disponibles:")
	fmt.Println("  reset    - Elimina y crea la base de datos configurada")
	fmt.Println("  up       - Ejecuta migraciones con 'goose up'")
	fmt.Println("  down     - Revierte migraciones con 'goose down'")
}
