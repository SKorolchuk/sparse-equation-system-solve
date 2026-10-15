package main

import (
	"fmt"
	"math"
)

// solveLinearSystem solves a system of linear algebraic equations using Gaussian elimination with partial pivoting.
func solveLinearSystem(matrix [][]float64, vector []float64) []float64 {
	size := len(vector)
	for i := range size {
		maxElement := math.Abs(matrix[i][i])
		maxRow := i
		for k := i + 1; k < size; k++ {
			if a := math.Abs(matrix[k][i]); a > maxElement {
				maxElement = a
				maxRow = k
			}
		}
		matrix[i], matrix[maxRow] = matrix[maxRow], matrix[i]
		vector[i], vector[maxRow] = vector[maxRow], vector[i]

		for k := i + 1; k < size; k++ {
			ratio := -matrix[k][i] / matrix[i][i]
			for j := i; j < size; j++ {
				if i == j {
					matrix[k][j] = 0
					continue
				}
				matrix[k][j] += ratio * matrix[i][j]
			}
			vector[k] += ratio * vector[i]
		}
	}

	solution := make([]float64, size)
	for i := size - 1; i >= 0; i-- {
		solution[i] = vector[i] / matrix[i][i]
		for k := i - 1; k >= 0; k-- {
			vector[k] -= matrix[k][i] * solution[i]
		}
	}
	return solution
}

// printReport formats and prints the results of the performance benchmarks.
func printReport(results []SampleResult) {
	fmt.Println("\n=========================================================================================================================")
	fmt.Println("                                           EFFICIENCY SUMMARY FOR METHODS")
	fmt.Println("=========================================================================================================================")

	fmt.Println("| Vertices |V| |  Edges |E| | Matrix Solution Kernel Size | Time (Decomposition) | Time (Lagrange Method) |  Speedup  | Max Deviation |")
	fmt.Println("| :----------: | :--------: | :-------------------------: | :------------------: | :--------------------: | :-------: | :-----------: |")

	for _, result := range results {
		method1FmtStr := fmt.Sprintf("%v", result.Decomposition)
		method2FmtStr := fmt.Sprintf("%v", result.Lagrange)

		fmt.Printf("| %12d | %10d | %27d | %20s | %22s | %8.2fx | %13f |\n",
			result.Vertices,
			result.Edges,
			result.Chords,
			method1FmtStr,
			method2FmtStr,
			result.Speedup,
			result.MaxDeviation,
		)
	}
}
