import React from 'react';
import type { Arc } from '../../types';
import { useGraphChart } from '../hooks/useGraphChart';

interface SpanningTreeCanvasProps {
  vertexCount: number;
  arcs: Arc[];
  balances: number[];
  treeEdgeIds: number[];
  rootNode: number;
}

export const SpanningTreeCanvas = React.memo(({
  vertexCount,
  arcs,
  balances,
  treeEdgeIds,
  rootNode
}: SpanningTreeCanvasProps): React.JSX.Element => {
  const containerRef = useGraphChart({
    title: "Extracted Spanning Tree",
    vertexCount,
    arcs,
    balances,
    treeEdgeIds,
    rootNode,
  });

  return <div ref={containerRef} className="chart-canvas-viewport" />;
});

SpanningTreeCanvas.displayName = 'SpanningTreeCanvas';
