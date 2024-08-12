package app

import (
	"backend/internal/lib/e"
	"backend/internal/lib/rr"
	as "backend/internal/modules/auth/service"
	br "backend/internal/modules/books/repository"
	bs "backend/internal/modules/books/service"
	"backend/internal/modules/superservice/controller"
	"backend/internal/modules/superservice/facade"
	ss "backend/internal/modules/superservice/service"
	ur "backend/internal/modules/users/repository"
	us "backend/internal/modules/users/service"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// include to use db drivers
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server interface {
	Start()
	Stop()
}

type App struct {
	host       string
	port       string
	jwtSecret  string
	dsn        string
	server     *http.Server
	signalChan chan os.Signal
	db         *sql.DB
	facade     facade.LibraryFacader
	controller controller.LibraryController
}

func NewApp() (*App, error) {
	a := &App{}

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Start() {
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			log.Println("error closing database connection")
		}
	}(a.db)

	fmt.Println("Started server on port", a.port)
	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func (a *App) Stop() {
	<-a.signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	<-ctx.Done()

	fmt.Println("Shutting down server gracefully")
}

func (a *App) init() (err error) {
	defer func() { err = e.WrapIfErr("failed to init app", err) }()

	if err = a.readConfig(); err != nil {
		return err
	}

	if a.db, err = a.connectToDB(); err != nil {
		return err
	}

	// service initialization
	booksService := bs.NewBooksService(br.NewPostgresDBRepo(a.db))
	usersService := us.NewUserService(ur.NewPostgresDBRepo(a.db))
	authService := as.NewAuthService(a.host, a.host, a.jwtSecret, a.host)
	a.facade = facade.NewFacade(ss.NewSuperService(booksService, usersService, authService))
	a.controller = controller.NewLibraryControl(a.facade, rr.NewReadRespond())

	// prepopulate tables with mockup data
	if err = a.facade.InitTables(); err != nil {
		return e.Wrap("couldn't init tables", err)
	}

	a.server = &http.Server{
		Addr:         ":" + a.port,
		Handler:      a.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	a.signalChan = make(chan os.Signal, 1)
	signal.Notify(a.signalChan, syscall.SIGINT, syscall.SIGTERM)

	return nil
}

func (a *App) readConfig(envPath ...string) (err error) {
	if len(envPath) > 0 {
		err = godotenv.Load(envPath[0])
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		return e.WrapIfErr("can't read .env file", err)
	}

	a.host = os.Getenv("HOST")
	a.port = os.Getenv("PORT")

	a.jwtSecret = os.Getenv("JWT_SECRET")

	a.dsn = fmt.Sprintf( // database source name
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5\n",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	return nil
}

func (a *App) connectToDB() (conn *sql.DB, err error) {
	defer func() { err = e.WrapIfErr("failed to connect to database", err) }()

	conn, err = sql.Open("pgx", a.dsn)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	return conn, nil
}
