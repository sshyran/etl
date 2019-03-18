package main_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/gonum/stat/combin"
)

func try(K int, N int, X float64) float64 {
	sum := 0.0

	for k := 0; k <= N; k++ {
		combinations := combin.Binomial(N, k)                               // number of cases with k active tasks
		likelihood := math.Pow(X, float64(N-k)) * math.Pow(1-X, float64(k)) // of this combination
		if k < K {
			sum += float64(k) * float64(combinations) * likelihood // k processors busy
		} else {
			sum += float64(K) * float64(combinations) * likelihood // K processors busy.
		}
	}
	fmt.Printf("%d %5.2f, %d, %6.3f %5.2f\n", K, X, N, sum/float64(K), float64(K)/(float64(N)*(1.0-X)))
	return sum / float64(K)
}
func TestMain(t *testing.T) {

	K := 4 // Number of processors
	for X := 0.1; X < 0.81; X += 0.1 {
		for N := 6; ; N++ {
			cpu := try(K, N, X)
			if cpu > 0.95 {
				break
			}
		}
		fmt.Println()
	}

}
