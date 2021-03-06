package main

import (
	"fmt"
	"github.com/ddbelyaev/gomeans/kmeans"
	"math/rand"
	"time"
)

func main() {

	km := kmeans.KMeans{K: 10, Iter: 500, Centroids: map[int][]float64{}}

	var rows int = 5500
	var columns int = 20

	var data [][]float64 = make([][]float64, rows)

	for j := 0; j < rows; j++ {
		for i := 0; i < columns; i++ {
			data[j] = append(data[j], float64(rand.Intn(1000)))
		}
	}

	start := time.Now()

	km.Fit(data)

	elapsed := time.Since(start)
	fmt.Printf("%s time took\n", elapsed)
}
