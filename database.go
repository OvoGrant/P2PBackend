package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const activePeers = "active_peers"

var redisOptions = &redis.Options{
	Addr:             "localhost:6379",
	Password:         "",
	DB:               0,
	DisableIndentity: true,
}

var Client *redis.Client

func initClient() {
	Client = redis.NewClient(redisOptions)
}

func SetPeerActive(peer string) bool {

	res := Client.SAdd(context.Background(), activePeers, peer)

	return res.Val() == 1
}

func AddPeerToFile(peer string, file string) bool {
	res := Client.SAdd(context.Background(), file, peer)

	return res.Val() == 1
}

func GetPeersWithFile(file string) ([]string, error) {
	//check if the file exists first of all
	exists := Client.Exists(context.Background(), file)

	result, err := exists.Result()

	//if no file is found return this
	if err != nil {
		return nil, err
	}

	if result == 0 {
		return nil, errors.New("file not found")
	}

	//if the file is found then return this
	val := Client.SInter(context.Background(), file, activePeers)

	fmt.Println(val.Val())

	//we will return only the first 10 peers for the sake of brevity
	return val.Val(), nil
}

func ClearPeerSet() {
	Client.Del(context.Background(), activePeers)
}
