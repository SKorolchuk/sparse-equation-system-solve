import React, { useEffect, useRef } from 'react';
import katex from 'katex';

interface SolutionSectionProps {
  title: string;
  texFormula: string;
  scrollable?: boolean;
}

export const SolutionSection = React.memo(({ 
  title, 
  texFormula, 
  scrollable = false 
}: SolutionSectionProps): React.JSX.Element => {
  const formulaRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (formulaRef.current && texFormula) {
      katex.render(texFormula, formulaRef.current, {
        displayMode: true,
        throwOnError: false
      });
    }
  }, [texFormula]);

  if (!texFormula) return <React.Fragment />;

  const viewportClassName = `latex-formula-viewport ${scrollable ? 'scrollable-matrix' : ''}`;

  return (
    <div>
      <span className="solution-section-title">{title}</span>
      <div ref={formulaRef} className={viewportClassName} />
    </div>
  );
});

SolutionSection.displayName = 'SolutionSection';
