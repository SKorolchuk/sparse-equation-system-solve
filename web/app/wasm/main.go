//go:build wasm

package main

import (
	"syscall/js"

	networksolver "sparse-equation-system-solve/internal/solver"
)

// jsSolveNetworkSystem orchestrates the bridge context between JS types and the underlying solver library.
func jsSolveNetworkSystem() js.Func {
	return js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 4 {
			return map[string]any{"error": "Missing parameters. Required sequential signature layout: (vertexCount, arcs, b, root)"}
		}

		vertexCount := args[0].Int()
		arcsData := args[1]
		bVectorData := args[2]
		rootID := args[3].Int()

		if arcsData.Length() == 0 {
			return map[string]any{"error": "Edge set (arcs topology configuration) cannot be empty"}
		}

		// Convert JS Arcs payload records into native arrays
		arcs := make([]networksolver.Arc, arcsData.Length())
		for i := 0; i < arcsData.Length(); i++ {
			arcData := arcsData.Index(i)
			arcs[i] = networksolver.Arc{
				ID:   arcData.Get("id").Int(),
				From: arcData.Get("from").Int(),
				To:   arcData.Get("to").Int(),
			}
		}

		// Standardize vector b envelope mapping bounds to handle 1-based indexes safely
		vectorB := make([]float64, vertexCount+1)
		switch bVectorData.Length() {
		case vertexCount:
			for i := 0; i < bVectorData.Length(); i++ {
				vectorB[i+1] = bVectorData.Index(i).Float()
			}
		case vertexCount + 1:
			for i := 0; i < bVectorData.Length(); i++ {
				vectorB[i] = bVectorData.Index(i).Float()
			}
		default:
			return map[string]any{"error": "Structural layout error: vector b items count size must match the total vertex count"}
		}

		graph := networksolver.NewGraph(vertexCount, arcs)

		particularSolution, err := graph.SolveParticular(vectorB, rootID)
		if err != nil {
			return map[string]any{"error": err.Error()}
		}

		spanningTree := graph.BuildSpanningTree(rootID)

		// Cache globally shared JS Array constructor interface to prevent execution stack thrashing
		jsArrayFactory := js.Global().Get("Array")

		jsTreeEdgeIDs := jsArrayFactory.New(0)
		for i := 1; i < len(spanningTree.ArcToPred); i++ {
			if spanningTree.Predecessors[i] != 0 {
				jsTreeEdgeIDs.Call("push", spanningTree.ArcToPred[i].ID)
			}
		}

		jsChordIDs := jsArrayFactory.New(0)
		for _, chord := range spanningTree.Chords {
			jsChordIDs.Call("push", chord.ID)
		}

		jsFlows := jsArrayFactory.New(len(particularSolution.Flows))
		for i, flow := range particularSolution.Flows {
			jsFlows.SetIndex(i, flow)
		}

		jsBalances := jsArrayFactory.New(len(vectorB))
		for i, balance := range vectorB {
			jsBalances.SetIndex(i, balance)
		}

		deltaMatrix := spanningTree.ComputeGeneralSolutionMatrix(len(arcs))

		jsMatrix := jsArrayFactory.New(len(deltaMatrix))
		for i, row := range deltaMatrix {
			jsRow := jsArrayFactory.New(len(row))
			for j, val := range row {
				jsRow.SetIndex(j, val)
			}
			jsMatrix.SetIndex(i, jsRow)
		}

		return map[string]any{
			"flows":       jsFlows,
			"delta":       jsMatrix,
			"treeEdgeIds": jsTreeEdgeIDs,
			"chordIds":    jsChordIDs,
			"balances":    jsBalances,
			"rootNode":    rootID,
		}
	})
}

func main() {
	keepAliveChannel := make(chan struct{})
	js.Global().Set("solveNetworkSystem", jsSolveNetworkSystem())
	<-keepAliveChannel
}
