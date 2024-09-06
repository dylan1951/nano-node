package ledger

import (
	"errors"
	"github.com/accept-nano/ed25519-blake2b"
	"node/blocks"
	"node/store"
	"node/types"
	"node/types/uint128"
)

type Account struct {
	PublicKey types.PublicKey
	Frontier  types.Hash
	Height    uint64
	Balance   uint128.Uint128
}

func AccountFromPublicKey(publicKey types.PublicKey) *Account {
	if record := store.GetAccount(publicKey); record != nil {
		return &Account{
			PublicKey: publicKey,
			Frontier:  record.Frontier,
			Height:    record.Height,
			Balance:   record.Balance,
		}
	}
	return &Account{
		PublicKey: publicKey,
	}
}

func PubKeyFromBlock(block blocks.Block) (*types.PublicKey, error) {
	switch block := block.(type) {
	case *blocks.StateBlock:
		return &block.Account, nil
	case *blocks.OpenBlock:
		return &block.Account, nil
	default:
		if previous := store.GetBlock(block.GetPrevious()); previous != nil {
			return &previous.Account, nil
		}
		return nil, &MissingDependency{Dependency: block.GetPrevious()}
	}
}

var Invalid = errors.New("invalid block")
var Fork = errors.New("fork")

type MissingDependency struct {
	Dependency types.Hash
}

func (e MissingDependency) Error() string {
	return "dependency missing"
}

func (account *Account) AddBlock(block blocks.Block) error {
	hash := block.Hash()

	// todo: verify work

	if !ed25519.Verify(account.PublicKey[:], hash[:], block.BlockCommon().Signature[:]) {
		return Invalid
	}

	if block.GetPrevious() != account.Frontier {
		if record := store.GetBlock(block.GetPrevious()); record == nil {
			return &MissingDependency{Dependency: block.GetPrevious()}
		}
		return Fork
	}

	newBlockRecord := store.BlockRecord{
		Block:   block,
		Account: account.PublicKey,
	}

	switch block := block.(type) {
	case *blocks.OpenBlock:
		if err := account.Receive(block.Source); err != nil {
			return err
		}

	case *blocks.SendBlock:
		if block.Balance.Cmp(account.Balance) >= 0 {
			return Invalid
		}

		newBlockRecord.Receivable = account.Balance.Sub(block.Balance)
		account.Balance = block.Balance

	case *blocks.ReceiveBlock:
		if err := account.Receive(block.Source); err != nil {
			return err
		}

	case *blocks.StateBlock:
		isReceiving := block.Balance.Cmp(account.Balance) > 0
		isSending := block.Balance.Cmp(account.Balance) < 0

		if isReceiving {
			if err := account.Receive(block.Link); err != nil {
				return err
			}
		} else if isSending {
			newBlockRecord.Receivable = account.Balance.Sub(block.Balance)
			account.Balance = block.Balance
		} else {
			if !block.Link.IsZero() {
				return Invalid
			}
		}
	}

	store.PutBlock(hash, newBlockRecord)
	store.SetAccount(account.PublicKey, store.AccountRecord{
		Frontier: hash,
		Height:   account.Height + 1,
		Balance:  account.Balance,
	})

	return nil
}

func (account *Account) Receive(sourceHash types.Hash) error {
	source := store.GetBlock(sourceHash)

	if source == nil {
		return &MissingDependency{Dependency: sourceHash}
	}

	if source.Receivable.IsZero() {
		return Fork
	}

	switch source := source.Block.(type) {
	case *blocks.SendBlock:
		if source.Destination != account.PublicKey {
			return Invalid
		}
	case *blocks.StateBlock:
		if source.Link != account.PublicKey {
			return Invalid
		}
	default:
		return Invalid
	}

	account.Balance.Add(source.Receivable)

	return nil
}

func (account *Account) Representative() types.PublicKey {
	record := store.GetBlock(account.Frontier)

	for {
		switch block := record.Block.(type) {
		case *blocks.StateBlock:
			return block.Representative
		case *blocks.ChangeBlock:
			return block.Representative
		case *blocks.OpenBlock:
			return block.Representative
		case *blocks.ReceiveBlock:
			record = store.GetBlock(block.Previous)
		case *blocks.SendBlock:
			record = store.GetBlock(block.Previous)
		}
	}
}
