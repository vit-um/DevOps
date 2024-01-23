package main

import (
	"database/sql"
	"os"

	"github.com/armon/go-metrics"
	"github.com/go-redis/redis"
	"github.com/nats-io/nats.go"
	"github.com/den-vasyliev/image2ascii/convert"
)

var (
	// API is api ref
	API = make(map[string]string)
	//Urls nats addresses
	Urls *string
	// AppRole app role
	AppRole *string
	//AppPort app port
	AppPort *string
	// AppLicense app
	AppLicense = os.Getenv("APP_LICENSE")
	// AppASCII app
	AppASCII = getEnv("APP_ASCII", "ascii")
	// AppDatastore app
	AppDatastore = getEnv("APP_DATASTORE", "data")
	// AppDbUser name
	AppDbUser = getEnv("AppDbPort", "root")
	// AppDbPort name
	AppDbPort = getEnv("APP_DB_PORT", "3306")
	// AppDbName name
	AppDbName = getEnv("APP_DB_NAME", "demo")
	// AppDb name
	AppDb = getEnv("APP_DB", AppDbUser+"@tcp(127.0.0.1:"+AppDbPort+")/"+AppDbName)
	// AppCache app
	AppCache = getEnv("APP_CACHE", "127.0.0.1")
	// AppCachePort app
	AppCachePort = getEnv("APP_CACHE_PORT", "6379")
	// AppCacheExpire app
	AppCacheExpire = getEnv("APP_CACHE_EXPIRE", "120s")
	// Version app
	Version = "BUILD"
	// Environment app
	Environment = ""
	// APIReg is a api map
	APIReg = make(map[string]string)
	// NC nats broker
	NC *nats.Conn
	// EC nats encoded
	EC *nats.EncodedConn
	// DB mysql conn
	DB *sql.DB
	// CACHE redis conn
	CACHE *redis.Client
	// Cache param
	Cache *string
	//REQ0 counter
	REQ0 float64
	//REQ1 counter
	REQ1 float64
	// INM metrics
	INM *metrics.InmemSink
	// STMTIns insert
	STMTIns *sql.Stmt
	// STMTSel select
	STMTSel *sql.Stmt
	// Role application name
	Role = ""
	// Wait timeout
	Wait *string

	imageFilename   string
	ratio           float64
	fixedWidth      int
	fixedHeight     int
	fitScreen       bool
	stretchedScreen bool
	colored         bool
	reversed        bool

	convertDefaultOptions = convert.DefaultOptions
)
