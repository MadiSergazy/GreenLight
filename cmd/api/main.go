package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"mado/internal/data"
	"mado/internal/jsonlog"
	"mado/internal/mailer"
)

const version = "1.0.0" // TODO: generate automatically at build time

type config struct {
	port int
	env  string
	db   struct {
		dsn          string // Data Source Name
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	limiter struct {
		rps     float64 // Requests Per Second
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	// Add a cors struct and trustedOrigins field with the type []string.
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development",
		"Environment (development|staging|production)")
	flag.StringVar(
		&cfg.db.dsn,
		"db-dsn",
		// TODO: create an environment variable, GREENLIGHT_DB_DSN, like
		// postgres://username:password@host/dbname
		"", //os.Getenv("GREENLIGHT_DB_DSN"),     ///it is already in makeFIle
		"PostgreSQL DSN",
	)
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25,
		"PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25,
		"PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m",
		"PostgreSQL max connection idle time")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2,
		"Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "0abf276416b183", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "d8672aa2264bb5", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.alexedwards.net>", "SMTP sender")
	/// Use the flag.Func() function to process the -cors-trusted-origins command line
	// flag. In this we use the strings.Fields() function to split the flag value into a
	// slice based on whitespace characters and assign it to our config struct.
	// Importantly, if the -cors-trusted-origins flag is not present, contains the empty
	// string, or contains only whitespace, then strings.Fields() will return an empty
	// []string slice.
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	flag.Parse()

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	//Publish a new "version" variable in the expvar handler containing our application
	// version number (currently the constant "1.0.0").
	expvar.NewString("version").Set(version)

	// Publish the number of active goroutines.
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	// Publish the database connection pool statistics.
	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))
	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
