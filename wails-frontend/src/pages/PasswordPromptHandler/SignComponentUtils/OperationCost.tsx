import { SyntheticEvent, useEffect, useState } from 'react';

import { toMAS } from '@massalabs/massa-web3';
import {
  AccordionCategory,
  AccordionContent,
  InlineMoney,
  Tooltip,
} from '@massalabs/react-ui-kit';
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
import { formatStandard, masToken } from '@/utils';

export interface OperationCostProps {
  coins?: string; // in nanoMAS
  fees: string; // in MAS
  defaultFees: string; // in MAS
  setFees: (fees: string) => void;
  isEditing?: boolean;
  setIsEditing: (isEditing: boolean) => void;
  allowFeeEdition: boolean;
}

export function OperationCost(props: OperationCostProps) {
  const hideCoins = props.coins === undefined;

  const coins = toMAS(props.coins || 0).toFixed(9);
  const {
    fees,
    setFees,
    defaultFees,
    isEditing,
    setIsEditing,
    allowFeeEdition,
  } = props;

  const [operationCost, setOperationCost] = useState(computeCost());

  if (fees === '') setFees(defaultFees);

  useEffect(() => {
    setOperationCost(computeCost());
  }, [fees, coins]);

  function handleEdit(e: SyntheticEvent) {
    e.preventDefault();
    setIsEditing(true);
  }

  function handleConfirm(e: SyntheticEvent) {
    e.preventDefault();
    setIsEditing(false);
  }

  function handleReset(e: SyntheticEvent) {
    e.preventDefault();
    setIsEditing(false);
    setFees(defaultFees);
  }

  function computeCost() {
    return new BigNumber(coins).plus(new BigNumber(fees)).toFixed(9);
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
          {formatStandard(operationCost)} {masToken}
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
                    icon={<FiInfo size={16} />}
                    content={
                      <>
                        {Intl.t('password-prompt.sign.fees-tooltip.1')}
                        <br />
                        {Intl.t('password-prompt.sign.fees-tooltip.2')}
                      </>
                    }
                  />
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
              {!hideCoins && (
                <div className="flex justify-between pb-2">
                  <p className="flex mas-caption">
                    <Tooltip
                      className="mr-1"
                      icon={<FiInfo size={16} />}
                      content={
                        <>
                          {Intl.t('password-prompt.sign.coins-tooltip')}
                          <br />
                          {Intl.t('password-prompt.sign.coins-tooltip-2')}
                        </>
                      }
                    />
                    {Intl.t('password-prompt.sign.coins')}
                  </p>
                  <InlineMoney
                    customClass="mas-caption"
                    disabled
                    value={formatStandard(coins)}
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
