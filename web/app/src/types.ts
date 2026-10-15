export interface Arc {
  id: number;
  from: number;
  to: number;
}

export interface NetworkFormData {
  vertexCount: number;
  edgeCount: number;
  rootNode: number;
  arcs: Arc[];
  balances: number[];
}

export interface SolverResponse {
  flows: number[];
  delta: number[][];
  treeEdgeIds: number[];
  chordIds: number[];
  rootNode: number;
  balances: number[];
  error?: string;
}
