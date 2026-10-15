import type { Arc, SolverResponse } from './types';

declare global {
  interface Window {
    Go: new () => { 
      run: (instance: WebAssembly.Instance) => Promise<void>; 
      importObject: WebAssembly.Imports 
    };
    solveNetworkSystem: (
      vertexCount: number, 
      arcs: Arc[], 
      b: number[], 
      root: number
    ) => SolverResponse;
  }
}

export {};
