import React, { memo, useCallback } from "react";

interface BalanceItemProps {
  index: number;
  value: string;
  onChange: (index: number, value: string) => void;
}

export const BalanceItem = memo(
  ({ index, value, onChange }: BalanceItemProps) => {
    const handleValueChange = useCallback(
      (e: React.ChangeEvent<HTMLInputElement>) => {
        onChange(index, e.target.value);
      },
      [index, onChange],
    );

    return (
      <div className="balance-input-card">
        <span className="balance-input-label">b_{index + 1}</span>
        <input
          type="text"
          value={value}
          onChange={handleValueChange}
          className="balance-field"
        />
      </div>
    );
  },
);
BalanceItem.displayName = "BalanceItem";
