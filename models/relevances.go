package models

import (
	"time"
	"strconv"
	"../config"
	"../utils/log"

	"github.com/garyburd/redigo/redis"
)

var simbasePool *redis.Pool

func initRelevances() {
	simbasePool = &redis.Pool{
		MaxIdle:   3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.SimbaseHost + ":" + config.SimbasePort)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	initRelevancesBasis()
}

func initRelevancesBasis() {
	simbaseClient := simbasePool.Get()
	defer simbaseClient.Close()

	blist, err := redis.Strings(simbaseClient.Do("blist"))
	if err != nil {
		log.Fatal(err)
		return
	}
	found := false
	for _, b := range blist {
		if b == "b512" {
			found = true
		}
	}
	if !found {
		tags := make([]string, 512)
		for i := 0; i < 512; i++ {
			tags[i] = "b" + strconv.Itoa(i)
		}
		_, err := simbaseClient.Do("bmk", redis.Args{}.Add("b512").AddFlat(tags)...)
		if err != nil {
			log.Fatal(err)
			return
		}

		_, err = simbaseClient.Do("vmk", "b512", "program")
		if err != nil {
			log.Fatal(err)
			return
		}

		_, err = simbaseClient.Do("rmk", "program", "program", "jensenshannon")
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func delRelevances() {
	simbasePool.Close()
}

func AddVector(id int, vector []float64) error {
	simbaseClient := simbasePool.Get()
	defer simbaseClient.Close()

	_, err := simbaseClient.Do("vget", "program", id)
	if err != nil {
		_, err = simbaseClient.Do("vset", redis.Args{}.Add("program").Add(id).AddFlat(vector)...)
	} else {
		_, err = simbaseClient.Do("vadd", redis.Args{}.Add("program").Add(id).AddFlat(vector)...)
	}
	return err
}

func getRelatedProgramsWithRelevances(id int, number int) ([]Program, error) {
	simbaseClient := simbasePool.Get()
	defer simbaseClient.Close()

	res, err := redis.Ints(simbaseClient.Do("rrec", "program", id, "program"))
	if err != nil {
		return nil, err
	}

	result := make([]Program, number)

	for i, p := range res {
		if i >= number {
			break
		}

		var program Program
		if program.Load(p) != nil {
			continue
		}

		result = append(result, program)
	}

	return result, nil
}
