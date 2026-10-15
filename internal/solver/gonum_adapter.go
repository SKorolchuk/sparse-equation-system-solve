package networksolver

import (
	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/mat"
)

// GonumAdapter links Gonum graph structures with the internal library structures
type GonumAdapter struct {
	gonumGraph    graph.Graph
	internalGraph *Graph
	arcs          []Arc
}

// NewGonumAdapter constructs and initializes a new GonumAdapter.
func NewGonumAdapter(g graph.Graph) *GonumAdapter {
	adapter := &GonumAdapter{
		gonumGraph: g,
		arcs:       make([]Arc, 0),
	}
	adapter.convertGraph()
	return adapter
}

// convertGraph translates Gonum topology into internal structures.
func (a *GonumAdapter) convertGraph() {
	var currentArcID int
	var maxNodeID int

	// 1. Extract all graph nodes as a flat slice
	allNodes := graph.NodesOf(a.gonumGraph.Nodes())

	// 2. Traverse each node to build adjacency list representation
	for _, fromNode := range allNodes {
		fromID := int(fromNode.ID())
		if fromID > maxNodeID {
			maxNodeID = fromID
		}

		// Retrieve all target nodes reachable from the current node
		targetNodes := graph.NodesOf(a.gonumGraph.From(fromNode.ID()))

		for _, toNode := range targetNodes {
			toID := int(toNode.ID())
			if toID > maxNodeID {
				maxNodeID = toID
			}

			arc := Arc{
				ID:   currentArcID,
				From: fromID,
				To:   toID,
			}
			a.arcs = append(a.arcs, arc)
			currentArcID++
		}
	}

	// Boundary guard: handle potential empty graph edge cases
	if maxNodeID == 0 {
		maxNodeID = 1
	}

	a.internalGraph = NewGraph(maxNodeID, a.arcs)
}

// SolveParticular accepts a intensity vector and computes a particular solution vector
func (a *GonumAdapter) SolveParticular(vectorB mat.Vector, rootID int64) (*mat.VecDense, error) {
	vertexCount := a.internalGraph.VertexCount
	internalB := make([]float64, vertexCount+1)

	for i := 1; i <= vertexCount; i++ {
		internalB[i] = vectorB.At(i, 0)
	}

	solution, err := a.internalGraph.SolveParticular(internalB, int(rootID))
	if err != nil {
		return nil, err
	}

	return mat.NewVecDense(len(solution.Flows), solution.Flows), nil
}

// ComputeGeneralSolutionMatrix constructs the fundamental system of solutions matrix
func (a *GonumAdapter) ComputeGeneralSolutionMatrix(rootID int64) [][]float64 {
	spanningTree := a.internalGraph.BuildSpanningTree(int(rootID))
	return spanningTree.ComputeGeneralSolutionMatrix(len(a.arcs))
}
