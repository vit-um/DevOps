package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	//"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/armon/go-metrics"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nats-io/nats.go"
	"github.com/valyala/fasthttp"
)

// Req Define the object
type Req struct {
	Token uint32
	Hextr string
	Reply string
	Db    string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	INM = metrics.NewInmemSink(10*time.Second, time.Minute)
	sig := metrics.DefaultInmemSignal(INM)

	defer sig.Stop()

	API["ascii"] = "curl -XPOST --data '{text:TEXT}' HOST/ascii/"
	API["img"] = "curl -F 'image=@IMAGE' HOST/img/"
	API["ml5"] = "curl HOST/ml5/"
	API["data"] = "broker message api"

	initOptions()

	AppName := flag.String("name", "k8sdiy", "application name")
	AppRole := flag.String("role", "api", "app role: api data ascii img ml5 iot")
	AppPort := flag.String("port", "8080", "application port")
	//AppPath := flag.String("path", "/static/", "path to serve static files")
	//AppDir := flag.String("dir", "./ml5", "the directory of static files to host")
	//ModelsPath := flag.String("mpath", "/models/", "path to serve models files")
	//ModelsDir := flag.String("mdir", "./ml5/models", "the directory of models files to host")
	Cache = flag.String("cache", "true", "cache enable")
	Wait = flag.String("wait", "2s", "wait timeout")

	Urls = flag.String("server", nats.DefaultURL, "The nats server URLs (separated by comma)")
	//var userCreds = flag.String("creds", "", "User Credentials File")
	var showTime = flag.Bool("timestamp", false, "Display timestamps")
	//var queueName = flag.String("queGroupName", "K8S-NATS-Q", "Queue Group Name")
	var showHelp = flag.Bool("help", false, "Show help message")
	var err error

	log.SetFlags(0)
	flag.Usage = usage

	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	// Environment app
	Role = *AppRole

	metrics.NewGlobal(metrics.DefaultConfig(Role), INM)

	// Perf

	REQ0 = 0.0
	REQ1 = 0.0
	//t0 := time.Now()

	go func() { // Daniel told me to write this handler this way.
		for {
			select {
			case <-time.After(time.Second * 1):
			//	ts := time.Since(t0)
			//	log.Println("[", Role, "] time: ", ts, " requests: ", REQ0, " rps: ", (REQ0-REQ1)/1, " throughput:", float64(REQ0)/ts.Seconds())
				REQ1 = REQ0
			}
		}
	}()

	Environment = fmt.Sprintf("%s-%s:%s", *AppName, Role, Version)
	// Connect Options.

	log.Printf("Listening on [%s]: %s port: %s", *AppRole, Environment, *AppPort)

	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	// the corresponding fasthttp code
	router := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			img(ctx)
		case "/version":
			version(ctx)
		case "/healthz":
			//		healthz(ctx)
		case "/readinez":
			//		readinez.HandlerFunc(ctx)
		case "/img":
			img(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}

	/*
		router := pat.New()
		router.Get("/version", http.HandlerFunc(version))
		router.Get("/healthz", http.HandlerFunc(healthz))
		router.Get("/readinez", http.HandlerFunc(readinez))
	*/
	switch *AppRole {

	/*	case "iot":
		svr := &service.Server{
			KeepAlive:        300,           // seconds
			ConnectTimeout:   2,             // seconds
			SessionsProvider: "mem",         // keeps sessions in memory
			Authenticator:    "mockSuccess", // always succeed
			TopicsProvider:   "mem",         // keeps topic subscriptions in memory
		}
		// Listen and serve connections at localhost:1883
		err := svr.ListenAndServe("tcp://:1883")
		fmt.Printf("%v", err)
	*/
	case "api":
		//router.Get("/", http.HandlerFunc(api))

		//router.HandleFunc("/", api)

	case "ascii":
		//router.Get("/", http.HandlerFunc(ascii))

		//router.HandleFunc("/", ascii)

	case "img":
		//if err := NC.Publish("api."+Environment, []byte(API["img"])); err != nil {
	//		log.Fatal(err)
	//		}
	//router.HandleFunc("/", img)
	//router.Get("/", http.HandlerFunc(img))

	case "ml5":

		//router.PathPrefix(*AppPath).Handler(http.StripPrefix(*AppPath, http.FileServer(http.Dir(*AppDir))))
		//router.PathPrefix(*ModelsPath).Handler(http.StripPrefix(*ModelsPath, http.FileServer(http.Dir(*ModelsDir))))

		//router.HandleFunc("/", ml5)
		//router.Get("/", http.HandlerFunc(ml5))

	case "data":

		//_, err = DB.Exec("CREATE TABLE IF NOT EXISTS demo (id INT NOT NULL AUTO_INCREMENT, token INT UNSIGNED NOT NULL, text TEXT, PRIMARY KEY(id, token))")

		if err != nil {
			log.Printf("CreateErr: %s", err) // proper error handling instead of panic in your app
		}

		//router.HandleFunc("/", dataHandler)
		//router.Get("/", http.HandlerFunc(data))

	}

	log.Fatal(fasthttp.ListenAndServe(":"+*AppPort, router))

	//log.Fatal(http.ListenAndServe(":"+*AppPort, router))

	// Setup the interrupt handler to drain so we don't miss
	// requests when scaling down.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	NC.Drain()
	log.Fatalf("Exiting")

}

func cache() {
	// Connect to cache
	CACHE = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", AppCache, AppCachePort),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := CACHE.Ping().Result()
	if err != nil {
		log.Print(err)
	}
}

func db() {
	// Connect to db

	DB, err := sql.Open("mysql", AppDb)
	//	DB.SetMaxIdleConns(10000)
	if err != nil {
		log.Print(err)
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		log.Print(err) // proper error handling instead of panic in your app
	}

	STMTIns, err = DB.Prepare("insert into demo values(null,?,?)")
	STMTSel, err = DB.Prepare("SELECT text FROM demo WHERE token = ? limit 1")
	defer STMTIns.Close()
	defer STMTSel.Close()
}

func mq() {
	subj, subjJSON, i := *AppRole+".*", *AppRole+".json.*", 0

	opts := []nats.Option{nats.Name(Role + " on " + subj)}
	opts = setupConnOptions(opts)

	// Use UserCredentials
	/*if *userCreds != "" {
		opts = append(opts, nats.UserCredentials(*userCreds))
	}
	*/

	// Connect to NATS
	var err error
	NC, err = nats.Connect(*Urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	if err := NC.LastError(); err != nil {
		log.Fatal(err)
	}
	defer NC.Close()

	EC, err = nats.NewEncodedConn(NC, nats.JSON_ENCODER)
	defer EC.Close()

	// Subscribe
	if _, err = EC.Subscribe(subjJSON, func(r *Req) {

		REQ0 = REQ0 + 1

		i++

		if *AppRole == "ascii" {

			go AsciiHandler(r, i)

		} else if *AppRole == "data" {

			go DataHandler(r, i)
		}

	}); err != nil {
		log.Print(err)
	}

	log.Printf("Listening on [%s]: %s port: %s", subj, Environment, *AppPort)

}
