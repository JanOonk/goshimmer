package tangle

import (
	"testing"
	"time"

	"github.com/iotaledger/goshimmer/packages/ledgerstate"

	"github.com/iotaledger/hive.go/identity"
)

func TestDependenciesConfirmed(t *testing.T) {
	tangle := newTestTangle()
	defer tangle.Shutdown()
	messages, wallets, walletsByAddress, genesisTransaction := setupEligibilityTests(tangle)

	transactions := make(map[string]*ledgerstate.Transaction)
	inputs := make(map[string]*ledgerstate.UTXOInput)
	outputs := make(map[string]*ledgerstate.SigLockedSingleOutput)
	outputsByID := make(map[ledgerstate.OutputID]ledgerstate.Output)

	scenarioMessagesApproveEmptyID(tangle, inputs, genesisTransaction, outputs, wallets, transactions, outputsByID, walletsByAddress, messages)
	mockUTXO := newUtxoDagMock(t, tangle.LedgerState.UTXODAG)
	mockUTXO.On("InclusionState", transactions["1"]).Return(ledgerstate.Confirmed)
	tangle.EligibilityManager.Events.MessageEligible.Attach(events.NewClosure(func() {
		fmt.Println("Test ")
	}))
	tangle.Solidifier.Events.MessageSolid.Trigger(messages["1"])

}

func setupEligibilityTests(tangle *Tangle) (map[string]*Message, map[string]wallet, map[ledgerstate.Address]wallet, *ledgerstate.Transaction) {
	tangle.EligibilityManager.Setup()

	messages := make(map[string]*Message)
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
	return messages, wallets, walletsByAddress, genesisTransaction
}

func scenarioMessagesApproveEmptyID(tangle *Tangle, inputs map[string]*ledgerstate.UTXOInput, genesisTransaction *ledgerstate.Transaction, outputs map[string]*ledgerstate.SigLockedSingleOutput, wallets map[string]wallet, transactions map[string]*ledgerstate.Transaction, outputsByID map[ledgerstate.OutputID]ledgerstate.Output, walletsByAddress map[ledgerstate.Address]wallet, messages map[string]*Message) {
	// Message 1
	inputs["GENESIS"] = ledgerstate.NewUTXOInput(ledgerstate.NewOutputID(genesisTransaction.ID(), 0))
	outputs["A"] = ledgerstate.NewSigLockedSingleOutput(1, wallets["A"].address)
	outputs["B"] = ledgerstate.NewSigLockedSingleOutput(1, wallets["B"].address)
	outputs["C"] = ledgerstate.NewSigLockedSingleOutput(1, wallets["C"].address)
	transactions["1"] = makeTransaction(ledgerstate.NewInputs(inputs["GENESIS"]), ledgerstate.NewOutputs(outputs["A"], outputs["B"], outputs["C"]), outputsByID, walletsByAddress, wallets["GENESIS"])
	messages["1"] = newTestParentsPayloadMessage(transactions["1"], []MessageID{EmptyMessageID}, []MessageID{})
	tangle.Storage.StoreMessage(messages["1"])

	// Message 2
	inputs["B"] = ledgerstate.NewUTXOInput(ledgerstate.NewOutputID(transactions["1"].ID(), selectIndex(transactions["1"], wallets["B"])))
	outputsByID[inputs["B"].ReferencedOutputID()] = ledgerstate.NewOutputs(outputs["B"])[0]

	inputs["Message 1"] = ledgerstate.NewUTXOInput(ledgerstate.NewOutputID(transactions["1"].ID(), 0))

	outputs["D"] = ledgerstate.NewSigLockedSingleOutput(1, wallets["D"].address)
	transactions["2"] = makeTransaction(ledgerstate.NewInputs(inputs["Message1"]), ledgerstate.NewOutputs(outputs["D"]), outputsByID, walletsByAddress, wallets["A"])
	messages["2"] = newTestParentsPayloadMessage(transactions["2"], []MessageID{EmptyMessageID}, []MessageID{})
	tangle.Storage.StoreMessage(messages["2"])
	return

}

type utxoDagMock struct {
	mock.Mock
	test *testing.T
}

func newUtxoDagMock(t *testing.T, utxoDag *ledgerstate.UTXODAG) *utxoDagMock {
	u := &utxoDagMock{
		test: t,
	}
	u.Test(t)

	return u
}

func (u *utxoDagMock) InclusionState(transactionID ledgerstate.TransactionID) (inclusionState ledgerstate.InclusionState, err error) {
	args := u.Called(transactionID)
	inclusionState = args.Get(0).(ledgerstate.InclusionState)
	return
}
