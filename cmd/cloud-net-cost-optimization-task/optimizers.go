package main

// runFlowOptimization performs least-squares flow optimization over the cycle space.
func runFlowOptimization(particular []float64, delta [][]float64, costs []float64) ([]float64, []float64) {
	edgeCount := len(particular)
	chordCount := len(delta)

	if chordCount == 0 {
		return particular, make([]float64, 0)
	}

	matrixS := make([][]float64, chordCount)
	for i := range matrixS {
		matrixS[i] = make([]float64, chordCount)
	}
	vectorR := make([]float64, chordCount)

	for i := range chordCount {
		for j := 0; j < chordCount; j++ {
			var sum float64
			for e := range edgeCount {
				sum += costs[e] * delta[i][e] * delta[j][e]
			}
			matrixS[i][j] = sum
		}

		var rhsSum float64
		for e := range edgeCount {
			rhsSum += costs[e] * particular[e] * delta[i][e]
		}
		vectorR[i] = -rhsSum
	}

	xUnits := solveLinearSystem(matrixS, vectorR)

	optimal := make([]float64, edgeCount)
	for e := range edgeCount {
		val := particular[e]
		for i := range chordCount {
			val += delta[i][e] * xUnits[i]
		}
		optimal[e] = val
	}

	return optimal, xUnits
}

// solveViaLagrangeDirect constructs an expanded Karush-Kuhn-Tucker (KKT) system
// using Lagrange multipliers to find the global cost optimum.
func solveViaLagrangeDirect(vertexCount int, edges []Edge, vectorB []float64, costs []float64) []float64 {
	edgeCount := len(edges)
	kktSize := edgeCount + vertexCount

	matrixKKT := make([][]float64, kktSize)
	for i := range matrixKKT {
		matrixKKT[i] = make([]float64, kktSize)
	}
	rhs := make([]float64, kktSize)

	for i := range edgeCount {
		matrixKKT[i][i] = 2.0 * costs[i]
	}

	for i, edge := range edges {
		matrixKKT[i][edgeCount+edge.From-1] = 1.0
		matrixKKT[i][edgeCount+edge.To-1] = -1.0

		matrixKKT[edgeCount+edge.From-1][i] = 1.0
		matrixKKT[edgeCount+edge.To-1][i] = -1.0
	}

	for i := 1; i <= vertexCount; i++ {
		rhs[edgeCount+i-1] = vectorB[i]
	}

	solution := solveLinearSystem(matrixKKT, rhs)
	return solution[:edgeCount]
}
