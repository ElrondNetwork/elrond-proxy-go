package facade_test

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/elrond-go/core"
	"github.com/ElrondNetwork/elrond-go/core/pubkeyConverter"
	"github.com/ElrondNetwork/elrond-go/crypto"
	"github.com/ElrondNetwork/elrond-go/crypto/signing"
	"github.com/ElrondNetwork/elrond-go/crypto/signing/ed25519"
	"github.com/ElrondNetwork/elrond-go/data/vm"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
	"github.com/ElrondNetwork/elrond-proxy-go/facade"
	"github.com/ElrondNetwork/elrond-proxy-go/facade/mock"
	"github.com/stretchr/testify/assert"
)

var publicKeyConverter, _ = pubkeyConverter.NewBech32PubkeyConverter(32)

func TestNewElrondProxyFacade_NilActionsProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		nil,
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilActionsProcessor, err)
}

func TestNewElrondProxyFacade_NilAccountProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		nil,
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilAccountProcessor, err)
}

func TestNewElrondProxyFacade_NilTransactionProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		nil,
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilTransactionProcessor, err)
}

func TestNewElrondProxyFacade_NilGetValuesProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		nil,
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilSCQueryService, err)
}

func TestNewElrondProxyFacade_NilHeartbeatProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		nil,
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilHeartbeatProcessor, err)
}

func TestNewElrondProxyFacade_NilValStatsProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		nil,
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilValidatorStatisticsProcessor, err)
}

func TestNewElrondProxyFacade_NilFaucetProcShouldErr(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		nil,
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilFaucetProcessor, err)
}

func TestNewElrondProxyFacade_NilNodeProcessor(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		nil,
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilNodeStatusProcessor, err)
}

func TestNewElrondProxyFacade_NilProofProcessor(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		nil,
		publicKeyConverter,
	)

	assert.Nil(t, epf)
	assert.Equal(t, facade.ErrNilProofProcessor, err)
}

func TestNewElrondProxyFacade_ShouldWork(t *testing.T) {
	t.Parallel()

	epf, err := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	assert.NotNil(t, epf)
	assert.Nil(t, err)
}

func TestElrondProxyFacade_GetAccount(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{
			GetAccountCalled: func(address string) (account *data.Account, e error) {
				wasCalled = true
				return &data.Account{}, nil
			},
		},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	_, _ = epf.GetAccount("")

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_SendTransaction(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{
			SendTransactionCalled: func(tx *data.Transaction) (int, string, error) {
				wasCalled = true

				return 0, "", nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	_, _, _ = epf.SendTransaction(&data.Transaction{})

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_SimulateTransaction(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{
			SimulateTransactionCalled: func(tx *data.Transaction, checkSignature bool) (*data.GenericAPIResponse, error) {
				wasCalled = true
				return nil, nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	_, _ = epf.SimulateTransaction(&data.Transaction{}, false)

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_SendUserFunds(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{
			GetAccountCalled: func(address string) (*data.Account, error) {
				return &data.Account{
					Nonce: uint64(0),
				}, nil
			},
		},
		&mock.TransactionProcessorStub{
			SendTransactionCalled: func(tx *data.Transaction) (int, string, error) {
				wasCalled = true
				return 0, "", nil
			},
		},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{
			SenderDetailsFromPemCalled: func(receiver string) (crypto.PrivateKey, string, error) {
				return getPrivKey(), "rcvr", nil
			},
			GenerateTxForSendUserFundsCalled: func(senderSk crypto.PrivateKey, senderPk string, senderNonce uint64, receiver string, value *big.Int, chainID string, version uint32) (*data.Transaction, error) {
				return &data.Transaction{}, nil
			},
		},
		&mock.NodeStatusProcessorStub{
			GetConfigMetricsCalled: func() (*data.GenericAPIResponse, error) {
				return &data.GenericAPIResponse{
					Data: map[string]interface{}{
						"config": map[string]interface{}{
							core.MetricChainId:               "chainID",
							core.MetricMinTransactionVersion: 1.0,
						},
					},
				}, nil
			},
		},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	_ = epf.SendUserFunds("", big.NewInt(0))

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_GetDataValue(t *testing.T) {
	t.Parallel()

	wasCalled := false
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{
			ExecuteQueryCalled: func(query *data.SCQuery) (*vm.VMOutputApi, error) {
				wasCalled = true
				return &vm.VMOutputApi{}, nil
			},
		},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	_, _ = epf.ExecuteSCQuery(nil)

	assert.True(t, wasCalled)
}

func TestElrondProxyFacade_GetHeartbeatData(t *testing.T) {
	t.Parallel()

	expectedResults := &data.HeartbeatResponse{
		Heartbeats: []data.PubKeyHeartbeat{
			{
				ReceivedShardID: 0,
				ComputedShardID: 1,
			},
		},
	}
	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{
			GetHeartbeatDataCalled: func() (*data.HeartbeatResponse, error) {
				return expectedResults, nil
			},
		},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	actualResult, _ := epf.GetHeartbeatData()

	assert.Equal(t, expectedResults, actualResult)
}

func TestElrondProxyFacade_ReloadObservers(t *testing.T) {
	t.Parallel()

	expectedResult := data.NodesReloadResponse{
		Description: "abc",
		Error:       "bca",
	}

	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{
			ReloadObserversCalled: func() data.NodesReloadResponse {
				return expectedResult
			},
		},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	actualResult := epf.ReloadObservers()

	assert.Equal(t, expectedResult, actualResult)
}

func TestElrondProxyFacade_ReloadFullHistoryObservers(t *testing.T) {
	t.Parallel()

	expectedResult := data.NodesReloadResponse{
		Description: "abc",
		Error:       "bca",
	}

	epf, _ := facade.NewElrondProxyFacade(
		&mock.ActionsProcessorStub{
			ReloadFullHistoryObserversCalled: func() data.NodesReloadResponse {
				return expectedResult
			},
		},
		&mock.AccountProcessorStub{},
		&mock.TransactionProcessorStub{},
		&mock.SCQueryServiceStub{},
		&mock.HeartbeatProcessorStub{},
		&mock.ValidatorStatisticsProcessorStub{},
		&mock.FaucetProcessorStub{},
		&mock.NodeStatusProcessorStub{},
		&mock.BlockProcessorStub{},
		&mock.ProofProcessorStub{},
		publicKeyConverter,
	)

	actualResult := epf.ReloadFullHistoryObservers()

	assert.Equal(t, expectedResult, actualResult)
}

func getPrivKey() crypto.PrivateKey {
	keyGen := signing.NewKeyGenerator(ed25519.NewEd25519())
	sk, _ := keyGen.GeneratePair()

	return sk
}
