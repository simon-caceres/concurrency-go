package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

func main() {
	// connect to the database
	db := initiateDB()
	db.Ping()
	// create sessions
	session := initiateSession()

	// create loggers
	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create waitgroup
	wg := sync.WaitGroup{}
	// set up the aplication config
	app := Config{
		Session:  session,
		DB:       db,
		Wait:     &wg,
		InfoLog:  infolog,
		ErrorLog: errorlog,
	}
	// set up mailer

	// listen for signals
	go app.listenForShutDown()

	// listen for web connections

	app.serve()

}

func (app *Config) serve() {
	// start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting server on port", webPort)

	err := srv.ListenAndServe()

	if err != nil {

		app.ErrorLog.Fatal(err)
		log.Panic(err)
	}

}

func initiateDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to the database")
	}

	return conn
}

func connectToDB() *sql.DB {
	counts := 0

	dsn := os.Getenv("DB_DSN")

	for {
		connction, err := openDB(dsn)

		if err != nil {
			log.Println("Postgrs not ready yet")
		} else {
			log.Println("Postgres is ready")
			return connction
		}

		if counts > 10 {
			return nil
		}

		log.Println("Waiting for postgres to start")
		time.Sleep(1 * time.Second)
		counts++
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}

func initiateSession() *scs.SessionManager {
	// set up session
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true

	return session
}

func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle:   10,
		MaxActive: 100,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS_DSN"))
		},
	}

	return redisPool
}

func (app *Config) listenForShutDown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	app.shutDown()

	os.Exit(0)
}

func (app *Config) shutDown() {
	// perform any cleanup tasks

	app.InfoLog.Println("Gracefully shutting down the server")

	// block until waitgroup is empy
	app.Wait.Wait()

	app.InfoLog.Println("Server shut down")
}
