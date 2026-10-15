import { useState, useEffect, useCallback } from 'react';
import type { Arc, SolverResponse } from '../../types';

interface NetworkSolverResult {
  isWasmReady: boolean;
  solverOutput: SolverResponse | null;
  solverError: string;
  networkParams: {
    vertexCount: number;
    arcs: Arc[];
    balances: number[];
  };
  handleComputeEquationSystem: (vertexCount: number, arcs: Arc[], balances: number[], root: number) => void;
}

export function useNetworkSolver(): NetworkSolverResult {
  const [isWasmReady, setIsWasmReady] = useState<boolean>(false);
  const [solverOutput, setSolverOutput] = useState<SolverResponse | null>(null);
  const [solverError, setSolverError] = useState<string>('');

  const [networkParams, setNetworkParams] = useState({
    vertexCount: 0,
    arcs: [] as Arc[],
    balances: [] as number[],
  });

  // Lifecycle stream to boot the WebAssembly core instance
  useEffect(() => {
    if (isWasmReady) return;

    // window.Go is loaded globally via wasm_exec.js script block inside index.html
    const go = new window.Go();
    const wasmUrl = `/network-solver/solver.wasm`;

    WebAssembly.instantiateStreaming(fetch(wasmUrl), go.importObject)
      .then((obj) => {
        go.run(obj.instance);
        setIsWasmReady(true);
        console.info("WASM Go module was loaded without errors: /network-solver/solver.wasm");
      })
      .catch((err) => {
        console.error('Critical error: failed to boot internal WASM kernel layer:', err);
        setSolverError('Failed to initialize WebAssembly engine.');
      });
  }, [isWasmReady]);

  // Execution callback invoking the Go function exported into the global JavaScript namespace
  const computeSystem = useCallback((vertexCount: number, arcs: Arc[], balances: number[], root: number): void => {
    setSolverError('');

    try {
      if (!window.solveNetworkSystem) {
        throw new Error('WASM solver function is not registered in the global scope.');
      }

      const response = window.solveNetworkSystem(vertexCount, arcs, balances, root);

      if (response && response.error) {
        setSolverError(response.error);
        setSolverOutput(null);
        return;
      }

      setNetworkParams({ vertexCount, arcs, balances });
      setSolverOutput(response);
    } catch (err) {
      const message = err instanceof Error
                        ? err.message
                        : String(err);
      setSolverError(`Runtime context exception: ${message}`);
      setSolverOutput(null);
    }
  }, []);

  return {
    isWasmReady,
    solverOutput,
    solverError,
    networkParams,
    handleComputeEquationSystem: computeSystem,
  };
}
