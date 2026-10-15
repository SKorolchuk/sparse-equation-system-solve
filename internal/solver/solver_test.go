package networksolver

import (
	"math"
	"strings"
	"testing"
)

func TestSolveParticular_TableDriven(t *testing.T) {
	type testCase struct {
		name        string
		vertexCount int
		arcs        []Arc
		b           []float64
		root        int
		expectError bool
		errContains string
	}

	tests := []testCase{
		{
			name:        "Positive case: Basic linear chain (3 nodes, 2 edges)",
			vertexCount: 3,
			arcs: []Arc{
				{ID: 0, From: 1, To: 2},
				{ID: 1, From: 2, To: 3},
			},
			b:           []float64{0.0, 10.0, 0.0, -10.0},
			root:        1,
			expectError: false,
		},
		{
			name:        "Positive case: Complex highly connected graph (9 nodes, 17 edges)",
			vertexCount: 9,
			arcs: []Arc{
				{ID: 0, From: 1, To: 2}, {ID: 1, From: 1, To: 3},
				{ID: 2, From: 2, To: 4}, {ID: 3, From: 2, To: 5},
				{ID: 4, From: 3, To: 6}, {ID: 5, From: 3, To: 7},
				{ID: 6, From: 5, To: 8}, {ID: 7, From: 7, To: 9},
				{ID: 8, From: 4, To: 5}, {ID: 9, From: 6, To: 7},
				{ID: 10, From: 8, To: 1}, {ID: 11, From: 9, To: 2},
				{ID: 12, From: 4, To: 1}, {ID: 13, From: 5, To: 6},
				{ID: 14, From: 8, To: 4}, {ID: 15, From: 9, To: 6},
				{ID: 16, From: 3, To: 2},
			},
			b:           []float64{0.0, 2.0, 4.0, -6.0, -10.0, 10.0, 9.0, -1.0, -3.0, -5.0},
			root:        5,
			expectError: false,
		},
		{
			name:        "Negative case: Vector b size violation (len(b) <= VertexCount)",
			vertexCount: 4,
			arcs: []Arc{
				{ID: 0, From: 1, To: 2},
				{ID: 1, From: 2, To: 3},
				{ID: 2, From: 3, To: 4},
			},
			b:           []float64{0.0, 5.0, -5.0},
			root:        1,
			expectError: true,
			errContains: "vector b size",
		},
		{
			name:        "Negative case: Network balance constraints violation (sum of b_i != 0)",
			vertexCount: 3,
			arcs: []Arc{
				{ID: 0, From: 1, To: 2},
				{ID: 1, From: 2, To: 3},
			},
			b:           []float64{0.0, 10.0, 5.0, -5.0},
			root:        1,
			expectError: true,
			errContains: "unbalanced network",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGraph(tc.vertexCount, tc.arcs)
			solution, err := g.SolveParticular(tc.b, tc.root)

			// Error handling verification
			if tc.expectError {
				if err == nil {
					t.Fatalf("Validation failed: expected error intercept, but computation completed successfully")
				}
				if !strings.Contains(err.Error(), tc.errContains) {
					t.Errorf("Error message mismatch:\nExpected substring: %q\nGot error: %q", tc.errContains, err.Error())
				}
				return
			}

			// Success handling verification
			if err != nil {
				t.Fatalf("Unexpected failure: %v", err)
			}

			if ok, report := VerifySolution(g, solution.Flows, tc.b); !ok {
				t.Errorf("Mathematical verification failed! Network balance constraints violated:%s", report)
			}
		})
	}
}

func TestComputeGeneralSolutionMatrix_Orthogonality(t *testing.T) {
	arcs := []Arc{
		{ID: 0, From: 1, To: 2},
		{ID: 1, From: 2, To: 3},
		{ID: 2, From: 3, To: 4},
		{ID: 3, From: 4, To: 1},
		{ID: 4, From: 2, To: 4},
	}

	g := NewGraph(4, arcs)
	tree := g.BuildSpanningTree(1)

	deltaMatrix := tree.ComputeGeneralSolutionMatrix(len(arcs))
	expectedCycles := len(arcs) - (g.VertexCount - 1)

	if len(deltaMatrix) != expectedCycles {
		t.Fatalf("Cycle space dimension mismatch: expected %d, got %d", expectedCycles, len(deltaMatrix))
	}

	// Numerical verification of basis vector orthogonality
	for rowIndex, cycleVector := range deltaMatrix {
		if len(cycleVector) != len(arcs) {
			t.Fatalf("Cycle vector %d length corrupted: expected %d components", rowIndex, len(arcs))
		}

		// Validate cycle coefficients against each node balance equation
		for node := 1; node <= g.VertexCount; node++ {
			nodeBalanceInCycle := 0.0
			for _, arc := range g.Arcs {
				if arc.From == node {
					nodeBalanceInCycle += cycleVector[arc.ID] // Outflow (+)
				}
				if arc.To == node {
					nodeBalanceInCycle -= cycleVector[arc.ID] // Inflow (-)
				}
			}

			if math.Abs(nodeBalanceInCycle) > 1e-9 {
				t.Errorf("Subspace orthogonality failure: cycle %d broken at node %d. Loop balance residual = %.6f", rowIndex, node, nodeBalanceInCycle)
			}
		}
	}
}

func TestSolution_StringFormatting(t *testing.T) {
	arcs := []Arc{
		{ID: 0, From: 1, To: 2},
	}
	g := NewGraph(2, arcs)
	sol := Solution{
		Graph: g,
		Flows: []float64{45.1234},
	}

	output := sol.String()

	expectedMarkers := []string{
		"Network Flows (Particular Solution):",
		"Edge 0 (1 -> 2): Flow = 45.1234",
	}

	for _, marker := range expectedMarkers {
		if !strings.Contains(output, marker) {
			t.Errorf("String layout adapter analysis failed. Missing token %q in output string. Received:\n%s", marker, output)
		}
	}
}

func TestVerifySolution_OutofBoundsProtection(t *testing.T) {
	arcs := []Arc{
		{ID: 0, From: 5, To: 2}, // Node 5 is out of bounds for a two-element balance vector b
	}
	g := NewGraph(5, arcs)
	b := []float64{0.0, 0.0}
	x := []float64{1.0}

	ok, report := VerifySolution(g, x, b)
	if ok {
		t.Errorf("Safety execution fail: verifier skipped incorrect slice boundaries without raising an error report")
	}
	if !strings.Contains(report, "error") {
		t.Errorf("Verifier diagnostic notification mismatch: %q", report)
	}
}

func BenchmarkSolveParticular_9Nodes17Edges(b *testing.B) {
	arcs := []Arc{
		{ID: 0, From: 1, To: 2}, {ID: 1, From: 1, To: 3},
		{ID: 2, From: 2, To: 4}, {ID: 3, From: 2, To: 5},
		{ID: 4, From: 3, To: 6}, {ID: 5, From: 3, To: 7},
		{ID: 6, From: 5, To: 8}, {ID: 7, From: 7, To: 9},
		{ID: 8, From: 4, To: 5}, {ID: 9, From: 6, To: 7},
		{ID: 10, From: 8, To: 1}, {ID: 11, From: 9, To: 2},
		{ID: 12, From: 4, To: 1}, {ID: 13, From: 5, To: 6},
		{ID: 14, From: 8, To: 4}, {ID: 15, From: 9, To: 6},
		{ID: 16, From: 3, To: 2},
	}
	balances := []float64{0.0, 2.0, 4.0, -6.0, -10.0, 10.0, 9.0, -1.0, -3.0, -5.0}
	g := NewGraph(9, arcs)

	for b.Loop() {
		// Benchmark runtime for the bottom-up tree flow balancing algorithm
		_, _ = g.SolveParticular(balances, 5)
	}
}
