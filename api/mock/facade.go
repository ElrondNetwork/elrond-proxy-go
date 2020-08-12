package mock

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-proxy-go/data"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
)

// Facade is the mock implementation of a node's router handler
type Facade struct {
	IsFaucetEnabledHandler                      func() bool
	GetAccountHandler                           func(address string) (*data.Account, error)
	GetValueForKeyHandler                       func(address string, key string) (string, error)
	GetTransactionsHandler                      func(address string) ([]data.DatabaseTransaction, error)
	GetTransactionHandler                       func(txHash string) (*data.FullTransaction, error)
	SendTransactionHandler                      func(tx *data.Transaction) (int, string, error)
	SendMultipleTransactionsHandler             func(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error)
	SendUserFundsCalled                         func(receiver string, value *big.Int) error
	ExecuteSCQueryHandler                       func(query *data.SCQuery) (*vmcommon.VMOutput, error)
	GetHeartbeatDataHandler                     func() (*data.HeartbeatResponse, error)
	ValidatorStatisticsHandler                  func() (map[string]*data.ValidatorApiResponse, error)
	TransactionCostRequestHandler               func(tx *data.Transaction) (string, error)
	GetTransactionStatusHandler                 func(txHash string, sender string) (string, error)
	GetConfigMetricsHandler                     func() (*data.GenericAPIResponse, error)
	GetNetworkMetricsHandler                    func(shardID uint32) (*data.GenericAPIResponse, error)
	GetBlockByShardIDAndNonceHandler            func(shardID uint32, nonce uint64) (data.AtlasBlock, error)
	GetTransactionByHashAndSenderAddressHandler func(txHash string, sndAddr string) (*data.FullTransaction, int, error)
	GetBlockByHashCalled                        func(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error)
	GetBlockByNonceCalled                       func(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error)
	GetHyperBlockByHashCalled                   func(hash string) (*data.HyperblockApiResponse, error)
	GetHyperBlockByNonceCalled                  func(nonce uint64) (*data.HyperblockApiResponse, error)
}

// IsFaucetEnabled -
func (f *Facade) IsFaucetEnabled() bool {
	if f.IsFaucetEnabledHandler != nil {
		return f.IsFaucetEnabledHandler()
	}

	return true
}

// GetNetworkStatusMetrics -
func (f *Facade) GetNetworkStatusMetrics(shardID uint32) (*data.GenericAPIResponse, error) {
	if f.GetNetworkMetricsHandler != nil {
		return f.GetNetworkMetricsHandler(shardID)
	}

	return nil, nil
}

// GetNetworkConfigMetrics -
func (f *Facade) GetNetworkConfigMetrics() (*data.GenericAPIResponse, error) {
	if f.GetConfigMetricsHandler != nil {
		return f.GetConfigMetricsHandler()
	}

	return nil, nil
}

// ValidatorStatistics is the mock implementation of a handler's ValidatorStatistics method
func (f *Facade) ValidatorStatistics() (map[string]*data.ValidatorApiResponse, error) {
	return f.ValidatorStatisticsHandler()
}

// GetAccount is the mock implementation of a handler's GetAccount method
func (f *Facade) GetAccount(address string) (*data.Account, error) {
	return f.GetAccountHandler(address)
}

// GetValueForKey --
func (f *Facade) GetValueForKey(address string, key string) (string, error) {
	return f.GetValueForKeyHandler(address, key)
}

// GetTransactions --
func (f *Facade) GetTransactions(address string) ([]data.DatabaseTransaction, error) {
	return f.GetTransactionsHandler(address)
}

// GetTransactionByHashAndSenderAddressCalled --
func (f *Facade) GetTransactionByHashAndSenderAddress(txHash string, sndAddr string) (*data.FullTransaction, int, error) {
	return f.GetTransactionByHashAndSenderAddressHandler(txHash, sndAddr)
}

// GetTransaction --
func (f *Facade) GetTransaction(txHash string) (*data.FullTransaction, error) {
	return f.GetTransactionHandler(txHash)
}

// SendTransaction is the mock implementation of a handler's SendTransaction method
func (f *Facade) SendTransaction(tx *data.Transaction) (int, string, error) {
	return f.SendTransactionHandler(tx)
}

// SendMultipleTransactions is the mock implementation of a handler's SendMultipleTransactions method
func (f *Facade) SendMultipleTransactions(txs []*data.Transaction) (data.MultipleTransactionsResponseData, error) {
	return f.SendMultipleTransactionsHandler(txs)
}

// TransactionCostRequest --
func (f *Facade) TransactionCostRequest(tx *data.Transaction) (string, error) {
	return f.TransactionCostRequestHandler(tx)
}

// GetTransactionStatus --
func (f *Facade) GetTransactionStatus(txHash string, sender string) (string, error) {
	return f.GetTransactionStatusHandler(txHash, sender)
}

// SendUserFunds is the mock implementation of a handler's SendUserFunds method
func (f *Facade) SendUserFunds(receiver string, value *big.Int) error {
	return f.SendUserFundsCalled(receiver, value)
}

// ExecuteSCQuery is a mock implementation.
func (f *Facade) ExecuteSCQuery(query *data.SCQuery) (*vmcommon.VMOutput, error) {
	return f.ExecuteSCQueryHandler(query)
}

// GetHeartbeatData is the mock implementation of a handler's GetHeartbeatData method
func (f *Facade) GetHeartbeatData() (*data.HeartbeatResponse, error) {
	return f.GetHeartbeatDataHandler()
}

// GetAtlasBlockByShardIDAndNonce -
func (f *Facade) GetAtlasBlockByShardIDAndNonce(shardID uint32, nonce uint64) (data.AtlasBlock, error) {
	return f.GetBlockByShardIDAndNonceHandler(shardID, nonce)
}

// GetBlockByHash -
func (f *Facade) GetBlockByHash(shardID uint32, hash string, withTxs bool) (*data.BlockApiResponse, error) {
	return f.GetBlockByHashCalled(shardID, hash, withTxs)
}

// GetBlockByHash -
func (f *Facade) GetBlockByNonce(shardID uint32, nonce uint64, withTxs bool) (*data.BlockApiResponse, error) {
	return f.GetBlockByNonceCalled(shardID, nonce, withTxs)
}

// GetHyperBlockByHash -
func (f *Facade) GetHyperBlockByHash(hash string) (*data.HyperblockApiResponse, error) {
	if f.GetHyperBlockByHashCalled != nil {
		return f.GetHyperBlockByHashCalled(hash)
	}

	panic("not implemented: GetHyperBlockByHash")
}

// GetHyperBlockByNonce -
func (f *Facade) GetHyperBlockByNonce(nonce uint64) (*data.HyperblockApiResponse, error) {
	if f.GetHyperBlockByNonceCalled != nil {
		return f.GetHyperBlockByNonceCalled(nonce)
	}

	panic("not implemented: GetHyperBlockByNonce")
}

// WrongFacade is a struct that can be used as a wrong implementation of the node router handler
type WrongFacade struct {
}
