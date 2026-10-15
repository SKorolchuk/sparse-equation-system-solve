import React from "react";
import type { Arc } from "../../types";
import { ArcRow } from "./ArcRow";

interface ArcsInputGroupProps {
  arcs: Arc[];
  vertexCount: number;
  onArcChange: (index: number, field: "from" | "to", value: number) => void;
}

export const ArcsInputGroup = React.memo(
  ({
    arcs,
    vertexCount,
    onArcChange,
  }: ArcsInputGroupProps): React.JSX.Element => {
    return (
      <div className="section-block">
        <div className="section-info-text">GRAPH CONFIGURATION:</div>
        <div className="arcs-scroll-container">
          <div className="arcs-responsive-grid">
            {arcs.map((arc, index) => (
              <ArcRow
                key={arc.id}
                index={index}
                arc={arc}
                vertexCount={vertexCount}
                onChange={onArcChange}
              />
            ))}
          </div>
        </div>
      </div>
    );
  },
);
ArcsInputGroup.displayName = "ArcsInputGroup";
