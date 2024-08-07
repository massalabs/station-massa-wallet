import { SyntheticEvent, useEffect, useMemo, useState } from 'react';

import { fromMAS, toMAS } from '@massalabs/massa-web3';
import {
  AccordionCategory,
  AccordionContent,
  InlineMoney,
  Tooltip,
} from '@massalabs/react-ui-kit';
import { formatAmount } from '@massalabs/react-ui-kit';
import { massaToken } from '@massalabs/react-ui-kit/src/lib/massa-react/const';
import BigNumber from 'bignumber.js';
import {
  FiCheck,
  FiChevronDown,
  FiChevronUp,
  FiEdit2,
  FiInfo,
  FiRotateCcw,
} from 'react-icons/fi';
import { NumberFormatValues } from 'react-number-format';

import Intl from '@/i18n/i18n';

export const MAX_FEES = 1_000_000_000_000_000_000n; // MASSA total supply

export interface OperationCostProps {
  coins?: string; // in nanoMAS
  fees: string; // in MAS
  minFees: string; // in MAS
  setFees: (fees: string) => void;
  feesError?: string;
  isEditing?: boolean;
  setIsEditing: (isEditing: boolean) => void;
  allowFeeEdition: boolean;
}

export function OperationCost(props: OperationCostProps) {
  const hideCoins = props.coins === undefined;

  const coins = toMAS(props.coins || 0).toFixed(9);
  const { fees, setFees, minFees, isEditing, setIsEditing, allowFeeEdition } =
    props;

  const [operationCost, setOperationCost] = useState(
    new BigNumber(coins).plus(new BigNumber(fees)).toFixed(9),
  );
  const [error, setError] = useState<string>();

  function getDefaultFees(): string {
    if (fees === '') {
      return minFees;
    }

    if (fromMAS(fees) < fromMAS(minFees)) {
      return minFees;
    }

    return fees;
  }

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const defaultFees = useMemo(getDefaultFees, [minFees]);

  if (fees === '') setFees(defaultFees);
  if (fromMAS(fees) < fromMAS(minFees)) setFees(minFees);

  useEffect(() => {
    setOperationCost(new BigNumber(coins).plus(new BigNumber(fees)).toFixed(9));
  }, [fees, coins]);

  function handleEdit(e: SyntheticEvent) {
    e.preventDefault();
    setIsEditing(true);
  }

  function validate(): boolean {
    if (fromMAS(fees) < fromMAS(minFees)) {
      setError(Intl.t('password-prompt.sign.fees-to-low'));
      return false;
    }

    if (fromMAS(fees) >= MAX_FEES) {
      setError(Intl.t('password-prompt.sign.fees-to-high'));
      return false;
    }

    setError('');
    return true;
  }

  function handleConfirm(e: SyntheticEvent) {
    e.preventDefault();

    if (!validate()) return;

    setIsEditing(false);
  }

  function handleReset(e: SyntheticEvent) {
    e.preventDefault();
    setIsEditing(false);
    setError('');
    setFees(defaultFees);
  }

  const feeEditionButtonsRow = allowFeeEdition ? (
    <div className="flex justify-end gap-1">
      {isEditing ? (
        <>
          <button className="flex hover:cursor-pointer" onClick={handleConfirm}>
            <FiCheck size={16} className="mr-1" />
            {Intl.t('password-prompt.sign.confirm-fees')}
          </button>
          <p className="px-1 hover:cursor-default">|</p>
          <button className="flex hover:cursor-pointer" onClick={handleReset}>
            <FiRotateCcw size={16} className="mr-1" />
            {Intl.t('password-prompt.sign.reset-fees')}
          </button>
        </>
      ) : (
        <button className="flex hover:cursor-pointer" onClick={handleEdit}>
          <FiEdit2 size={16} className="mr-1" />
          {Intl.t('password-prompt.sign.edit-fees')}
        </button>
      )}
    </div>
  ) : null;

  return (
    <div className="w-full">
      <div className="flex w-full justify-between pb-2">
        <p>{Intl.t('password-prompt.sign.operation-const')}</p>
        <p>
          {formatAmount(operationCost).full} {massaToken}
        </p>
      </div>

      <div className="flex flex-col w-full h-fit">
        <AccordionCategory
          isChild={false}
          iconOpen={<FiChevronDown />}
          iconClose={<FiChevronUp />}
          customClass="px-0 py-0"
          categoryTitle={
            <p className="p-0 pb-2 mas-caption text-f-disabled-1">
              {Intl.t('password-prompt.sign.view-details')}
            </p>
          }
        >
          <AccordionContent customClass="px-0 py-0">
            <div className="flex flex-col gap-1 text-f-disabled-1">
              <div className="flex justify-between">
                <p className="flex mas-caption">
                  <Tooltip
                    className="mr-1"
                    body={
                      <>
                        {Intl.t('password-prompt.sign.fees-tooltip.1')}
                        <br />
                        {Intl.t('password-prompt.sign.fees-tooltip.2')}
                      </>
                    }
                  >
                    <FiInfo size={16} />
                  </Tooltip>
                  <label className="mas-caption">
                    {Intl.t('password-prompt.sign.fees')}
                  </label>
                </p>
                <InlineMoney
                  customClass="mas-caption"
                  disabled={!isEditing}
                  value={fees}
                  onValueChange={(event: NumberFormatValues) =>
                    setFees(event.value)
                  }
                />
              </div>
              {error && (
                <div className="flex justify-between text-s-error">{error}</div>
              )}
              {!hideCoins && (
                <div className="flex justify-between pb-2">
                  <p className="flex mas-caption">
                    <Tooltip
                      className="mr-1"
                      body={
                        <>
                          {Intl.t('password-prompt.sign.coins-tooltip')}
                          <br />
                          {Intl.t('password-prompt.sign.coins-tooltip-2')}
                        </>
                      }
                    >
                      <FiInfo size={16} />
                    </Tooltip>
                    {Intl.t('password-prompt.sign.coins')}
                  </p>
                  <InlineMoney
                    customClass="mas-caption"
                    disabled
                    value={formatAmount(coins).full}
                  />
                </div>
              )}
              {feeEditionButtonsRow}
            </div>
          </AccordionContent>
        </AccordionCategory>
      </div>
    </div>
  );
}
