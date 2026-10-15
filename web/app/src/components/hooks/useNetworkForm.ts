import { useCallback, useMemo, useState } from "react";
import type { Arc } from "../../types";

interface UseNetworkFormProps {
  onCompute: (
    vertexCount: number,
    arcs: Arc[],
    balances: number[],
    root: number,
  ) => void;
}

interface NetworkFormResult {
  vertexCount: number;
  edgeCount: number;
  rootNode: number;
  balances: string[];
  arcs: Arc[];
  localError: string;
  balanceResidual: number;
  isFormInvalid: boolean;
  setRootNode: React.Dispatch<React.SetStateAction<number>>;
  handleVertexCountChange: (value: number) => void;
  handleEdgeCountChange: (value: number) => void;
  handleBalanceChange: (index: number, value: string) => void;
  handleArcChange: (index: number, field: "from" | "to", value: number) => void;
  validateAndSubmit: (e: React.FormEvent) => void;
}

export function useNetworkForm({
  onCompute,
}: UseNetworkFormProps): NetworkFormResult {
  const [vertexCount, setVertexCount] = useState<number>(9);
  const [edgeCount, setEdgeCount] = useState<number>(17);
  const [rootNode, setRootNode] = useState<number>(5);
  const [balances, setBalances] = useState<string[]>([
    "4",
    "-6",
    "-10",
    "2",
    "10",
    "9",
    "-1",
    "-3",
    "-5",
  ]);
  const [arcs, setArcs] = useState<Arc[]>([
    { id: 0, from: 1, to: 2 },
    { id: 1, from: 1, to: 3 },
    { id: 2, from: 2, to: 4 },
    { id: 3, from: 2, to: 5 },
    { id: 4, from: 3, to: 6 },
    { id: 5, from: 3, to: 7 },
    { id: 6, from: 5, to: 8 },
    { id: 7, from: 7, to: 9 },
    { id: 8, from: 4, to: 5 },
    { id: 9, from: 6, to: 7 },
    { id: 10, from: 8, to: 1 },
    { id: 11, from: 9, to: 2 },
    { id: 12, from: 4, to: 1 },
    { id: 13, from: 5, to: 6 },
    { id: 14, from: 8, to: 4 },
    { id: 15, from: 9, to: 6 },
    { id: 16, from: 3, to: 2 },
  ]);

  const [localError, setLocalError] = useState<string>("");

  const balanceResidual = useMemo(() => {
    return balances.reduce(
      (acc, current) => acc + (parseFloat(current) || 0.0),
      0.0,
    );
  }, [balances]);

  const handleVertexCountChange = useCallback((value: number): void => {
    const nextCount = Math.max(2, value);
    setVertexCount(nextCount);

    setBalances((prev) => {
      const next = [...prev];
      if (next.length < nextCount) {
        while (next.length < nextCount) next.push("0");
      } else {
        next.length = nextCount;
      }
      return next;
    });
  }, []);

  const handleEdgeCountChange = useCallback((value: number): void => {
    const nextCount = Math.max(1, value);
    setEdgeCount(nextCount);

    setArcs((prev) => {
      const next = [...prev];
      if (next.length < nextCount) {
        while (next.length < nextCount) {
          next.push({ id: next.length, from: 1, to: 2 });
        }
      } else {
        next.length = nextCount;
      }
      return next;
    });
  }, []);

  const handleBalanceChange = useCallback(
    (index: number, value: string): void => {
      setBalances((prev) => {
        const updated = [...prev];
        updated[index] = value;
        return updated;
      });
    },
    [],
  );

  const handleArcChange = useCallback(
    (index: number, field: "from" | "to", value: number): void => {
      setArcs((prev) => {
        const updated = [...prev];
        updated[index] = { ...updated[index], [field]: value };
        return updated;
      });
    },
    [],
  );

  const validateAndSubmit = useCallback(
    (e: React.FormEvent): void => {
      e.preventDefault();
      setLocalError("");

      if (Math.abs(balanceResidual) > 1e-9) {
        setLocalError(
          `Validation error: vector b elements must sum up to exactly 0: ${balanceResidual.toFixed(4)}`,
        );
        return;
      }

      for (const arc of arcs) {
        if (
          arc.from < 1 ||
          arc.from > vertexCount ||
          arc.to < 1 ||
          arc.to > vertexCount
        ) {
          setLocalError(
            `Topology error at Edge ${arc.id}: nodes must be within range [1 .. ${vertexCount}]`,
          );
          return;
        }
        if (arc.from === arc.to) {
          setLocalError(
            `Topology error at Edge ${arc.id}: graph loops (from node ${arc.from} to itself) are prohibited`,
          );
          return;
        }
      }

      if (rootNode < 1 || rootNode > vertexCount) {
        setLocalError(
          `Root error: target root index must be within range [1 .. ${vertexCount}]`,
        );
        return;
      }

      const floatBalances = [0.0, ...balances.map((v) => parseFloat(v) || 0.0)];
      onCompute(vertexCount, arcs, floatBalances, rootNode);
    },
    [arcs, balanceResidual, balances, onCompute, rootNode, vertexCount],
  );

  const isFormInvalid = Math.abs(balanceResidual) > 1e-9;

  return {
    vertexCount,
    edgeCount,
    rootNode,
    balances,
    arcs,
    localError,
    balanceResidual,
    isFormInvalid,
    setRootNode,
    handleVertexCountChange,
    handleEdgeCountChange,
    handleBalanceChange,
    handleArcChange,
    validateAndSubmit,
  };
}
