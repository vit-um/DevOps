package main

import (
	"bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"hash/fnv"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CrowdSurge/banner"
	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/den-vasyliev/image2ascii/convert"
)

func initOptions() {

	flag.StringVar(&imageFilename,
		"f",
		"",
		"Image filename to be convert")
	flag.Float64Var(&ratio,
		"r",
		convertDefaultOptions.Ratio,
		"Ratio to scale the image, ignored when use -w or -g")
	flag.IntVar(&fixedWidth,
		"w",
		convertDefaultOptions.FixedWidth,
		"Expected image width, -1 for image default width")
	flag.IntVar(&fixedHeight,
		"g",
		convertDefaultOptions.FixedHeight,
		"Expected image height, -1 for image default height")
	flag.BoolVar(&fitScreen,
		"s",
		convertDefaultOptions.FitScreen,
		"Fit the terminal screen, ignored when use -w, -g, -r")
	flag.BoolVar(&colored,
		"c",
		convertDefaultOptions.Colored,
		"Colored the ascii when output to the terminal")
	flag.BoolVar(&reversed,
		"i",
		convertDefaultOptions.Reversed,
		"Reversed the ascii when output to the terminal")
	flag.BoolVar(&stretchedScreen,
		"t",
		convertDefaultOptions.StretchedScreen,
		"Stretch the picture to overspread the screen")
}

func usage() {
	log.Printf("Usage: app [-name name] [-role role] [-port port] \n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func natsErrHandler(NC *nats.Conn, sub *nats.Subscription, natsErr error) {
	fmt.Printf("error: %v\n", natsErr)
	if natsErr == nats.ErrSlowConsumer {
		pendingMsgs, _, err := sub.Pending()
		if err != nil {
			fmt.Printf("couldn't get pending messages: %v", err)
			return
		}
		fmt.Printf("Falling behind with %d pending messages on subject %q.\n",
			pendingMsgs, sub.Subject)
		// Log error, notify operations...
	}
	// check for other errors
}

func file() {
	router := mux.NewRouter()
	path := flag.String("p", "/static/", "path to serve static files")
	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()

	router.PathPrefix(*path).Handler(http.StripPrefix(*path, http.FileServer(http.Dir(*directory))))
	//http.Handle("/static/", http.StripPrefix(strings.TrimRight(path, "/"), http.FileServer(http.Dir(*directory))))
	log.Fatal(http.ListenAndServe(":8880", router))

}

func hash(decodedStr string) (uint32, string) {
	///defer metrics.MeasureSince([]string{"API"}, time.Now())
	//log.Print("DecodedStr: ", decodedStr)
	encodedStr := hex.EncodeToString([]byte(banner.PrintS(decodedStr)))
	//log.Print("EncodedStr: ", encodedStr)
	h := fnv.New32a()
	h.Write([]byte(encodedStr))
	//hashStr := fmt.Sprintf("%x", md5.Sum([]byte(encodedStr)))
	return h.Sum32(), encodedStr
}

func rest(url string, jsonStr string) []byte {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	log.Print("response Status:", resp.Status)
	log.Print("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	return body
}

func parseOptions() (*convert.Options, error) {

	// config  the options
	convertOptions := &convert.Options{
		Ratio:           ratio,
		FixedWidth:      fixedWidth,
		FixedHeight:     fixedHeight,
		FitScreen:       fitScreen,
		StretchedScreen: stretchedScreen,
		Colored:         colored,
		Reversed:        reversed,
	}
	return convertOptions, nil
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	opts = append(opts, nats.ErrorHandler(natsErrHandler))
	return opts
}

// getEnv get key environment variable if exist otherwise return defalutValue
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

//return []byte("0")
/*
	cached, err := CACHE.Get(strconv.FormatUint(uint64(hashStr), 10)).Result()

	if *Cache == "false" {
		err = errors.New("Processing")
		cached = "636163686564"
	}

	if err != nil {
		log.Print(err)

		sec, _ := time.ParseDuration(AppCacheExpire)

		CACHE.Set(strconv.FormatUint(uint64(hashStr), 10), encodedStr, sec)

		msg, err := NC.Request(AppDatastore+".hash", []byte(fmt.Sprintf(`{"hash":"%s"}`, strconv.FormatUint(uint64(hashStr), 10))), 2*time.Second)
		if err != nil {
			log.Printf("ErrRequest: %e", err)
		}

		var reply []byte

		if err != nil {
			log.Printf("ErrReply: %e", err)
			reply = []byte(fmt.Sprintf("{Reply:%s}", err))
		} else {
			reply = msg.Data
		}

		REQ0 = REQ0 + 1

		return reply
	}

	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Print(err)
		return []byte("undef")
	}

	return []byte(string(decoded))
*/

/*
	case "POST":
		b, _ := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
		if err := json.Unmarshal(b, &m); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			if err := json.NewEncoder(w).Encode(err); err != nil {
				panic(err)
			}
		}
		w.Write([]byte(dataStore(m.Hash)))*/
