package networksolver

import (
	"bytes"
	"fmt"
)

// Arc represents a directed edge in the graph
type Arc struct {
	ID   int
	From int
	To   int
}

// SpanningTree stores topological structures and network arrays for computation
type SpanningTree struct {
	Root           int
	Predecessors   []int // Parent node index for each node
	Depth          []int // Depth level of each node in the spanning tree
	ArcToPred      []Arc // Graph edge connecting each node to its predecessor
	TraversalOrder []int // Post-visit order (bottom-up traversal sequence)
	Chords         []Arc // Non-tree edges
}

// Graph represents a network topology configuration
type Graph struct {
	VertexCount int
	Arcs        []Arc
	adjacency   [][]Arc
}

// NewGraph constructs and initializes a new Graph instance
func NewGraph(vertexCount int, arcs []Arc) *Graph {
	graph := &Graph{
		VertexCount: vertexCount,
		Arcs:        arcs,
		adjacency:   make([][]Arc, vertexCount+1),
	}
	for _, arc := range arcs {
		graph.adjacency[arc.From] = append(graph.adjacency[arc.From], arc)
		graph.adjacency[arc.To] = append(graph.adjacency[arc.To], arc)
	}
	return graph
}

// BuildSpanningTree creates a SpanningTree structure utilizing a Depth-First Search algorithm
func (g *Graph) BuildSpanningTree(root int) *SpanningTree {
	tree := &SpanningTree{
		Root:           root,
		Predecessors:   make([]int, g.VertexCount+1),
		Depth:          make([]int, g.VertexCount+1),
		ArcToPred:      make([]Arc, g.VertexCount+1),
		TraversalOrder: make([]int, 0, g.VertexCount),
	}

	visited := make([]bool, g.VertexCount+1)
	inTree := make([]bool, len(g.Arcs))

	var dfs func(currentVertex, previousVertex, depth int)
	dfs = func(v, p, d int) {
		visited[v] = true
		tree.Predecessors[v] = p
		tree.Depth[v] = d

		for _, arc := range g.adjacency[v] {
			next := arc.To
			if arc.To == v {
				next = arc.From
			}

			if !visited[next] {
				inTree[arc.ID] = true
				tree.ArcToPred[next] = arc
				dfs(next, v, d+1)
			}
		}
		// POST-VISIT: Append node to traversal order sequence upon recursion exit
		tree.TraversalOrder = append(tree.TraversalOrder, v)
	}

	dfs(root, 0, 0)

	for _, arc := range g.Arcs {
		if !inTree[arc.ID] {
			tree.Chords = append(tree.Chords, arc)
		}
	}
	return tree
}

// String formats the arc details (implements fmt.Stringer)
func (a Arc) String() string {
	return fmt.Sprintf("Edge %d [%d -> %d]", a.ID, a.From, a.To)
}

// String formats the graph structure data (implements fmt.Stringer)
func (g *Graph) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Network Graph (Vertices: %d, Edges: %d):\n", g.VertexCount, len(g.Arcs))
	for _, arc := range g.Arcs {
		fmt.Fprintf(&buf, "  %s\n", arc)
	}
	return buf.String()
}

// String formats the spanning tree arrays into a matrix (implements fmt.Stringer)
func (st *SpanningTree) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Spanning Tree (Root: %d):\n", st.Root)

	buf.WriteString("  Node:\t\t")
	for i := 1; i < len(st.Predecessors); i++ {
		fmt.Fprintf(&buf, "%d\t", i)
	}
	buf.WriteString("\n  Predecessors:\t")
	for i := 1; i < len(st.Predecessors); i++ {
		fmt.Fprintf(&buf, "%d\t", st.Predecessors[i])
	}
	buf.WriteString("\n  Depth:\t")
	for i := 1; i < len(st.Depth); i++ {
		fmt.Fprintf(&buf, "%d\t", st.Depth[i])
	}
	buf.WriteString("\n  ArcToPred:\t")
	for i := 1; i < len(st.ArcToPred); i++ {
		if st.Predecessors[i] == 0 {
			buf.WriteString("-\t")
		} else {
			fmt.Fprintf(&buf, "%d\t", st.ArcToPred[i].ID)
		}
	}
	buf.WriteString("\n")

	fmt.Fprintf(&buf, "  Traversal Order (Post-visit): %v\n", st.TraversalOrder)
	buf.WriteString("  Chords (Non-tree):            ")
	if len(st.Chords) == 0 {
		buf.WriteString("none\n")
	} else {
		for _, chord := range st.Chords {
			fmt.Fprintf(&buf, "%d ", chord.ID)
		}
		buf.WriteString("\n")
	}

	return buf.String()
}
