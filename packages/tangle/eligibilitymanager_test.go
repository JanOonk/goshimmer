package tangle

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"

	"github.com/iotaledger/hive.go/identity"
)

func TestDependenciesConfirmed(t *testing.T) {
	tangle := newTestTangle()
	defer tangle.Shutdown()
	wallets, walletsByAddress, genesisTransaction := setupEligibilityTests(t, tangle)

	messages, transactions, _, _, _ := scenario1MessagesApproveEmptyID(t, tangle, genesisTransaction, wallets, walletsByAddress)
	mockUTXO := ledgerstate.NewUtxoDagMock(t, tangle.LedgerState.UTXODAG)
	mockUTXO.On("InclusionState", genesisTransaction.ID()).Return(ledgerstate.Confirmed)
	mockUTXO.On("InclusionState", transactions["2"].ID()).Return(ledgerstate.Confirmed)
	// TODO mocking does not work
	tangle.LedgerState.UTXODAG = mockUTXO
	err := tangle.EligibilityManager.checkEligibility(messages["1"].ID())
	assert.NoError(t, err)
	var eligibilityResult1 bool
	tangle.Storage.MessageMetadata(messages["1"].ID()).Consume(func(messageMetadata *MessageMetadata) {
		eligibilityResult1 = messageMetadata.IsEligible()
	})

	assert.True(t, eligibilityResult1)
	// TODO check if event message eligible was triggered

}

func TestDataMessageAlwaysEligible(t *testing.T) {
	tangle := newTestTangle()
	defer tangle.Shutdown()
	tangle.EligibilityManager.Setup()

	message := newTestDataMessage("data")
	tangle.Storage.StoreMessage(message)

	err := tangle.EligibilityManager.checkEligibility(message.ID())
	assert.NoError(t, err)

	var eligibilityResult1 bool
	tangle.Storage.MessageMetadata(message.ID()).Consume(func(messageMetadata *MessageMetadata) {
		eligibilityResult1 = messageMetadata.IsEligible()
	})

	assert.True(t, eligibilityResult1)
}

func TestDependencyDirectApproval(t *testing.T) {
	tangle := newTestTangle()
	defer tangle.Shutdown()

	//wallets, walletsByAddress, genesisTransaction := setupEligibilityTests(t, tangle)
	//messages, transactions, _, _, _ := scenario2MessagesApproveDependency(t, tangle, genesisTransaction, wallets, walletsByAddress)

}

func setupEligibilityTests(t *testing.T, tangle *Tangle) (map[string]wallet, map[ledgerstate.Address]wallet, *ledgerstate.Transaction) {
	tangle.EligibilityManager.Setup()

	// branches := make(map[string]ledgerstate.BranchID)

	wallets := make(map[string]wallet)
	walletsByAddress := make(map[ledgerstate.Address]wallet)
	w := createWallets(10)
	wallets["GENESIS"] = w[0]
	wallets["A"] = w[1]
	wallets["B"] = w[2]
	wallets["C"] = w[3]
	wallets["D"] = w[4]
	wallets["E"] = w[5]
	wallets["F"] = w[6]
	wallets["H"] = w[7]
	wallets["I"] = w[8]
	wallets["J"] = w[9]
	for _, wallet := range wallets {
		walletsByAddress[wallet.address] = wallet
	}
	genesisBalance := ledgerstate.NewColoredBalances(
		map[ledgerstate.Color]uint64{
			ledgerstate.ColorIOTA: 3,
		})
	genesisEssence := ledgerstate.NewTransactionEssence(
		0,
		time.Unix(DefaultGenesisTime, 0),
		identity.ID{},
		identity.ID{},
		ledgerstate.NewInputs(ledgerstate.NewUTXOInput(ledgerstate.NewOutputID(ledgerstate.GenesisTransactionID, 0))),
		ledgerstate.NewOutputs(ledgerstate.NewSigLockedColoredOutput(genesisBalance, wallets["GENESIS"].address)),
	)

	genesisTransaction := ledgerstate.NewTransaction(genesisEssence, ledgerstate.UnlockBlocks{ledgerstate.NewReferenceUnlockBlock(0)})
	stored, _, _ := tangle.LedgerState.UTXODAG.StoreTransaction(genesisTransaction)
	assert.True(t, stored, "genesis transaction stored")
	fmt.Println("genesis: ", genesisTransaction.ID().Base58())

	return wallets, walletsByAddress, genesisTransaction
}

func scenario1MessagesApproveEmptyID(t *testing.T, tangle *Tangle, genesisTransaction *ledgerstate.Transaction, wallets map[string]wallet, walletsByAddress map[ledgerstate.Address]wallet) (map[string]*Message, map[string]*ledgerstate.Transaction, map[string]*ledgerstate.UTXOInput, map[string]*ledgerstate.SigLockedSingleOutput, map[ledgerstate.OutputID]ledgerstate.Output) {
	messages := make(map[string]*Message)
	transactions := make(map[string]*ledgerstate.Transaction)
	inputs := make(map[string]*ledgerstate.UTXOInput)
	outputs := make(map[string]*ledgerstate.SigLockedSingleOutput)
	outputsByID := make(map[ledgerstate.OutputID]ledgerstate.Output)

	// Message 1
	inputs["GENESIS"] = ledgerstate.NewUTXOInput(ledgerstate.NewOutputID(genesisTransaction.ID(), 0))
	outputs["1A"] = ledgerstate.NewSigLockedSingleOutput(1, wallets["A"].address)
	outputs["1B"] = ledgerstate.NewSigLockedSingleOutput(1, wallets["B"].address)
	outputs["1C"] = ledgerstate.NewSigLockedSingleOutput(1, wallets["C"].address)

	transactions["1"] = makeTransaction(ledgerstate.NewInputs(inputs["GENESIS"]), ledgerstate.NewOutputs(outputs["1A"], outputs["1B"], outputs["1C"]), outputsByID, walletsByAddress, wallets["GENESIS"])
	messages["1"] = newTestParentsPayloadMessage(transactions["1"], []MessageID{EmptyMessageID}, []MessageID{})

	stored, _, _ := tangle.LedgerState.UTXODAG.StoreTransaction(transactions["1"])
	assert.True(t, stored, "transaction stored")
	tangle.Storage.StoreMessage(messages["1"])

	//_, stored = tangle.Storage.StoreAttachment(transactions["1"].ID(), messages["1"].ID())
	//assert.True(t, stored, "attachment stored")
	// Message 2
	inputs["2A"] = ledgerstate.NewUTXOInput(ledgerstate.NewOutputID(transactions["1"].ID(), selectIndex(transactions["1"], wallets["A"])))
	outputsByID[inputs["2A"].ReferencedOutputID()] = ledgerstate.NewOutputs(outputs["1A"])[0]

	outputs["2D"] = ledgerstate.NewSigLockedSingleOutput(1, wallets["D"].address)
	transactions["2"] = makeTransaction(ledgerstate.NewInputs(inputs["2A"]), ledgerstate.NewOutputs(outputs["2D"]), outputsByID, walletsByAddress)
	messages["2"] = newTestParentsPayloadMessage(transactions["2"], []MessageID{EmptyMessageID}, []MessageID{})

	tangle.Storage.StoreMessage(messages["2"])
	stored, _, _ = tangle.LedgerState.UTXODAG.StoreTransaction(transactions["2"])
	assert.True(t, stored, "transaction stored")
	//_, stored = tangle.Storage.StoreAttachment(transactions["2"].ID(), messages["2"].ID())
	//assert.True(t, stored, "attachment stored")

	fmt.Println("txs: ", transactions["1"].ID(), transactions["2"].ID())
	fmt.Println("msgs: ", messages["1"].ID(), messages["2"].ID())

	return messages, transactions, inputs, outputs, outputsByID
}

func scenario2MessagesApproveDependency(t *testing.T, tangle *Tangle, genesisTransaction *ledgerstate.Transaction, wallets map[string]wallet, walletsByAddress map[ledgerstate.Address]wallet) (map[string]*Message, map[string]*ledgerstate.Transaction, map[string]*ledgerstate.UTXOInput, map[string]*ledgerstate.SigLockedSingleOutput, map[ledgerstate.OutputID]ledgerstate.Output) {
	messages := make(map[string]*Message)
	transactions := make(map[string]*ledgerstate.Transaction)
	inputs := make(map[string]*ledgerstate.UTXOInput)
	outputs := make(map[string]*ledgerstate.SigLockedSingleOutput)
	outputsByID := make(map[ledgerstate.OutputID]ledgerstate.Output)

	// TODO
	return messages, transactions, inputs, outputs, outputsByID
}

//type utxoDagMock struct {
//	mock.Mock
//	test *testing.T
//	utxoDag *ledgerstate.UTXODAG
//}
//
//func newUtxoDagMock(t *testing.T, utxoDag *ledgerstate.UTXODAG) *utxoDagMock {
//	u := &utxoDagMock{
//		test: t,
//	}
//	u.Test(t)
//	u.utxoDag = utxoDag
//	return u
//}
//
//func (u *utxoDagMock) InclusionState(transactionID ledgerstate.TransactionID) (inclusionState ledgerstate.InclusionState, err error) {
//	args := u.Called(transactionID)
//	inclusionState = args.Get(0).(ledgerstate.InclusionState)
//	return
//}
//
//func (u *utxoDagMock) Shutdown() {
//	u.utxoDag.Shutdown()
//	return
//}
//
//func (u *utxoDagMock) StoreTransaction(transaction *ledgerstate.Transaction) (stored bool, solidityType ledgerstate.SolidityType, err error) {
//	return u.utxoDag.StoreTransaction(transaction)
//
//}
//
//func (u *utxoDagMock) CheckTransaction(transaction *ledgerstate.Transaction) (err error) {
//	return u.utxoDag.CheckTransaction(transaction)
//}
//
//func (u *utxoDagMock) CachedTransaction(transactionID ledgerstate.TransactionID) (cachedTransaction *ledgerstate.CachedTransaction) {
//	return u.utxoDag.CachedTransaction(transactionID)
//}
//
//func (u *utxoDagMock) Transaction(transactionID ledgerstate.TransactionID) (transaction *ledgerstate.Transaction) {
//	return u.utxoDag.Transaction(transactionID)
//}
//
//func (u *utxoDagMock) Transactions() (transactions map[ledgerstate.TransactionID]*ledgerstate.Transaction) {
//	return u.utxoDag.Transactions()
//}
//
//func (u *utxoDagMock) CachedTransactionMetadata(transactionID ledgerstate.TransactionID, computeIfAbsentCallback ...func(transactionID ledgerstate.TransactionID) *ledgerstate.TransactionMetadata) (cachedTransactionMetadata *CachedTransactionMetadata) {
//	return u.utxoDag.CachedTransactionMetadata(transactionID, computeIfAbsentCallback...)
//}
//
//func (u *utxoDagMock) CachedOutput(outputID ledgerstate.OutputID) (cachedOutput *ledgerstate.CachedOutput) {
//	return u.utxoDag.CachedOutput(outputID)
//}
//
//func (u *utxoDagMock) CachedOutputMetadata(outputID ledgerstate.OutputID) (cachedOutput *ledgerstate.CachedOutputMetadata) {
//	return u.utxoDag.CachedOutputMetadata(outputID)
//}
//
//func (u *utxoDagMock) CachedConsumers(outputID ledgerstate.OutputID, optionalSolidityType ...ledgerstate.SolidityType) (cachedConsumers ledgerstate.CachedConsumers) {
//	return u.utxoDag.CachedConsumers(outputID, optionalSolidityType...)
//}
//
//func (u *utxoDagMock) LoadSnapshot(snapshot *ledgerstate.Snapshot) {
//	u.utxoDag.LoadSnapshot(snapshot)
//}
//
//func (u *utxoDagMock) CachedAddressOutputMapping(address ledgerstate.Address) (cachedAddressOutputMappings ledgerstate.CachedAddressOutputMappings) {
//	return u.utxoDag.CachedAddressOutputMapping(address)
//}
//
//func (u *utxoDagMock) SetTransactionConfirmed(transactionID ledgerstate.TransactionID) (err error) {
//	return u.utxoDag.SetTransactionConfirmed(transactionID)
//}
