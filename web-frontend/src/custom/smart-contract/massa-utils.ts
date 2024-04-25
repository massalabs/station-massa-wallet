import { Client } from '@massalabs/massa-web3';

import {
  MASSA_EXPLORER_URL,
  MASSA_EXPLO_EXTENSION,
  MASSA_EXPLO_URL,
} from '../../const/const';

export function generateExplorerLink(opId: string, isMainnet = true): string {
  const buildnetExplorerUrl = `${MASSA_EXPLO_URL}${opId}${MASSA_EXPLO_EXTENSION}`;
  const mainnetExplorerUrl = `${MASSA_EXPLORER_URL}${opId}`;
  const explorerUrl = isMainnet ? mainnetExplorerUrl : buildnetExplorerUrl;

  return explorerUrl;
}

export function logSmartContractEvents(
  client: Client,
  operationId: string,
): void {
  client
    .smartContracts()
    .getFilteredScOutputEvents({
      emitter_address: null,
      start: null,
      end: null,
      original_caller_address: null,
      original_operation_id: operationId,
      is_final: null,
    })
    .then((events) => {
      events.map((l) =>
        console.error(`opId ${operationId}: execution error ${l.data}`),
      );
    });
}
