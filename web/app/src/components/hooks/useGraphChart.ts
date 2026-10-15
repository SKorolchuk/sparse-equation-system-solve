import { useEffect, useRef, useMemo } from 'react';
import * as echarts from 'echarts';
import type { Arc } from '../../types';

// internal helper to construct ECharts options from raw graph primitives
const buildChartOption = (
  title: string,
  vertexCount: number,
  edgeList: Arc[],
  bVector: number[],
  treeEdgeIds: number[],
  rootNode: number
): echarts.EChartsOption => {
  const nodes = [];
  for (let i = 1; i <= vertexCount; i++) {
    const intensity = bVector[i] !== undefined ? bVector[i] : 0;
    const isRoot = i === rootNode;
    
    nodes.push({
      id: i.toString(),
      name: `Node ${i}\nb_${i}=${intensity.toFixed(1)}`,
      symbolSize: isRoot ? 50 : 38,
      itemStyle: {
        color: isRoot ? '#ffcdbe' : '#f0f0f0',
        borderColor: isRoot ? '#474746' : '#666666',
        borderWidth: isRoot ? 2 : 1,
      },
    });
  }

  const hasTreeConstraints = treeEdgeIds.length > 0;
  const links = edgeList.map((arc) => {
    const isTreeEdge = treeEdgeIds.includes(arc.id);
    const lineType: 'solid' | 'dashed' = isTreeEdge ? 'solid' : 'dashed';

    return {
      source: arc.from.toString(),
      target: arc.to.toString(),
      label: { 
        show: true, 
        formatter: `e${arc.id}`,
        fontSize: 10,
        fontFamily: 'monospace',
        backgroundColor: '#ffffff',
        padding: 10,
        borderRadius: 2
      },
      lineStyle: {
        color: isTreeEdge ? '#38e246' : (hasTreeConstraints ? '#cccccc' : '#444444'),
        width: isTreeEdge ? 3 : 1.5,
        type: hasTreeConstraints ? lineType : 'solid',
        curveness: 0.1,
      },
    };
  });

  const explicitArrowSizes: [number, number] = Array.of(0, 10) as [number, number];

  return {
    title: { 
      text: title, 
      textStyle: { fontSize: 13, fontFamily: 'monospace', fontWeight: 'bold' }, 
      left: 'center' 
    },
    tooltip: {
      trigger: 'item',
      formatter: (params: echarts.TooltipComponentOption & { dataType?: string }) => {
        if (params && params.dataType === 'node' && typeof params.name === 'string') {
          return params.name.replace('\n', '<br/>');
        }
        return params.name || '';
      }
    },
    series: [{
      type: 'graph',
      layout: 'force',
      force: { 
        repulsion: 600,
        edgeLength: 120,
        gravity: 0.02,
        friction: 0.4,
        layoutAnimation: true 
      },
      roam: true,
      draggable: true,
      label: { 
        show: true, 
        position: 'inside',
        fontSize: 10, 
        fontFamily: 'monospace',
        color: '#111111'
      },
      edgeSymbol: ['none', 'arrow'],
      edgeSymbolSize: explicitArrowSizes, 
      data: nodes,
      links: links,
    }]
  } as echarts.EChartsOption;
};

interface UseGraphChartProps {
  title: string;
  vertexCount: number;
  arcs: Arc[];
  balances: number[];
  treeEdgeIds: number[];
  rootNode: number;
}

export const useGraphChart = ({
  title,
  vertexCount,
  arcs,
  balances,
  treeEdgeIds,
  rootNode,
}: UseGraphChartProps): React.RefObject<HTMLDivElement | null> => {
  const containerRef = useRef<HTMLDivElement>(null);
  const chartInstance = useRef<echarts.ECharts | null>(null);

  const option = useMemo(() => {
    if (vertexCount === 0 || arcs.length === 0) return null;
    return buildChartOption(title, vertexCount, arcs, balances, treeEdgeIds, rootNode);
  }, [title, vertexCount, arcs, balances, treeEdgeIds, rootNode]);

  useEffect(() => {
    if (!containerRef.current || !option) return;

    if (!chartInstance.current) {
      chartInstance.current = echarts.init(containerRef.current);
    }

    chartInstance.current.setOption(option);
  }, [option]);

  useEffect(() => {
    const element = containerRef.current;
    if (!element) return;

    const observer = new ResizeObserver(() => {
      chartInstance.current?.resize();
    });
    observer.observe(element);

    return () => {
      observer.disconnect();
      chartInstance.current?.dispose();
      chartInstance.current = null;
    };
  }, []);

  return containerRef;
}
