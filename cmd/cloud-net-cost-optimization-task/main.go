package main

import "fmt"

// ====================================================================
//  PROBLEM STATEMENT: CLOUD BALANCING
// ====================================================================

//  GIVEN:
//    1. A network consisting of cloud infrastructure resources, modeled as
//       a directed graph G = (V, E).
//    2. A vector 'b' representing the infrastructure workload:
//       - b[i] > 0: A route with a shortage of computing resources.
//       - b[i] < 0: A route with an excess of computing resources.
//       - b[i] = 0: A balanced route.
//       * Modeling constraint: Sum([i in V], b[i]) = 0, ensuring consistency across all routes.
//    3. A cost vector 'C' of size |E|, where C[e] is the tariff cost
//       for routing along edge 'e'.

//  REQUIRED TO FIND:
//    The optimal workload distribution across routes x = (x_1, x_2, ..., x_n), n=|E|, Ax = b,
//    min F(x) <- Sum_{e in E} ( C[e] * (x[e])^2 )

func main() {
	fmt.Println("Comparison: OPTIMIZATION METHOD WITH INITIAL DECOMPOSITION || METHOD OF LAGRANGE MULTIPLIERS")

	inputParameters := []struct{ vertexCount, edgeCount int }{
		{100, 300},
		{200, 600},
		{300, 900},
		{400, 1200},
		{500, 1500},
		{600, 1800},
	}

	results := make([]SampleResult, 0, len(inputParameters))

	for _, inputParameter := range inputParameters {
		fmt.Printf("[BENCHMARK] Testing dimension: |V| = %d, |E| = %d...\n", inputParameter.vertexCount, inputParameter.edgeCount)
		res := runBenchmarkIteration(inputParameter.vertexCount, inputParameter.edgeCount)
		results = append(results, res)
	}

	printReport(results)
}
