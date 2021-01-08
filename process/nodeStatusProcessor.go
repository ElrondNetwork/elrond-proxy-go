package process

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

// NetworkStatusPath represents the path where an observer exposes his network metrics
const NetworkStatusPath = "/network/status"

// NetworkConfigPath represents the path where an observer exposes his network metrics
const NetworkConfigPath = "/network/config"

// EconomicsDataPath represents the path where an observer exposes his economics data
const EconomicsDataPath = "/network/economics"

// NetworkConfigPath represents the path where an observer exposes his node status metrics
const NodeStatusPath = "/node/status"

// TotalStakedPath represents the path where an observer exposes the total staked value from validators contract
const TotalStakedPath = "/network/total-staked"

// NodeStatusProcessor handles the action needed for fetching data related to status metrics from nodes
type NodeStatusProcessor struct {
	proc Processor
}

// NewNodeStatusProcessor creates a new instance of NodeStatusProcessor
func NewNodeStatusProcessor(processor Processor) (*NodeStatusProcessor, error) {
	if check.IfNil(processor) {
		return nil, ErrNilCoreProcessor
	}

	return &NodeStatusProcessor{
		proc: processor,
	}, nil
}

// GetNetworkStatusMetrics will simply forward the network status metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseNetworkMetrics *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, NetworkStatusPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("network metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network metrics request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil

	}

	return nil, ErrSendingRequest
}

// GetTotalStaked will simply forward the total staked value from a metachain observer
func (nsp *NodeStatusProcessor) GetTotalStaked() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseTotalStaked *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, TotalStakedPath, &responseTotalStaked)
		if err != nil {
			log.Error("total staked request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("total staked request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseTotalStaked, nil

	}

	return nil, ErrSendingRequest
}

// GetNetworkConfigMetrics will simply forward the network config metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetAllObservers()
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseNetworkMetrics *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, NetworkConfigPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("network metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("network metrics request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil

	}

	return nil, ErrSendingRequest
}

// GetNetworkConfigMetrics will simply forward the network config metrics from an observer in the given shard
func (nsp *NodeStatusProcessor) GetEconomicsDataMetrics() (*data.GenericAPIResponse, error) {
	metaObservers, err := nsp.proc.GetObservers(core.MetachainShardId)
	if err != nil {
		return nil, err
	}

	metaResponse, err := nsp.getEconomicsDataMetrics(metaObservers)
	if err == nil {
		return metaResponse, nil
	}

	log.Warn("cannot get economics data metrics from metachain observer. will try with all observers",
		"error", err)

	allObservers, err := nsp.proc.GetAllObservers()
	if err != nil {
		return nil, err
	}

	return nsp.getEconomicsDataMetrics(allObservers)
}

func (nsp *NodeStatusProcessor) getEconomicsDataMetrics(observers []*data.NodeData) (*data.GenericAPIResponse, error) {
	for _, observer := range observers {
		var responseNetworkMetrics *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, EconomicsDataPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("economics data request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("economics data request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil
	}

	return nil, ErrSendingRequest
}

func (nsp *NodeStatusProcessor) getNodeStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	observers, err := nsp.proc.GetObservers(shardID)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var responseNetworkMetrics *data.GenericAPIResponse

		_, err := nsp.proc.CallGetRestEndPoint(observer.Address, NodeStatusPath, &responseNetworkMetrics)
		if err != nil {
			log.Error("node status metrics request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("node status metrics request", "shard id", observer.ShardId, "observer", observer.Address)
		return responseNetworkMetrics, nil

	}

	return nil, ErrSendingRequest
}

// GetLatestFullySynchronizedHyperblockNonce will compute nonce of the latest hyperblock that can be returned
func (nsp *NodeStatusProcessor) GetLatestFullySynchronizedHyperblockNonce() (uint64, error) {
	shardsIDs, err := nsp.getShardsIDs()
	if err != nil {
		return 0, err
	}

	nonces := make([]uint64, 0)
	for shardID := range shardsIDs {
		nodeStatusResponse, err := nsp.getNodeStatusMetrics(shardID)
		if err != nil {
			return 0, err
		}

		if nodeStatusResponse.Error != "" {
			return 0, errors.New(nodeStatusResponse.Error)
		}

		var nonce uint64
		var ok bool
		if shardID == core.MetachainShardId {
			nonce, ok = getNonceFromMetachainStatus(nodeStatusResponse.Data)
		} else {
			nonce, ok = getNonceFromShardStatus(nodeStatusResponse.Data)
		}
		if !ok {
			return 0, ErrCannotParseNodeStatusMetrics
		}

		nonces = append(nonces, nonce)
	}

	return getMinNonce(nonces), nil
}

func getMinNonce(noncesSlice []uint64) uint64 {
	// initialize min with max uint64 value
	min := uint64(math.MaxUint64)
	for _, value := range noncesSlice {
		if value < min {
			min = value
		}
	}

	return min
}

func (nsp *NodeStatusProcessor) getShardsIDs() (map[uint32]struct{}, error) {
	observers, err := nsp.proc.GetAllObservers()
	if err != nil {
		return nil, err
	}

	shardsIDs := make(map[uint32]struct{})
	for _, observer := range observers {
		shardsIDs[observer.ShardId] = struct{}{}
	}

	if len(shardsIDs) == 0 {
		return nil, ErrMissingObserver
	}

	return shardsIDs, nil
}

func getNonceFromShardStatus(nodeStatusData interface{}) (uint64, bool) {
	metric, ok := getMetric(nodeStatusData, core.MetricCrossCheckBlockHeight)
	if !ok {
		return 0, false
	}

	return parseMetricCrossCheckBlockHeight(metric)
}

func getNonceFromMetachainStatus(nodeStatusData interface{}) (uint64, bool) {
	metric, ok := getMetric(nodeStatusData, core.MetricNonce)
	if !ok {
		return 0, false
	}

	return getUint(metric), true
}

func getMetric(nodeStatusData interface{}, metric string) (interface{}, bool) {
	metricsMapI, ok := nodeStatusData.(map[string]interface{})
	if !ok {
		return nil, false
	}

	metricsMap, ok := metricsMapI["metrics"]
	if !ok {
		return nil, false
	}

	metrics, ok := metricsMap.(map[string]interface{})
	if !ok {
		return nil, false
	}

	value, ok := metrics[metric]
	if !ok {
		return nil, false
	}

	return value, true
}

func parseMetricCrossCheckBlockHeight(value interface{}) (uint64, bool) {
	valueStr, ok := value.(string)
	if !ok {
		return 0, false
	}

	// metric looks like that
	// "meta 886717"
	values := strings.Split(valueStr, " ")
	if len(values) < 2 {
		return 0, false
	}

	nonce, err := strconv.ParseUint(values[1], 10, 64)
	if err != nil {
		return 0, false
	}

	return nonce, true
}

func getUint(value interface{}) uint64 {
	valueFloat, ok := value.(float64)
	if !ok {
		return 0
	}

	return uint64(valueFloat)
}
