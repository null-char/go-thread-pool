package main

import (
	"errors"
	"hash/fnv"
	"math/rand"
	"time"

	j "github.com/null-char/go-thread-pool/job"
	"github.com/null-char/go-thread-pool/pool"
)

// These functions and stuff are just here so that we can get some dummy jobs going.

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// GenerateRandStr generates a random string of size n.
func GenerateRandStr(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

// Process takes in a string, hashes it and returns the value.
func Process(r j.Resource) (uint32, error) {
	h := fnv.New32a()
	h.Write([]byte(r))
	time.Sleep(time.Second)

	// We'll just return nil and no error.
	return h.Sum32(), nil
}

// PostProcess takes in a Result, and does some (simulated) extra processing on it.
func PostProcess(r j.Result) j.Result {
	// Sleep for a few milliseconds to simulate some work that takes a while.
	if sleepDur, err := time.ParseDuration("15ms"); err == nil {
		time.Sleep(sleepDur)

		// We're simply just going to increment the value.
		return j.Result{Job: r.Job, Value: r.Value + 1, Err: nil}
	}

	r.Err = errors.New("parsing duration for sleep while post-processing failed")
	return r
}

// Ideally numWorkers should be lower than numResources
const numWorkers = 73
const numResources = 150

func main() {
	// Build resources
	var resources []j.Resource = make([]j.Resource, numResources)
	for i := 0; i < numResources; i++ {
		resources[i] = j.Resource(GenerateRandStr(8))
	}

	p := pool.New(numWorkers)
	p.BeginWork(resources, Process, PostProcess)
}
