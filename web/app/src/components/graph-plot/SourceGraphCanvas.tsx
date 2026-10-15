import React from 'react';
import type { Arc } from '../../types';
import { useGraphChart } from '../hooks/useGraphChart';

interface SourceGraphCanvasProps {
  vertexCount: number;
  arcs: Arc[];
  balances: number[];
  rootNode: number;
}

export const SourceGraphCanvas = React.memo(({
  vertexCount,
  arcs,
  balances,
  rootNode
}: SourceGraphCanvasProps): React.JSX.Element => {
  const containerRef = useGraphChart({
    title: "Graph Topology",
    vertexCount,
    arcs,
    balances,
    treeEdgeIds: [],
    rootNode,
  });

  return <div ref={containerRef} className="chart-canvas-viewport" />;
});

SourceGraphCanvas.displayName = 'SourceGraphCanvas';
