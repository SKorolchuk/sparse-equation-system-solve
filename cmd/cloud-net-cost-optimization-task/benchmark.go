package main

import (
	"log"
	"math"
	"math/rand/v2"
	"time"

	"gonum.org/v1/gonum/graph/simple"
	"gonum.org/v1/gonum/mat"

	networksolver "sparse-equation-system-solve/internal/solver"
)

// Edge represents a directed connection between two vertices.
type Edge struct {
	From int
	To   int
}

// SampleResult holds the performance and precision metrics for a single benchmark iteration.
type SampleResult struct {
	Vertices      int
	Edges         int
	Chords        int
	Decomposition time.Duration
	Lagrange      time.Duration
	Speedup       float64
	MaxDeviation  float64
}

// computeMaxDeviation calculates the maximum absolute difference between two vectors.
func computeMaxDeviation(vecA, vecB []float64) float64 {
	if len(vecA) != len(vecB) {
		return math.MaxFloat64
	}
	var maxDiff float64
	for i := range vecA {
		diff := math.Abs(vecA[i] - vecB[i])
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	return maxDiff
}

// runBenchmarkIteration executes a comparative benchmark and measures accuracy.
func runBenchmarkIteration(vertexCount, edgeCount int) SampleResult {
	rng := rand.New(rand.NewPCG(42, 42))
	const rootNode int64 = 1

	graph := simple.NewDirectedGraph()
	for i := 1; i <= vertexCount; i++ {
		graph.AddNode(simple.Node(i))
	}

	edges := make([]Edge, 0, edgeCount)
	for i := 2; i <= vertexCount; i++ {
		parent := rng.IntN(i-1) + 1
		var edge Edge
		if rng.Float64() > 0.5 {
			edge = Edge{From: parent, To: i}
		} else {
			edge = Edge{From: i, To: parent}
		}
		graph.SetEdge(simple.Edge{F: simple.Node(edge.From), T: simple.Node(edge.To)})
		edges = append(edges, edge)
	}

	for len(edges) < edgeCount {
		u := rng.IntN(vertexCount) + 1
		v := rng.IntN(vertexCount) + 1
		if u == v || graph.HasEdgeBetween(int64(u), int64(v)) {
			continue
		}
		graph.SetEdge(simple.Edge{F: simple.Node(u), T: simple.Node(v)})
		edges = append(edges, Edge{From: u, To: v})
	}

	vectorBData := make([]float64, vertexCount+1)
	var totalSum float64
	for i := 1; i <= vertexCount; i++ {
		roll := rng.Float64()
		switch {
		case roll < 0.3:
			vectorBData[i] = float64(rng.IntN(4)+1) + rng.Float64()
		case roll < 0.6:
			vectorBData[i] = -(float64(rng.IntN(4)+1) + rng.Float64())
		default:
			vectorBData[i] = 0.0
		}
		totalSum += vectorBData[i]
	}
	vectorBData[vertexCount] -= totalSum
	gonumB := mat.NewVecDense(len(vectorBData), vectorBData)

	edgeCosts := make([]float64, edgeCount)
	for i := 0; i < edgeCount; i++ {
		edgeCosts[i] = rng.Float64()*1.4 + 0.1
	}

	// Benchmark 1: Decomposition Method
	startDecomposition := time.Now()
	adapter := networksolver.NewGonumAdapter(graph)

	particularSol, err := adapter.SolveParticular(gonumB, rootNode)
	if err != nil {
		log.Fatalf("Critical error: %v", err)
	}

	flowsData := make([]float64, edgeCount)
	for i := 0; i < edgeCount; i++ {
		flowsData[i] = particularSol.AtVec(i)
	}

	deltaMatrix := adapter.ComputeGeneralSolutionMatrix(rootNode)
	decompResult, _ := runFlowOptimization(flowsData, deltaMatrix, edgeCosts)
	decompositionDuration := time.Since(startDecomposition)

	// Benchmark 2: Lagrange Multipliers
	startLagrange := time.Now()
	lagrangeResult := solveViaLagrangeDirect(vertexCount, edges, vectorBData, edgeCosts)
	lagrangeDuration := time.Since(startLagrange)

	chordCount := len(deltaMatrix)
	speedup := float64(lagrangeDuration.Nanoseconds()) / float64(decompositionDuration.Nanoseconds())

	maxDeviation := computeMaxDeviation(decompResult, lagrangeResult)

	return SampleResult{
		Vertices:      vertexCount,
		Edges:         edgeCount,
		Chords:        chordCount,
		Decomposition: decompositionDuration,
		Lagrange:      lagrangeDuration,
		Speedup:       speedup,
		MaxDeviation:  maxDeviation,
	}
}
