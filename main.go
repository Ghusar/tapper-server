package main

import (
	"os"
	"net/http"
	"context"
    "github.com/go-redis/redis/v9"
	"encoding/json"
)

var ctx = context.Background()

var redisClient = redis.NewClient(&redis.Options{
	Addr:      os.Getenv("REDIS_URL"),
	Password:  os.Getenv("REDIS_PASS"),
	DB:       0,  // use default DB
})

type add_score_request_body struct {
    Id string
	Score float64
}

type score_response struct {
	Id string
	Score float64
}

func main() {
	http.HandleFunc("/add", handleAddScore)
	http.HandleFunc("/get", handleGetAllScores)
	http.ListenAndServe(":8080", nil)
}

func handleAddScore(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var body add_score_request_body
	decodeErr := decoder.Decode(&body)

	if decodeErr != nil {
		panic(decodeErr)

		return
	}
	
	err := redisClient.ZAdd(ctx, "test_lb", redis.Z{Score: body.Score, Member: body.Id}).Err()

	if err != nil {
		panic(err)
	}
}


func handleGetAllScores(w http.ResponseWriter, r *http.Request) {
	vals, err := redisClient.ZRevRangeByScoreWithScores(ctx, "test_lb", &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result()

    if err != nil {
		panic(err)

		return
	}

	data := make([]score_response, len(vals))
	for index := 0; index < len(vals); index++ {
		data[index] = score_response{Id: vals[index].Member.(string), Score: vals[index].Score }
	} 
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}