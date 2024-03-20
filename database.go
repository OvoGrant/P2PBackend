package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const ACTIVE_PEERS = "active_peers"

var redis_options = &redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
	DisableIndentity: true,
}

var Client *redis.Client

func initClient() {
	Client = redis.NewClient(redis_options)
}

func Close() {
	Client.Close()
}

func PeerActive(peer string) bool {
	res := Client.SIsMember(context.Background(), ACTIVE_PEERS, peer)

	return res.Val()
}
func SetPeerActive(peer string) bool {
	fmt.Println(peer)

	res := Client.SAdd(context.Background(), ACTIVE_PEERS, peer)

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

	fmt.Println("Result:", result)

	//if the file is found then return this
	val := Client.SInter(context.Background(), file, ACTIVE_PEERS)

	//we will return only the first 10 peers for the sake of brevity
	return val.Val(), nil
}
