import { Server } from 'miragejs';

import { routesForAccounts } from './account';
import { routesForConfig } from './config';
import { routesForSignRules } from './signRules';
import { AppSchema } from '../types';

const handlers = {
  accounts: routesForAccounts,
  signRules: routesForSignRules,
  config: routesForConfig,
};

const otherDomainHandlers = (server: Server) => {
  server.get('https://station.massa/massa/node', () => {
    return {
      chainId: 77658366,
      network: 'buildnet',
      url: 'https://buildnet.massa.net/api/v2',
    };
  });
  server.get('https://station.massa/plugin-manager', () => {
    return [
      {
        author: 'Massa Labs',
        description: 'Massa blockchain official wallet',
        home: '/plugin/massa-labs/massa-wallet/',
        id: '4957',
        logo: '/usr/local/share/massastation/plugins/wallet-plugin/wallet.svg',
        name: 'Massa Wallet',
        status: 'Up',
        version: '0.3.4',
      },
    ];
  });

  server.get(
    'https://station.massa/plugin/massa-labs/massa-wallet/api/accounts',
    (schema: AppSchema) => {
      const { models: accounts } = schema.all('account');

      return accounts;
    },
  );

  server.post('https://buildnet.massa.net/api/v2', (_, request) => {
    const { method } = JSON.parse(request.requestBody);
    if (method === 'get_addresses') {
      return {
        jsonrpc: '2.0',
        result: [
          {
            address: 'AU1Bknx3Du4aiGiHaeh1vo7LfwEPRF3piAwotRkdK975qCBxWwLs',
            thread: 3,
            final_balance: '2714.073821874',
            final_roll_count: 0,
            final_datastore_keys: [],
            candidate_balance: '2714.073821874',
            candidate_roll_count: 0,
            candidate_datastore_keys: [],
            deferred_credits: [],
            next_block_draws: [],
            next_endorsement_draws: [],
            created_blocks: [],
            created_operations: [],
            created_endorsements: [],
            cycle_infos: [
              {
                cycle: 4973,
                is_final: true,
                ok_count: 0,
                nok_count: 0,
                active_rolls: null,
              },
            ],
          },
        ],
        id: 0,
      };
    }

    if (method === 'execute_read_only_call') {
      return {
        jsonrpc: '2.0',
        result: [
          {
            executed_at: { period: 637344, thread: 26 },
            result: { Ok: [] },
            output_events: [
              {
                context: {
                  slot: { period: 637344, thread: 26 },
                  block: null,
                  read_only: true,
                  index_in_slot: 0,
                  call_stack: [
                    'AU1Bknx3Du4aiGiHaeh1vo7LfwEPRF3piAwotRkdK975qCBxWwLs',
                    'AS12k8viVmqPtRuXzCm6rKXjLgpQWqbuMjc37YHhB452KSUUb9FgL',
                  ],
                  origin_operation_id: null,
                  is_final: false,
                  is_error: false,
                },
                data: 'TRANSFER SUCCESS',
              },
            ],
            gas_cost: 2100000,
            state_changes: {
              ledger_changes: {
                AS12k8viVmqPtRuXzCm6rKXjLgpQWqbuMjc37YHhB452KSUUb9FgL: {
                  Update: {
                    balance: 'Keep',
                    bytecode: 'Keep',
                    datastore: [
                      [
                        [66],
                        {
                          Set: [122],
                        },
                      ],
                      [
                        [66],
                        {
                          Set: [158],
                        },
                      ],
                    ],
                  },
                },
              },
              async_pool_changes: [],
              pos_changes: {
                seed_bits: {
                  order: 'bitvec::order::Lsb0',
                  head: { width: 8, index: 0 },
                  bits: 0,
                  data: [],
                },
                roll_changes: {},
                production_stats: {},
                deferred_credits: { credits: {} },
              },
              executed_ops_changes: {},
              executed_denunciations_changes: [],
              execution_trail_hash_change: {
                Set: '2M8iBp8qZRh1yovM9MQmWCNMVkZZjkFUF7NpEioswEN5HMU7fc',
              },
            },
          },
        ],
        id: 0,
      };
    }

    if (method === 'get_operations') {
      return {
        jsonrpc: '2.0',
        result: [
          {
            id: 'O1YUEyP4fmSyKdczcThr3rwvxKaPCRfQC6Rz7g8L5uxWGe53Tsm',
            in_pool: false,
            in_blocks: ['Baaa'],
            is_operation_final: true,
            thread: 3,
            operation: {
              content: {
                fee: '0.01',
                expire_period: 638742,
                op: {
                  CallSC: {
                    target_addr:
                      'AS1gt69gqYD92dqPyE6DBRJ7KjpnQHqFzFs2YCkBcSnuxX5bGhBC',
                    target_func: 'transfer',
                    param: [52, 0, 0],
                    max_gas: 2520000,
                    coins: '0',
                  },
                },
              },
              signature:
                '1YzWNWdhuZCuKWozSaYBfggp6PKScmbfABv7FLCDBmC4TNWspkWCJutaBxmPJ8yZzJyZknfeVtjmy2bfg2exEuCtqgb8DN',
              content_creator_pub_key:
                'P12wGSJoFJuP4ozRMphLU3VvU3PqQL336a4vA8eSyZJxcvQur4Cp',
              content_creator_address:
                'AU1Bknx3Du4aiGiHaeh1vo7LfwEPRF3piAwotRkdK975qCBxWwLs',
              id: 'O1YUEyP4fmSyKdczcThr3rwvxKaPCRfQC6Rz7g8L5uxWGe53Tsm',
            },
            op_exec_status: true,
          },
        ],
        id: 0,
      };
    }

    if (method === 'get_status') {
      return {
        jsonrpc: '2.0',
        result: {
          node_id: 'N12sNdL7YwSawpnJrk9XCWDjKbgfNamAobp62AX5qfkgpBkGh2wC',
          node_ip: '149.202.84.39',
          version: 'DEVN.28.3',
          current_time: 1732005865705,
          current_cycle: 13533,
          current_cycle_time: 1732005384000,
          next_cycle_time: 1732007432000,
          connected_nodes: {
            N1DZb3ao8BEtdsYP1KYyWacpTENHDrQboxGaYDL4U8MQppaxvzo: [
              '::ffff:149.202.65.130',
              false,
            ],
            N1NnuSW48GKGaYZamAVKXfXbbnt3StxWoHpYtBZSJvY9e8U1BTC: [
              '37.187.156.118',
              true,
            ],
            N1kKfgrCveVnosUkxTzaBw5cf9f2cbTvK3R5Ssb2Pf76au8xwmH: [
              '149.202.84.7',
              true,
            ],
          },
          last_slot: {
            period: 1732254,
            thread: 3,
          },
          next_slot: {
            period: 1732254,
            thread: 4,
          },
          consensus_stats: {
            start_timespan: 1732005805705,
            end_timespan: 1732005865705,
            final_block_count: 120,
            stale_block_count: 0,
            clique_count: 1,
          },
          pool_stats: [0, 0],
          network_stats: {
            in_connection_count: 1,
            out_connection_count: 2,
            known_peer_count: 9,
            banned_peer_count: 0,
            active_node_count: 3,
          },
          execution_stats: {
            time_window_start: 1732005805705,
            time_window_end: 1732005865705,
            final_block_count: 120,
            final_executed_operations_count: 3,
            active_cursor: {
              period: 1732253,
              thread: 31,
            },
            final_cursor: {
              period: 1732252,
              thread: 2,
            },
          },
          config: {
            genesis_timestamp: 1704289800000,
            end_timestamp: null,
            thread_count: 32,
            t0: 16000,
            delta_f0: 1088,
            operation_validity_periods: 10,
            periods_per_cycle: 128,
            block_reward: '1.02',
            roll_price: '100',
            max_block_size: 300000,
          },
          chain_id: 77658366,
          minimal_fees: '0.01',
        },
        id: 0,
      };
    }

    return {
      jsonrpc: '2.0',
      result: [],
      id: 1,
    };
  });

  server.post(
    'https://station.massa/plugin/massa-labs/massa-wallet/api/accounts/bridge/transfer',
    () => {
      return {
        operationId: 'O15nBJ7b6t5tNGqE992H6gMnBF2M3t6HuZ4whkzg5sBAmkMd2eL',
      };
    },
  );

  server.post('https://station.massa/cmd/executeFunction', () => {
    return {
      firstEvent: {
        address: 'AS1gt69gqYD92dqPyE6DBRJ7KjpnQHqFzFs2YCkBcSnuxX5bGhBC',
        data: 'Function called successfully but did not wait for event',
      },
      operationId: 'O1YUEyP4fmSyKdczcThr3rwvxKaPCRfQC6Rz7g8L5uxWGe53Tsm',
    };
  });
};

export { handlers, otherDomainHandlers };
