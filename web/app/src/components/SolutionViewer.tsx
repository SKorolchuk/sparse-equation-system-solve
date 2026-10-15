import React from 'react';
import type { Arc, SolverResponse } from '../types';
import { SolutionSection } from './SolutionSection';
import { useSolutionCompiler } from './hooks/useSolutionCompiler';

import 'katex/dist/katex.min.css';
import './SolutionViewer.css';

interface SolutionViewerProps {
  output: SolverResponse | null;
  arcs: Arc[];
  vertexCount: number;
}

export const SolutionViewer = ({ output, arcs, vertexCount }: SolutionViewerProps): React.JSX.Element => {
  const {
    particularSolutionTex,
    incidenceSystemTex,
    vectorGeneralSolutionTex,
    verificationSystemTex,
  } = useSolutionCompiler({ output, arcs, vertexCount });

  if (!output) {
    return (
      <div className="solution-viewer-placeholder">
        Awaiting submit of input parameters...
      </div>
    );
  }

  return (
    <div className="solution-viewer-container">
      <div className="solution-viewer-header">
        SOLUTION CALCULATION REPORT:
      </div>

      <SolutionSection 
        title="1. Particular solution vector of the tree subsystem:" 
        texFormula={particularSolutionTex} 
      />

      <SolutionSection 
        title="2. Node-arc incidence constraints of the network system (Ax = b):" 
        texFormula={incidenceSystemTex} 
      />

      <SolutionSection 
        title="3. General solution:" 
        texFormula={vectorGeneralSolutionTex} 
        scrollable 
      />

      <SolutionSection 
        title="4. Numerical equivalence verification:" 
        texFormula={verificationSystemTex} 
        scrollable 
      />
    </div>
  );
}
SolutionViewer.displayName = 'SolutionViewer';