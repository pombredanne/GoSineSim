package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"sort"
)

type Item struct {
	Id   string
	Data map[string]float64
}

type Items []Item

type Result struct {
	Similarity float64
	Id         string
	Data       map[string]float64
}

type GoSignSimResults []Result

func (slice GoSignSimResults) Len() int {
	return len(slice)
}

func (slice GoSignSimResults) Less(i, j int) bool {
	return slice[i].Similarity < slice[j].Similarity
}

func (slice GoSignSimResults) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func norm(obj Item) float64 {
	var norm float64 = 0

	for _, v := range obj.Data {
		norm += v * v
	}

	return math.Sqrt(norm)
}

func dotProduct(source, other Item) float64 {
	var product float64 = 0

	for k, v := range source.Data {
		product += v * other.Data[k]
	}

	return product
}

func pad(source, other Item) (Item, Item) {
	for k, _ := range source.Data {
		_, okay := other.Data[k]

		if okay == false {
			other.Data[k] = 0
		}
	}

	for k, _ := range other.Data {
		_, okay := source.Data[k]

		if okay == false {
			source.Data[k] = 0
		}
	}

	return source, other
}

func getScore(source, other Item, c chan float64) {
	source, other = pad(source, other)
	dem := norm(source) * norm(other)

	if dem > 0 {
		c <- (dotProduct(source, other) / dem) * 100
	} else {
		c <- 0
	}
}

func CoseineSimilarity(source Item, pool Items, threshold float64) GoSignSimResults {
	var results = make(GoSignSimResults, len(pool))

	for i, other := range pool {
		go func(i int, other Item) {
			score_c := make(chan float64, 1)

			getScore(source, other, score_c)

			score := <-score_c

			if score >= threshold {
				res := Result{Similarity: score, Data: other.Data, Id: other.Id}
				results[i] = res
			}
		}(i, other)
	}

	sort.Sort(sort.Reverse(results))

	return results
}

func main() {
	source := flag.String("source", "", "The source JSON object to compare")
	pool := flag.String("pool", "", "The data that will be compared against to source")
	pool_file := flag.String("pool_file", "", "An optional file to read the pool data from")
	threshold := flag.Float64("threshold", 0.0, "The lower limit ")

	flag.Parse()

	var obj Item
	string_bytes := []byte(*source)
	err := json.Unmarshal(string_bytes, &obj)

	if err != nil {
		log.Fatal(err)
	}

	pf := *pool_file
	var pool_obj Items

	if pf != "" {
		pool_file_bytes, pool_file_err := ioutil.ReadFile(string(*pool_file))

		if pool_file_err != nil {
			fmt.Print(pool_file_err)
		}

		pool_err := json.Unmarshal(pool_file_bytes, &pool_obj)

		if pool_err != nil {
			log.Fatal(pool_err)
		}
	} else {
		pool_bytes := []byte(*pool)
		pool_err := json.Unmarshal(pool_bytes, &pool_obj)

		if pool_err != nil {
			log.Fatal(pool_err)
		}
	}

	results := CoseineSimilarity(obj, pool_obj, float64(*threshold))
	results_json, results_err := json.Marshal(results)

	if results_err != nil {
		log.Fatal(results_err)
	}

	fmt.Printf("%s", results_json)
}
