package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/beyachad-maan/auth-service/pkg/postgres"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runCmd = cobra.Command{
	Use:  "run",
	Run:  runService,
	Long: "Run the auth service for the 'beyachad-maan' application",
}

var dbConfig postgres.Config

const serverPort = 8443

// Duration to wait for connections to close once
// a signal was recieved to kill the process.
const connectionsWaitDuration time.Duration = 30 * time.Second

func init() {
	// Add flags to the run command.
	runCmd.Flags().String("database-name", "auth-service-db", "Database name")
	runCmd.Flags().Int("database-port", 5432, "Database port")
	runCmd.Flags().String("database-user", "postgres", "Database user")
	runCmd.Flags().String("database-password", "postgres", "Database password")
	runCmd.Flags().String("certificate-path", "/etc/tls/tls.crt", "Certificate file path")
	runCmd.Flags().String("key-path", "/etc/tls/tls.key", "Key file path")
	runCmd.Flags().String("jwt-public-key-path", "/etc/jwt/jwt_public_key.pem", "JWT public key file path")
	runCmd.Flags().String("jwt-private-key-path", "/etc/jwt/jwt_private_key.pem", "JWT private key file path")

	// Bind the flags to the viper configuration.
	viper.BindPFlag("database-name", runCmd.Flags().Lookup("database-name"))
	viper.BindPFlag("database-port", runCmd.Flags().Lookup("database-port"))
	viper.BindPFlag("database-user", runCmd.Flags().Lookup("database-user"))
	viper.BindPFlag("database-password", runCmd.Flags().Lookup("database-password"))
	viper.BindPFlag("certificate-path", runCmd.Flags().Lookup("certificate-path"))
	viper.BindPFlag("key-path", runCmd.Flags().Lookup("key-path"))
	viper.BindPFlag("jwt-public-key-path", runCmd.Flags().Lookup("jwt-public-key-path"))
	viper.BindPFlag("jwt-private-key-path", runCmd.Flags().Lookup("jwt-private-key-path"))

	// Set default values for the database configuration environment variables.
	// This is crucial for the viper to consider the environment variables.
	viper.BindEnv("DATABASE_USER")
	viper.SetDefault("DATABASE_USER", "postgres")
	viper.BindEnv("DATABASE_PASSWORD")
	viper.SetDefault("DATABASE_PASSWORD", "postgres")
	viper.BindEnv("DATABASE_NAME")
	viper.SetDefault("DATABASE_NAME", "auth-service-db")
	viper.BindEnv("DATABASE_PORT")
	viper.SetDefault("DATABASE_PORT", 5432)

	viper.AutomaticEnv()
	viper.Unmarshal(&dbConfig)
}

func runService(cmd *cobra.Command, args []string) {
	db, err := postgres.ConnectDB(dbConfig.DatabaseUser, dbConfig.DatabasePassword, postgresHost, dbConfig.DatabasePort, dbConfig.DatabaseName)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// jwtPrivateKey, err := jwt.LoadPrivateKey(viper.GetString("jwt-private-key-path"))
	// if err != nil {
	// 	log.Fatalf("Failed to load JWT private key: %v", err)
	// }

	// jwtPublicKey, err := jwt.LoadPublicKey(viper.GetString("jwt-public-key-path"))
	// if err != nil {
	// 	log.Fatalf("Failed to load JWT private key: %v", err)
	// }

	// Load the certificate and key files
	cert, err := tls.LoadX509KeyPair(viper.GetString("certificate-path"), viper.GetString("key-path"))
	if err != nil {
		log.Fatalf("Failed to load TLS certificate and key: %v", err)
	}

	// Initialize the router.
	r := mux.NewRouter()
	// r.HandleFunc("/", handlers.RootHandler).Methods(http.MethodGet)
	// apiV1Router := r.PathPrefix("/api/v1").Subrouter()

	// usersDAO := dao.NewUsers(db)
	// usersHandler := handlers.NewUserHandler(usersDAO, jwtPublicKey)
	//loginHandler := handlers.NewLoginHandler(usersDAO, jwtPrivateKey)
	//apiV1Router.HandleFunc("/users", usersHandler.CreateUser).Methods(http.MethodPost)
	//apiV1Router.HandleFunc("/users/{id}", usersHandler.GetUserByID).Methods(http.MethodGet)
	//apiV1Router.HandleFunc("/users/{id}", usersHandler.DeleteUserByID).Methods(http.MethodDelete)
	//apiV1Router.HandleFunc("/users/{id}/score", usersHandler.CreateUserScoreById).Methods(http.MethodPost)
	//apiV1Router.HandleFunc("/login", loginHandler.LoginUser).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", serverPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	// Run our server in a goroutine so that it doesn't block.
	log.Println("Starting server...")
	go func() {
		if err := srv.ListenAndServeTLS(viper.GetString("certificate-path"), viper.GetString("key-path")); err != nil {
			log.Println(err)
		}
	}()

	done := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(done, os.Interrupt)

	// Block until we receive our signal.
	<-done

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), connectionsWaitDuration)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	log.Println("shutting down...")
	os.Exit(0)

}
