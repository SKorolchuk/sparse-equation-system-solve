import { memo } from "react";
import { BalanceItem } from "./BalanceItem";

interface BalancesInputGroupProps {
  balances: string[];
  balanceResidual: number;
  onBalanceChange: (index: number, value: string) => void;
}

export const BalancesInputGroup = memo(
  ({ balances, balanceResidual, onBalanceChange }: BalancesInputGroupProps) => {
    const isUnbalanced = Math.abs(balanceResidual) > 1e-9;

    return (
      <div className="section-block">
        <div className="section-info-text">
          {"BALANCE VECTOR — TOTAL SUM: "}
          <span
            className={isUnbalanced ? "status-unbalanced" : "status-balanced"}
          >
            {balanceResidual.toFixed(2)}
          </span>
        </div>
        <div className="balance-grid-viewport">
          {balances.map((value, index) => (
            <BalanceItem
              key={index}
              index={index}
              value={value}
              onChange={onBalanceChange}
            />
          ))}
        </div>
      </div>
    );
  },
);
BalancesInputGroup.displayName = "BalancesInputGroup";
