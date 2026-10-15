import type { Arc } from '../types';
import { BalancesInputGroup } from './inputs/BalancesInputGroup';
import { ArcsInputGroup } from './inputs/ArcsInputGroup';
import { useNetworkForm } from './hooks/useNetworkForm';
import './NetworkForm.css';

interface NetworkFormProps {
  isWasmReady: boolean;
  onCompute: (vertexCount: number, arcs: Arc[], balances: number[], root: number) => void;
  solverError: string;
}

export const NetworkForm = ({ isWasmReady, onCompute, solverError }: NetworkFormProps): React.JSX.Element => {
  const {
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
  } = useNetworkForm({ onCompute });

  return (
    <form onSubmit={validateAndSubmit} className="network-form-container">
      <div className="form-row-header">
        <label className="number-input-field">
          Vertices Count: 
          <input 
            type="number" 
            min={2} 
            value={vertexCount} 
            onChange={e => handleVertexCountChange(parseInt(e.target.value, 10) || 2)} 
            className="number-field"
          />
        </label>
        <label className="number-input-field">
          Edges Count: 
          <input 
            type="number" 
            min={1} 
            value={edgeCount} 
            onChange={e => handleEdgeCountChange(parseInt(e.target.value, 10) || 1)} 
            className="number-field"
          />
        </label>
        <label className="number-input-field">
          Root Node ID: 
          <input 
            type="number" 
            min={1} 
            max={vertexCount} 
            value={rootNode} 
            onChange={e => setRootNode(parseInt(e.target.value, 10) || 1)} 
            className="number-field"
          />
        </label>
      </div>

      <BalancesInputGroup 
        balances={balances}
        balanceResidual={balanceResidual}
        onBalanceChange={handleBalanceChange}
      />

      <ArcsInputGroup 
        arcs={arcs}
        vertexCount={vertexCount}
        onArcChange={handleArcChange}
      />

      <div className="form-footer-actions">
        <button 
          type="submit" 
          disabled={!isWasmReady || isFormInvalid} 
          className="submit-solver-btn"
        >
          {isWasmReady ? "EXECUTE SOLVER" : "LOADING WASM ENGINE CONTEXT..."}
        </button>
        {localError && <span className="error-message-text">{localError}</span>}
        {solverError && <span className="error-message-text">{solverError}</span>}
      </div>
    </form>
  );
}
NetworkForm.displayName = "NetworkForm";