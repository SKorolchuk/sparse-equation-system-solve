import React from 'react';
import type { Arc, SolverResponse } from '../types';
import { SourceGraphCanvas } from './graph-plot/SourceGraphCanvas';
import { SpanningTreeCanvas } from './graph-plot/SpanningTreeCanvas';

import './GraphVisualizer.css';

interface GraphVisualizerProps {
  vertexCount: number;
  arcs: Arc[];
  balances: number[];
  output: SolverResponse | null;
}

export const GraphVisualizer = ({ vertexCount, arcs, balances, output }: GraphVisualizerProps): React.JSX.Element => {
  const activeRoot = output ? output.rootNode : 1;
  const treeEdgeIds = output ? output.treeEdgeIds : [];

  return (
    <div className="graph-visualizer-container">
      <SourceGraphCanvas 
        vertexCount={vertexCount} 
        arcs={arcs} 
        balances={balances} 
        rootNode={activeRoot} 
      />

      <SpanningTreeCanvas 
        vertexCount={vertexCount} 
        arcs={arcs} 
        balances={balances} 
        treeEdgeIds={treeEdgeIds} 
        rootNode={activeRoot} 
      />
    </div>
  );
}
GraphVisualizer.displayName = 'GraphVisualizer';