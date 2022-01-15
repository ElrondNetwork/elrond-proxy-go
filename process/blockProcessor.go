package process

import (
	"fmt"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/check"
	"github.com/ElrondNetwork/elrond-proxy-go/common"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

const (
	blockByHashPath  = "/block/by-hash"
	blockByNoncePath = "/block/by-nonce"
	withTxsParamTrue = "?withTxs=true"

	internalMetaBlockByHashPath  = "/internal/%s/metablock/by-hash"
	internalShardBlockByHashPath = "/internal/%s/shardblock/by-hash"

	internalMetaBlockByNoncePath  = "/internal/%s/metablock/by-nonce"
	internalShardBlockByNoncePath = "/internal/%s/shardblock/by-nonce"
)

const (
	jsonPathStr = "json"
	rawPathStr  = "raw"
)

// BlockProcessor handles blocks retrieving
type BlockProcessor struct {
	proc     Processor
	dbReader ExternalStorageConnector
}

// NewBlockProcessor will create a new block processor
func NewBlockProcessor(dbReader ExternalStorageConnector, proc Processor) (*BlockProcessor, error) {
	if check.IfNil(dbReader) {
		return nil, ErrNilDatabaseConnector
	}
	if check.IfNil(proc) {
		return nil, ErrNilCoreProcessor
	}

	return &BlockProcessor{
		dbReader: dbReader,
		proc:     proc,
	}, nil
}

// GetAtlasBlockByShardIDAndNonce return the block byte shardID and nonce
func (bp *BlockProcessor) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	return bp.dbReader.GetAtlasBlockByShardIDAndNonce(shardID, nonce)
}

// GetBlockByHash will return the block based on its hash
func (bp *BlockProcessor) GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%s", blockByHashPath, hash)
	if withTxs {
		path += withTxsParamTrue
	}

	for _, observer := range observers {
		var response data.BlockApiResponse

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("block request", "shard id", observer.ShardId, "hash", hash, "observer", observer.Address)
		return &response, nil

	}

	return nil, ErrSendingRequest
}

// GetBlockByNonce will return the block based on the nonce
func (bp *BlockProcessor) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("%s/%d", blockByNoncePath, nonce)
	if withTxs {
		path += withTxsParamTrue
	}

	for _, observer := range observers {
		var response data.BlockApiResponse

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("block request", "shard id", observer.ShardId, "nonce", nonce, "observer", observer.Address)
		return &response, nil

	}

	return nil, ErrSendingRequest
}

func (bp *BlockProcessor) getObserversOrFullHistoryNodes(shardID uint32) ([]*data.NodeData, error) {
	fullHistoryNodes, err := bp.proc.GetFullHistoryNodes(shardID)
	if err == nil {
		return fullHistoryNodes, nil
	}

	return bp.proc.GetObservers(shardID)
}

// GetHyperBlockByHash returns the hyperblock by hash
func (bp *BlockProcessor) GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error) {
	builder := &HyperblockBuilder{}

	metaBlockResponse, err := bp.GetBlockByHash(core.MetachainShardId, hash, true)
	if err != nil {
		return nil, err
	}

	metaBlock := metaBlockResponse.Data.Block
	builder.addMetaBlock(&metaBlock)

	for _, notarizedBlock := range metaBlock.NotarizedBlocks {
		shardBlockResponse, err := bp.GetBlockByHash(notarizedBlock.Shard, notarizedBlock.Hash, true)
		if err != nil {
			return nil, err
		}

		builder.addShardBlock(&shardBlockResponse.Data.Block)
	}

	hyperblock := builder.build()
	return data.NewHyperblockApiResponse(hyperblock), nil
}

// GetHyperBlockByNonce returns the hyperblock by nonce
func (bp *BlockProcessor) GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error) {
	builder := &HyperblockBuilder{}

	metaBlockResponse, err := bp.GetBlockByNonce(core.MetachainShardId, nonce, true)
	if err != nil {
		return nil, err
	}

	metaBlock := metaBlockResponse.Data.Block
	builder.addMetaBlock(&metaBlock)

	for _, notarizedBlock := range metaBlock.NotarizedBlocks {
		shardBlockResponse, err := bp.GetBlockByHash(notarizedBlock.Shard, notarizedBlock.Hash, true)
		if err != nil {
			return nil, err
		}

		builder.addShardBlock(&shardBlockResponse.Data.Block)
	}

	hyperblock := builder.build()
	return data.NewHyperblockApiResponse(hyperblock), nil
}

// GetInternalBlockByHash will return the block based on its hash
func (bp *BlockProcessor) GetInternalBlockByHash(shardID uint32, hash string, format common.OutportFormat) (*data.InternalBlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path, err := getInternalBlockPath(shardID, format, hash)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var response data.InternalBlockApiResponse

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("block request", "shard id", observer.ShardId, "hash", hash, "observer", observer.Address)
		return &response, nil

	}

	return nil, ErrSendingRequest
}

// GetInternalBlockByNonce will return the block based on its nonce
func (bp *BlockProcessor) GetInternalBlockByNonce(shardID uint32, nonce uint64, format common.OutportFormat) (*data.InternalBlockApiResponse, error) {
	observers, err := bp.getObserversOrFullHistoryNodes(shardID)
	if err != nil {
		return nil, err
	}

	path, err := getInternalBlockPath(shardID, format, nonce)
	if err != nil {
		return nil, err
	}

	for _, observer := range observers {
		var response data.InternalBlockApiResponse

		_, err := bp.proc.CallGetRestEndPoint(observer.Address, path, &response)
		if err != nil {
			log.Error("block request", "observer", observer.Address, "error", err.Error())
			continue
		}

		log.Info("block request", "shard id", observer.ShardId, "round", nonce, "observer", observer.Address)
		return &response, nil

	}

	return nil, ErrSendingRequest
}

func getInternalBlockPath(shardID uint32, format common.OutportFormat, value interface{}) (string, error) {
	var path string

	var outportStr string
	switch format {
	case common.Internal:
		outportStr = jsonPathStr
	case common.Proto:
		outportStr = rawPathStr
	default:
		return "", ErrInvalidOutportFormat
	}

	var strFormat string
	switch value.(type) {
	case uint64:
		strFormat = "%s/%d"
	case string:
		strFormat = "%s/%s"
	default:
		strFormat = "%s/%v"
	}

	if shardID == core.MetachainShardId {
		path = fmt.Sprintf(internalMetaBlockByNoncePath, outportStr)
	} else {
		path = fmt.Sprintf(internalShardBlockByNoncePath, outportStr)
	}

	return fmt.Sprintf(strFormat, path, value), nil
}
