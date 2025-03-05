import {
  SyntheticEvent,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react';

import { Mas, StorageCost } from '@massalabs/massa-web3';
import {
  AccordionCategory,
  AccordionContent,
  InlineMoney,
  Tooltip,
  formatAmount,
  massaToken,
} from '@massalabs/react-ui-kit';
import { LogPrint } from '@wailsjs/runtime/runtime';
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
  fees: Mas.Mas;
  minFees: Mas.Mas;
  setFees: (fees: Mas.Mas) => void;
  feesError?: string;
  isEditing?: boolean;
  setIsEditing: (isEditing: boolean) => void;
  allowFeeEdition: boolean;
  DeployedByteCodeSize: number; // for executeSC of type deploySC
  DeployedCoins: string; // for executeSC of type deploySC, in nanoMAS
}

export function OperationCost(props: OperationCostProps) {
  const hideCoins = props.coins ? props.coins === '0' : true;
  const hideByteCodeCost = props.DeployedByteCodeSize === 0;
  const hideDeployedCoins = props.DeployedCoins
    ? props.DeployedCoins === '0'
    : true;

  LogPrint(
    '\nOperationCost props DeployedByteCodeSize:' + props.DeployedByteCodeSize,
  );
  LogPrint('\nOperationCost props DeployedCoins:' + props.DeployedCoins);
  LogPrint('\nOperationCost hideByteCodeCost:' + hideByteCodeCost);
  LogPrint('\nOperationCost hideDeployedCoins:' + hideDeployedCoins);

  const coins = BigInt(props.coins ?? 0);
  const byteCodeStorageCost = props.DeployedByteCodeSize
    ? StorageCost.smartContract(props.DeployedByteCodeSize)
    : 0n;

  const deployedCoins = BigInt(props.DeployedCoins ?? 0);

  const { fees, setFees, minFees, isEditing, setIsEditing, allowFeeEdition } =
    props;

  const [error, setError] = useState<string>();

  /* Compute fees*/
  function getDefaultFees(): Mas.Mas {
    if (fees === 0n) {
      return minFees;
    }

    if (fees < minFees) {
      return minFees;
    }

    return fees;
  }

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const defaultFees = useMemo(getDefaultFees, [minFees]);

  if (fees === 0n) setFees(defaultFees);
  if (fees < minFees) setFees(minFees);

  /* Handle operation cost*/
  const computeOperationCost = useCallback(() => {
    return coins + fees + deployedCoins + byteCodeStorageCost;
  }, [fees, coins, deployedCoins, byteCodeStorageCost]);

  const [operationCost, setOperationCost] = useState(computeOperationCost());

  useEffect(() => {
    setOperationCost(computeOperationCost());
  }, [computeOperationCost]);

  /* handleConfirmTemplate callback functions*/
  function handleEdit(e: SyntheticEvent) {
    e.preventDefault();
    setIsEditing(true);
  }

  function validate(): boolean {
    if (fees < minFees) {
      setError(Intl.t('password-prompt.sign.fees-to-low'));
      return false;
    }

    if (fees >= MAX_FEES) {
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

  function handleFeesChange(event: NumberFormatValues) {
    setFees(Mas.fromString(event.value));
  }

  /* Set fee edition button if fee edition is allowed */
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
                  value={Mas.toString(fees)}
                  onValueChange={handleFeesChange}
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
              {!hideDeployedCoins && (
                <div className="flex justify-between pb-2">
                  <p className="flex mas-caption">
                    <Tooltip
                      className="mr-1"
                      body={Intl.t(
                        'password-prompt.sign.deployed-coins-tooltip',
                      )}
                    >
                      <FiInfo size={16} />
                    </Tooltip>
                    {Intl.t('password-prompt.sign.deployed-coins')}
                  </p>
                  <InlineMoney
                    customClass="mas-caption"
                    disabled
                    value={formatAmount(deployedCoins).full}
                  />
                </div>
              )}
              {!hideByteCodeCost && (
                <div className="flex justify-between pb-2">
                  <p className="flex mas-caption">
                    <Tooltip
                      className="mr-1"
                      body={Intl.t(
                        'password-prompt.sign.deployed-bytecode-cost-tooltip',
                      )}
                    >
                      <FiInfo size={16} />
                    </Tooltip>
                    {Intl.t('password-prompt.sign.deployed-bytecode-cost')}
                  </p>
                  <InlineMoney
                    customClass="mas-caption"
                    disabled
                    value={formatAmount(byteCodeStorageCost).full}
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
