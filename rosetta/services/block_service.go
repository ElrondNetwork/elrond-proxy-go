package services

import (
	"context"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/client"
	"github.com/ElrondNetwork/elrond-proxy-go/rosetta/configuration"
	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
)

type blockAPIService struct {
	elrondClient client.ElrondClientHandler
	txsParser    *transactionsParser
}

// NewBlockAPIService will create a new instance of blockAPIService
func NewBlockAPIService(elrondClient client.ElrondClientHandler, cfg *configuration.Configuration) server.BlockAPIServicer {
	return &blockAPIService{
		elrondClient: elrondClient,
		txsParser:    newTransactionParser(cfg),
	}
}

// Block implements the /block endpoint.
func (bas *blockAPIService) Block(
	_ context.Context,
	request *types.BlockRequest,
) (*types.BlockResponse, *types.Error) {
	if request.BlockIdentifier.Index != nil {
		return bas.getBlockByNonce(*request.BlockIdentifier.Index)
	}

	if request.BlockIdentifier.Hash != nil {
		return bas.getBlockByHash(*request.BlockIdentifier.Hash)
	}

	return nil, ErrMustQueryByIndexOrByHash
}

func (bas *blockAPIService) getBlockByNonce(nonce int64) (*types.BlockResponse, *types.Error) {
	hyperBlock, err := bas.elrondClient.GetBlockByNonce(nonce)
	if err != nil {
		return nil, wrapErr(ErrUnableToGetBlock, err)
	}

	return bas.parseHyperBlock(hyperBlock)
}

func (bas *blockAPIService) getBlockByHash(hash string) (*types.BlockResponse, *types.Error) {
	hyperBlock, err := bas.elrondClient.GetBlockByHash(hash)
	if err != nil {
		return nil, wrapErr(ErrUnableToGetBlock, err)
	}

	return bas.parseHyperBlock(hyperBlock)
}

func (bas *blockAPIService) parseHyperBlock(hyperBlock *data.Hyperblock) (*types.BlockResponse, *types.Error) {
	var parentBlockIdentifier *types.BlockIdentifier
	if hyperBlock.Nonce != 0 {
		parentBlockIdentifier = &types.BlockIdentifier{
			Index: int64(hyperBlock.Nonce - 1),
			Hash:  hyperBlock.PrevBlockHash,
		}
	}

	return &types.BlockResponse{
		Block: &types.Block{
			BlockIdentifier: &types.BlockIdentifier{
				Index: int64(hyperBlock.Nonce),
				Hash:  hyperBlock.Hash,
			},
			ParentBlockIdentifier: parentBlockIdentifier,
			Timestamp:             bas.elrondClient.CalculateBlockTimestampUnix(hyperBlock.Round),
			Transactions:          bas.txsParser.parseTxsFromHyperBlock(hyperBlock),
			Metadata: objectsMap{
				"epoch": hyperBlock.Epoch,
				"round": hyperBlock.Round,
			},
		},
	}, nil
}

// BlockTransaction - not implemented
// We dont need this method because all transactions are returned by method Block
func (bas *blockAPIService) BlockTransaction(
	_ context.Context,
	_ *types.BlockTransactionRequest,
) (*types.BlockTransactionResponse, *types.Error) {
	return nil, ErrNotImplemented
}
