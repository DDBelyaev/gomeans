package kmeans

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/google/go-cmp/cmp"
)

type KMeans struct {
	K         int
	Iter      int
	Centroids map[int][]float64
	Clusters  []int
}

func distance(row1, row2 []float64) (float64, error) {
	var dist float64

	if len(row1) != len(row2) {
		return dist, errors.New("Distance is fucked up")
	}

	for i := 0; i < len(row1); i++ {
		dist += math.Pow((row1[i] - row2[i]), 2)
	}

	return dist, nil
}

func array_equal(arr1, arr2 []float64) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

func array_equal_int(arr1, arr2 []int) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}

func init_centroids(data *[][]float64, k int) (map[int][]float64, error) {
	centroids := make(map[int][]float64, k)

	if len(*data) < k {
		return centroids, errors.New("NOT ENOUGH DATA")
	}

	for i := 0; i < k; i++ {
		centroids[i] = (*data)[rand.Intn(len(*data))]
		for j := 0; j < i; j++ {
			fmt.Printf("%v, %v\n", i, j)
			if cmp.Equal(centroids[i], centroids[j]) {
				centroids[i] = (*data)[rand.Intn(len(*data))]
				j = -1
			}
		}
	}

	//	fmt.Println(centroids)
	fmt.Println("CENTROIDS INITIALIZED")
	return centroids, nil
}

func array_mean(aoa [][]float64) ([]float64, error) {
	arr := make([]float64, len(aoa[0]))

	for i := 0; i < len(aoa[0]); i++ {
		var mean float64

		for j := 0; j < len(aoa)/2; j++ {
			mean += aoa[j][i] + aoa[len(aoa)-1-j][i]
		}

		mean = mean / float64(len(aoa))
		arr[i] = mean
	}

	return arr, nil
}

func calculate_clusters(centroids map[int][]float64, data *[][]float64) []int {
	clusters := make([]int, len(*data))
	var wg sync.WaitGroup
	wg.Add(len(*data))
	for i, row := range *data {
		// THIS IS IT :D
		go func(i int, row []float64) {
			defer wg.Done()
			var cluster int
			var dist float64 = math.Inf(1)
			for cl, centroid := range centroids {
				if d, _ := distance(row, centroid); d < dist {
					dist = d
					cluster = cl
					//				fmt.Printf("DISTANCE %v, CLUSTER %v\n", d, cl)
				} else {
					//				fmt.Printf("CLUSTER %v OMMITED\n", cl)
				}
			}
			//		fmt.Printf("FINAL CLUSTER %v\n", cluster)
			clusters[i] = cluster
		}(i, row)
	}

	fmt.Println("CALCULATING NEW CLUSTERS...")
	wg.Wait()
	return clusters
}

func calculate_centroids(data *[][]float64, clusters []int, k int) map[int][]float64 {
	centroids := make(map[int][]float64, k)

	for i := 0; i < k; i++ {
		var aoa [][]float64
		for j := 0; j < len(clusters); j++ {
			if clusters[j] == i {
				aoa = append(aoa, (*data)[j])
			}
		}

		centroids[i], _ = array_mean(aoa)
	}

	return centroids
}

func (km *KMeans) Fit(data *[][]float64) {
	centroids, _ := init_centroids(data, km.K)
	clusters := calculate_clusters(centroids, data)
	var old_clusters []int
	var max_it, total_it int

	for it := 0; it < km.Iter; it++ {
		centroids = calculate_centroids(data, clusters, km.K)
		old_clusters, clusters = clusters, calculate_clusters(centroids, data)

		if cmp.Equal(old_clusters, clusters) {
			max_it++
			if max_it == 5 {
				total_it++
				break
			}
		} else {
			max_it = 0
		}

		total_it++
	}

	km.Centroids = centroids
	km.Clusters = clusters

	//fmt.Printf("CENTROIDS %+v", km.Centroids)
	fmt.Printf("CLUSTERS %+v\n", km.Clusters[:50])
	fmt.Printf("TOTAL ITERATIONS %v\n", total_it)
}
