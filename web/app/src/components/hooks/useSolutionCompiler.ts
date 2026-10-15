import { useMemo } from 'react';
import type { Arc, SolverResponse } from '../../types';

interface UseSolutionCompilerProps {
  output: SolverResponse | null;
  arcs: Arc[];
  vertexCount: number;
}

interface UseSolutionCompilerResult {
  particularSolutionTex: string;
  incidenceSystemTex: string;
  vectorGeneralSolutionTex: string;
  verificationSystemTex: string;
}

export function useSolutionCompiler({ output, arcs, vertexCount }: UseSolutionCompilerProps): UseSolutionCompilerResult {
  
  const particularSolutionTex = useMemo(() => {
    if (!output || !output.flows) return '';
    let tex = String.raw`\tilde{x} = \begin{pmatrix} `;
    tex += output.flows.map(f => f.toFixed(2)).join(" & ");
    tex += String.raw` \end{pmatrix}^T`;
    return tex;
  }, [output]);

  const incidenceSystemTex = useMemo(() => {
    if (!output || arcs.length === 0 || vertexCount === 0 || !output.balances) return '';

    const matrixA: number[][] = Array.from({ length: vertexCount }, () => 
      Array(arcs.length).fill(0)
    );

    arcs.forEach((arc) => {
      if (arc.id >= arcs.length) return;
      matrixA[arc.from - 1][arc.id] = 1;  
      matrixA[arc.to - 1][arc.id] = -1; 
    });

    let texA = String.raw`\begin{pmatrix} `;
    matrixA.forEach((row, rowIndex) => {
      texA += row.join(" & ");
      if (rowIndex < matrixA.length - 1) texA += String.raw` \\ `;
    });
    texA += String.raw` \end{pmatrix}`;

    let texX = String.raw`\begin{pmatrix} `;
    texX += arcs.map(arc => String.raw`x_{${arc.from}\to${arc.to}}`).join(String.raw` \\ `);
    texX += String.raw` \end{pmatrix}`;

    let texB = String.raw`\begin{pmatrix} `;
    const bElements: string[] = [];
    for (let i = 1; i <= vertexCount; i++) {
      const val = output.balances[i] !== undefined ? output.balances[i] : 0;
      bElements.push(val.toFixed(2));
    }
    texB += bElements.join(String.raw` \\ `);
    texB += String.raw` \end{pmatrix}`;

    return String.raw`A \cdot x = b \implies ${texA} \cdot ${texX} = ${texB}`;
  }, [output, arcs, vertexCount]);

  const vectorGeneralSolutionTex = useMemo(() => {
    if (!output || arcs.length === 0 || !output.chordIds || !output.delta || !output.flows) return '';

    const equationsTexRows: string[] = [];
    const chordIdsSet = new Set(output.chordIds);
    
    arcs.forEach((arc) => {
      if (arc.id >= output.flows.length) return;
      if (chordIdsSet.has(arc.id)) return;
      
      const particularFlow = output.flows[arc.id];
      let chordsTexAccumulator = '';
      
      output.chordIds.forEach((chordId, chordIdx) => {
        if (chordIdx >= output.delta.length || arc.id >= output.delta[chordIdx].length) return;
        
        const coeff = output.delta[chordIdx][arc.id];
        if (coeff !== 0) {
          const sign = coeff > 0 ? "+" : "-";
          const absCoeff = Math.abs(coeff) === 1 ? "" : Math.abs(coeff).toString();
          const chordArc = arcs[chordId];
          
          if (chordArc) {
            chordsTexAccumulator += String.raw` ${sign} ${absCoeff}x^*_{${chordArc.from}\to ${chordArc.to}}`;
          }
        }
      });
      
      let rowTex = String.raw`x_{${arc.from}\to ${arc.to}} = ${particularFlow.toFixed(2)}`;
      rowTex += chordsTexAccumulator;
      equationsTexRows.push(rowTex);
    });

    if (equationsTexRows.length === 0) return '';
    return String.raw`\begin{cases} ` + equationsTexRows.join(String.raw` \\ `) + String.raw` \end{cases}`;
  }, [output, arcs]);

  // 4. Section 4
  const verificationSystemTex = useMemo(() => {
    if (!output || arcs.length === 0 || vertexCount === 0 || !output.balances || !output.flows) return '';

    const verificationRows: string[] = [];

    for (let node = 1; node <= vertexCount; node++) {
      const expectedBalance = output.balances[node] !== undefined ? output.balances[node] : 0;
      
      const terms: string[] = [];
      let calculatedSum = 0;

      arcs.forEach((arc) => {
        const flowValue = output.flows[arc.id];

        if (arc.from === node) {
          // Outflow direction (Theoretical sign: +1)
          calculatedSum += flowValue;
          
          // Format with explicit leading operational sign rules
          if (flowValue >= 0) {
            terms.push(`+ ${flowValue.toFixed(2)}`);
          } else {
            terms.push(`- ${Math.abs(flowValue).toFixed(2)}`);
          }
        } else if (arc.to === node) {
          // Inflow direction (Theoretical sign: -1)
          calculatedSum -= flowValue;
          
          // Invert formatting sign based on multiplication mechanics (-1 * flowValue)
          if (flowValue >= 0) {
            terms.push(`- ${flowValue.toFixed(2)}`);
          } else {
            terms.push(`+ ${Math.abs(flowValue).toFixed(2)}`);
          }
        }
      });

      // Joint individual strings fragments into a sequential expressions stream
      let lhsTex = terms.join(" ");
      
      // Clean up leading block artifacts (+ sign at the very beginning of the node equation)
      if (lhsTex.startsWith("+ ")) {
        lhsTex = lhsTex.substring(2);
      } else if (lhsTex.startsWith("- ")) {
        // Keeps the minus but removes the extra blank spacing for canonical TeX appearance
        lhsTex = "-" + lhsTex.substring(2);
      }
      
      if (terms.length === 0) {
        lhsTex = "0.00";
      }

      // Format final evaluation equivalence statement string
      const rowTex = String.raw`\text{Equation } ${node}: \quad ${lhsTex} = ${calculatedSum.toFixed(2)} \equiv ${expectedBalance.toFixed(2)}`;
      verificationRows.push(rowTex);
    }

    return String.raw`\begin{cases} ` + verificationRows.join(String.raw` \\ `) + String.raw` \end{cases}`;
  }, [output, arcs, vertexCount]);

  return {
    particularSolutionTex,
    incidenceSystemTex,
    vectorGeneralSolutionTex,
    verificationSystemTex,
  };
}
