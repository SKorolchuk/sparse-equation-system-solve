import React, { memo, useCallback } from "react";
import type { Arc } from "../../types";

interface ArcRowProps {
  index: number;
  arc: Arc;
  vertexCount: number;
  onChange: (index: number, field: "from" | "to", value: number) => void;
}

export const ArcRow = memo(
  ({ index, arc, vertexCount, onChange }: ArcRowProps) => {
    const handleFromChange = useCallback(
      (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = parseInt(e.target.value, 10) || 1;
        onChange(index, "from", value);
      },
      [index, onChange],
    );

    const handleToChange = useCallback(
      (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = parseInt(e.target.value, 10) || 1;
        onChange(index, "to", value);
      },
      [index, onChange],
    );

    return (
      <div className="arc-input-row">
        <span className="arc-row-label">e_{arc.id}:</span>
        <input
          type="number"
          min={1}
          max={vertexCount}
          value={arc.from}
          onChange={handleFromChange}
          className="arc-node-field"
        />
        <span className="arc-arrow">&rarr;</span>
        <input
          type="number"
          min={1}
          max={vertexCount}
          value={arc.to}
          onChange={handleToChange}
          className="arc-node-field"
        />
      </div>
    );
  },
);
ArcRow.displayName = "ArcRow";
