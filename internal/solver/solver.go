package networksolver

import (
	"strings"
	"bytes"
	"fmt"
	"math"
)

// Solution is an wrapper for the flow vector (particular solution)
// used for verification and plain text output formatting
type Solution struct {
	Graph *Graph
	Flows []float64
}

// String formats the flows through the edges into a text table (implements fmt.Stringer)
func (s Solution) String() string {
	var buf bytes.Buffer
	buf.WriteString("Network Flows (Particular Solution):\n")
	for i, flow := range s.Flows {
		if i >= len(s.Graph.Arcs) {
			break
		}
		arc := s.Graph.Arcs[i]
		fmt.Fprintf(&buf, "  Edge %d (%d -> %d): Flow = %.4f\n", arc.ID, arc.From, arc.To, flow)
	}
	return buf.String()
}

// SolveParticular is a method that builds the spanning tree,
// calculates the particular solution, and verifies the network flow balance constraints
func (g *Graph) SolveParticular(vectorB []float64, root int) (Solution, error) {
	if len(vectorB) <= g.VertexCount {
		return Solution{}, fmt.Errorf("error: vector b size (%d) must be at least VertexCount+1 (%d) due to 1-based indexing", len(vectorB), g.VertexCount+1)
	}

	var sum float64
	for _, val := range vectorB {
		sum += val
	}
	if math.Abs(sum) > 1e-9 {
		return Solution{}, fmt.Errorf("error: unbalanced network, sum of b = %.4f (must be 0)", sum)
	}

	spanningTree := g.BuildSpanningTree(root)
	flows := spanningTree.ComputeParticularSolution(vectorB, len(g.Arcs))

	return Solution{Graph: g, Flows: flows}, nil
}

// ComputeParticularSolution calculates the particular solution
func (st *SpanningTree) ComputeParticularSolution(vectorB []float64, edgeCount int) []float64 {
	particularSolution := make([]float64, edgeCount)
	currentBalance := make([]float64, len(vectorB))

	// Bottom-up traversal from leaves to root based on the post-visit traversal order
	for _, nodeID := range st.TraversalOrder {
		if nodeID == st.Root {
			continue
		}

		parentID := st.Predecessors[nodeID]
		arc := st.ArcToPred[nodeID]
		edgeID := arc.ID

		if arc.From == nodeID {
			// Edge direction: nodeID -> parentID (leaves nodeID, enters parentID)
			particularSolution[edgeID] = vectorB[nodeID] - currentBalance[nodeID]
			currentBalance[nodeID] += particularSolution[edgeID]
			currentBalance[parentID] -= particularSolution[edgeID]
		} else {
			// Edge direction: parentID -> nodeID (enters nodeID, leaves parentID)
			particularSolution[edgeID] = currentBalance[nodeID] - vectorB[nodeID]
			currentBalance[nodeID] -= particularSolution[edgeID]
			currentBalance[parentID] += particularSolution[edgeID]
		}
	}

	return particularSolution
}

// Determining finds the fundamental cycle for a non-tree edge (chord) between nodes (nodeK, nodeL).
// Returns a map where key is tree edge ID and value is the cycle coefficient (+1 or -1)
func (st *SpanningTree) Determining(nodeK, nodeL int) map[int]int {
	deltaCycle := make(map[int]int)
	i, j := nodeK, nodeL

	// Traverse up the spanning tree from nodes k and l to their lowest common ancestor
	for i != j {
		if st.Depth[i] > st.Depth[j] {
			arc := st.ArcToPred[i]
			// Check direction of the edge in the graph relative to the parent node
			if arc.From == st.Predecessors[i] {
				deltaCycle[arc.ID] = 1
			} else {
				deltaCycle[arc.ID] = -1
			}
			i = st.Predecessors[i]
		} else if st.Depth[j] > st.Depth[i] {
			arc := st.ArcToPred[j]
			// Signs are inverted because traversal goes up the opposite path of the cycle (from node l)
			if arc.From == st.Predecessors[j] {
				deltaCycle[arc.ID] = -1
			} else {
				deltaCycle[arc.ID] = 1
			}
			j = st.Predecessors[j]
		} else {
			// If depths are equal, lift both vertices simultaneously
			arcI := st.ArcToPred[i]
			if arcI.From == st.Predecessors[i] {
				deltaCycle[arcI.ID] = 1
			} else {
				deltaCycle[arcI.ID] = -1
			}

			arcJ := st.ArcToPred[j]
			if arcJ.From == st.Predecessors[j] {
				deltaCycle[arcJ.ID] = -1
			} else {
				deltaCycle[arcJ.ID] = 1
			}

			i = st.Predecessors[i]
			j = st.Predecessors[j]
		}
	}
	return deltaCycle
}

// ComputeGeneralSolutionMatrix builds the fundamental system of solutions.
// Each row corresponds to a single non-tree edge (chord) and defines its fundamental cycle vector
func (st *SpanningTree) ComputeGeneralSolutionMatrix(edgeCount int) [][]float64 {
	matrix := make([][]float64, len(st.Chords))

	for index, chord := range st.Chords {
		row := make([]float64, edgeCount)

		// The chord itself always takes a coefficient of 1.0 in its fundamental cycle
		row[chord.ID] = 1.0

		// Evaluate coefficients for the tree edges included in this cycle
		treeEdgesInCycle := st.Determining(chord.From, chord.To)
		for arcID, coeff := range treeEdgesInCycle {
			row[arcID] = float64(coeff)
		}

		matrix[index] = row
	}

	return matrix
}

// VerifySolution checks whether the flow vector satisfies the node balance equations (Ax = b).
func VerifySolution(g *Graph, solutionX []float64, vectorB []float64) (bool, string) {
	calculatedB := make([]float64, len(vectorB))

	for _, arc := range g.Arcs {
		if arc.From >= len(calculatedB) || arc.To >= len(calculatedB) {
			return false, fmt.Sprintf("error: node index in edge (%d->%d) is out of bounds for balance vector b", arc.From, arc.To)
		}
		calculatedB[arc.From] += solutionX[arc.ID] // Outflow (+)
		calculatedB[arc.To] -= solutionX[arc.ID]   // Inflow (-)
	}

	var resultInfo strings.Builder
	success := true
	for i := 1; i < len(vectorB); i++ {
		if i > g.VertexCount {
			continue
		}
		if math.Abs(calculatedB[i]-vectorB[i]) > 1e-9 {
			success = false
			resultInfo.WriteString(fmt.Sprintf("\n  [Node %d] Expected: %.2f, Got: %.2f (Diff: %.2f)", i, vectorB[i], calculatedB[i], calculatedB[i]-vectorB[i]))
		}
	}
	return success, resultInfo.String()
}
