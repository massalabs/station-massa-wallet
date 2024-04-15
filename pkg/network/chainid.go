package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/massalabs/station/pkg/node"
)

var ErrChainIDNotInStatus = errors.New("chain id not in status")

func GetNodeInfo() (uint64, string, error) {
	client, err := NewMassaClient()
	if err != nil {
		return 0, "", err
	}

	status, err := status(client)
	if err != nil {
		return 0, "", err
	}

	chainID := status.ChainID

	if chainID == nil {
		return 0, "", ErrChainIDNotInStatus
	}

	var minimalFees string
	if status.MinimalFees == nil {
		minimalFees = "0"
	} else {
		minimalFees = *status.MinimalFees
	}

	return uint64(*chainID), minimalFees, nil
}

// TODO: remove the copied code from station lib

type State struct {
	Config         *Config         `json:"config"`
	ConsensusStats *ConsensusStats `json:"consensus_stats"`
	CurrentCycle   *uint           `json:"current_cycle"`
	CurrentTime    *uint           `json:"current_time"`
	LastSlot       *node.Slot      `json:"last_slot"`
	NetworkStats   *NetworkStats   `json:"network_stats"`
	NextSlot       *node.Slot      `json:"next_slot"`
	NodeID         *string         `json:"node_id"`
	NodeIP         *string         `json:"node_ip"`
	PoolStats      *[]uint         `json:"pool_stats"`
	Version        *string         `json:"version"`
	ChainID        *uint           `json:"chain_id"`
	MinimalFees    *string         `json:"minimal_fees"`
}

//nolint:tagliatelle
type Config struct {
	BlockReward             *string `json:"block_reward"`
	DeltaF0                 *uint   `json:"delta_f0"`
	EndTimeStamp            *uint   `json:"end_timestamp"`
	GenesisTimestamp        *uint   `json:"genesis_timestamp"`
	OperationValidityParios *uint   `json:"operation_validity_periods"`
	PeriodsPerCycle         *uint   `json:"periods_per_cycle"`
	PosLockCycles           *uint   `json:"pos_lock_cycles"`
	PosLookbackCycle        *uint   `json:"pos_lookback_cycles"`
	RollPrice               *string `json:"roll_price"`
	T0                      *uint   `json:"t0"`
	ThreadCount             *uint   `json:"thread_count"`
}

type ConsensusStats struct {
	CliqueCount         *uint `json:"clique_count"`
	EndTimespan         *uint `json:"end_timespan"`
	FinalBlockCount     *uint `json:"final_block_count"`
	FinalOperationCount *uint `json:"final_operation_count"`
	StakerCount         *uint `json:"staker_count"`
	StaleBlockCount     *uint `json:"stale_block_count"`
	StartTimespan       *uint `json:"start_timespan"`
}

//nolint:tagliatelle
type NetworkStats struct {
	ActiveNodeCount    *uint `json:"active_node_count"`
	BannedPeerCount    *uint `json:"banned_peer_count"`
	InConnectionCount  *uint `json:"in_connection_count"`
	KnowPeerCount      *uint `json:"known_peer_count"`
	OutConnectionCount *uint `json:"out_connection_count"`
}

func status(client *node.Client) (*State, error) {
	rawResponse, err := client.RPCClient.Call(
		context.Background(),
		"get_status",
	)
	if err != nil {
		return nil, fmt.Errorf("calling get_status: %w", err)
	}

	if rawResponse.Error != nil {
		return nil, rawResponse.Error
	}

	var resp State

	err = rawResponse.GetObject(&resp)
	if err != nil {
		return nil, fmt.Errorf("parsing get_status jsonrpc response '%+v': %w", rawResponse, err)
	}

	return &resp, nil
}
