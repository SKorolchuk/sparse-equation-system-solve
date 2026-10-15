import React from "react";
import { NetworkForm } from "./components/NetworkForm";
import { GraphVisualizer } from "./components/GraphVisualizer";
import { SolutionViewer } from "./components/SolutionViewer";
import { useNetworkSolver } from "./components/hooks/useNetworkSolver";

import "./App.css";

export default function App(): React.JSX.Element {
  const {
    isWasmReady,
    solverOutput,
    solverError,
    networkParams,
    handleComputeEquationSystem,
  } = useNetworkSolver();

  return (
    <div className="app-container">
      <h3 className="app-title">NETWORK-TYPE SPARSE SYSTEM OF LINEAR ALGEBRAIC EQUATIONS CALCULATOR (WASM ENGINE)</h3>
      <div className="main-layout">
        <div className="configuration-panel">
          <div className="panel-title">
            PARAMETERS CONFIGURATION:
          </div>
          <NetworkForm
            isWasmReady={isWasmReady}
            onCompute={handleComputeEquationSystem}
            solverError={solverError}
          />
        </div>
        <div className="output-panel">
          <div className="visualizer-wrapper">
            <GraphVisualizer
              vertexCount={networkParams.vertexCount}
              arcs={networkParams.arcs}
              balances={networkParams.balances}
              output={solverOutput}
            />
          </div>
          <div className="viewer-wrapper">
            <SolutionViewer output={solverOutput} arcs={networkParams.arcs} vertexCount={networkParams.vertexCount} />
          </div>
        </div>
      </div>
    </div>
  );
}
