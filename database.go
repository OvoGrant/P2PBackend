package main

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
)

// activePeers refers to the set containing the active peers and their information
const activePeers = "active_peers"

// redisOptions contains the information for connecting to redis
var redisOptions = &redis.Options{
	Addr:             "localhost:6379",
	Password:         "",
	DB:               0,
	DisableIndentity: true,
}

// Client provides access to local installation of redis
var Client *redis.Client

// initClient creates a new redis client
func initClient() {
	Client = redis.NewClient(redisOptions)
}

// SetPeerActive this function takes in peer information and adds it to the active peer set
func SetPeerActive(peer string) bool {

	res := Client.SAdd(context.Background(), activePeers, peer)

	return res.Val() == 1
}

//AddPeerToFile takes in a peer and a filename and indicates that the peer has the file
func AddPeerToFile(peer string, file string) bool {
	res := Client.SAdd(context.Background(), file, peer)

	return res.Val() == 1
}

//GetPeersWithFile returns the set intersection between active peers and the set of peers with a file
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

	//we will return only the first 10 peers for the sake of brevity
	return val.Val(), nil
}

//ClearPeerSet empties the active peer set
func ClearPeerSet() {
	Client.Del(context.Background(), activePeers)
}
