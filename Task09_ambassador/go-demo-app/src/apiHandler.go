package main

import (
	"encoding/hex"
	"errors"
	"hash/fnv"
	_ "image/jpeg"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/valyala/fasthttp"
	//metrics "github.com/armon/go-metrics"
)

func api(ctx *fasthttp.RequestCtx) {
	//	fmt.Fprintf(ctx, "Hi there! RequestURI is %q", ctx.RequestURI())
	//}

	//func api(w http.ResponseWriter, r *http.Request) {
	// increment counter
	REQ0 = REQ0 + 1
	var reply []byte
	var hexEncodedStr, cached string
	var token uint32
	h := fnv.New32a()

	// parse uri
	u, err := url.Parse(string(ctx.RequestURI()))
	if err != nil {
		log.Print(err)
	}

	// get uri parameters
	q := u.Query()

	if len(q.Get("text")) > 0 {
		// we won't cache
		if *Cache == "false" {
			// create bin hash from text+request_siq
			h.Write([]byte(q.Get("text") + strconv.FormatFloat(REQ0, 'f', 0, 32)))
			// we need token as a string
			tokenStr := strconv.FormatUint(uint64(h.Sum32()), 10)
			// bin token
			token = h.Sum32()
			// encode text to hex
			hexEncodedStr = hex.EncodeToString([]byte(q.Get("text") + strconv.FormatFloat(REQ0, 'f', 0, 32)))
			// need this for the next check
			err = errors.New("NoCache")
			// define reply
			cached = tokenStr
			// default cache check first
		} else {
			h.Write([]byte(q.Get("text")))
			tokenStr := strconv.FormatUint(uint64(h.Sum32()), 10)
			token = h.Sum32()
			hexEncodedStr = hex.EncodeToString([]byte(q.Get("text")))
			cached, err = CACHE.Get(tokenStr).Result()

		}
		// if cache found - reply
		if err == nil {
			reply, err = hex.DecodeString(cached)
			ctx.Write(reply)
			// if cache not found - processing
		} else {
			// Create a unique subject name for replies.
			uniqueReplyTo := nats.NewInbox()

			// Listen for a single response
			sub, err := NC.SubscribeSync(uniqueReplyTo)
			if err != nil {
				log.Print(err)
			}
			// Send the request.
			// If processing is synchronous, use Request() which returns the response message.
			if err := EC.Publish("ascii.json.banner", &Req{Token: token, Hextr: hexEncodedStr, Reply: uniqueReplyTo, Db: q.Get("db")}); err != nil {
				log.Print(err)
			}
			// Read the reply for Wait seconds
			sec, _ := time.ParseDuration(*Wait)
			msg, err := sub.NextMsg(sec)
			if err != nil {
				log.Print(err)
			}
			// get result from data service
			cached, err := CACHE.Get(string(msg.Data)).Result()
			// decode cached reply
			reply, _ = hex.DecodeString(cached)

			ctx.Write(reply)

		}

	} else {
		ctx.Write(append([]byte(""), Environment...))
	}
}
