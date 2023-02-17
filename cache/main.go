package main

import (
	"fmt"
	"sort"
)

type example struct {
	data      interface{}
	timestamp int64
}

func main() {
	m := map[string]example{
		"Bob":   example{data: "15", timestamp: 11},
		"Alice": example{data: "20", timestamp: 10},
		"Mark":  example{data: "30", timestamp: 5},
	}

	keys := []string{}
	for key, _ := range m {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return m[keys[i]].timestamp < m[keys[j]].timestamp
	})

	for _, key := range keys {
		fmt.Println(key, m[key])
	}

}
