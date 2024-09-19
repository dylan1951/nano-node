package ledger

import (
	"errors"
	"fmt"
	"github.com/accept-nano/ed25519-blake2b"
	"github.com/cockroachdb/pebble"
	"log"
	"node/blocks"
	"node/config"
	"node/store"
	"node/types"
	"node/types/uint128"
	"node/utils"
)

type Account struct {
	PublicKey types.PublicKey
	Frontier  types.Hash
	Height    uint64
	Balance   uint128.Uint128
	Version   uint8
}

func AccountFromPublicKey(publicKey types.PublicKey) *Account {
	if record := store.GetAccount(publicKey); record != nil {
		return &Account{
			PublicKey: publicKey,
			Frontier:  record.Frontier,
			Height:    record.Height,
			Balance:   record.Balance,
			Version:   record.Version,
		}
	}
	return &Account{
		PublicKey: publicKey,
	}
}

func GetUnsyncedAccount() *Account {
	publicKey, err := store.GetUnsyncedAccount()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", publicKey.GoString())
	return AccountFromPublicKey(*publicKey)
}

var Invalid = errors.New("invalid block")
var Fork = errors.New("fork")
var Old = errors.New("old")

type MissingDependency struct {
	Dependency types.Hash
}

func (e MissingDependency) Error() string {
	return "dependency missing"
}

func (account *Account) AddBlock(batch *pebble.Batch, block blocks.Block) error {
	hash := block.Hash()

	if store.GetBlock(hash) != nil {
		fmt.Printf("already processed %s\n", hash.GoString())
		return Old
	}

	if block.GetPrevious() != account.Frontier {
		if record := store.GetBlock(block.GetPrevious()); record == nil {
			fmt.Printf("block %s is missing previous block: %s\n", block.Hash().GoString(), block.GetPrevious().GoString())
			return &MissingDependency{Dependency: block.GetPrevious()}
		}
		log.Fatalf("previous %s doesn't match frontier: %s\n", block.GetPrevious().GoString(), account.Frontier.GoString())
		return Fork
	}

	newBlockRecord := store.BlockRecord{
		Block:   block,
		Account: account.PublicKey,
	}

	signer := account.PublicKey[:]

	if block, ok := block.(*blocks.StateBlock); ok && account.Balance == block.Balance {
		switch block.Link {
		case config.EpochV1:
			signer = config.Network.Genesis.Account[:]
		case config.EpochV2:
			signer = config.Network.EpochV2Signer[:]
		}
	}

	if !ed25519.Verify(signer, hash[:], block.Common().Signature[:]) {
		fmt.Printf("signer: %s\n", utils.PubKeyToAddress(signer, false))
		block.Print()
		log.Fatalf("invalid signature")
		return Invalid
	}

	switch block := block.(type) {
	case *blocks.OpenBlock:
		if err := account.Receive(block.Source); err != nil {
			return err
		}

	case *blocks.SendBlock:
		if block.Balance.Cmp(account.Balance) >= 0 {
			log.Fatalf("sendblock invalid balance")
			return Invalid
		}

		newBlockRecord.Receivable = account.Balance.Sub(block.Balance)
		account.Balance = block.Balance

		if !block.Destination.IsZero() {
			// ensure not burn address
			store.MarkAccountUnsynced(batch, block.Destination)
		}

	case *blocks.ReceiveBlock:
		if err := account.Receive(block.Source); err != nil {
			return err
		}

	case *blocks.StateBlock:
		isReceiving := block.Balance.Cmp(account.Balance) > 0
		isSending := block.Balance.Cmp(account.Balance) < 0

		switch {
		case isReceiving:
			if err := account.Receive(block.Link); err != nil {
				return err
			}

			if account.Balance != block.Balance {
				return Invalid
			}
		case isSending:
			newBlockRecord.Receivable = account.Balance.Sub(block.Balance)
			account.Balance = block.Balance
			if !block.Link.IsZero() {
				// ensure not burn address
				store.MarkAccountUnsynced(batch, block.Link)
			}
		default:
			isRepUnchanged := block.Representative == account.Representative()
			switch {
			case block.Link == config.EpochV1 && account.Version == 0 && isRepUnchanged:
				account.Version = 1
			case block.Link == config.EpochV2 && account.Version == 1 && isRepUnchanged:
				account.Version = 2
			case !block.Link.IsZero():
				log.Fatalf("link not zero")
				return Invalid
			}
		}
	}

	account.Frontier = hash
	account.Height++

	store.PutBlock(batch, hash, newBlockRecord)
	store.SetAccount(batch, account.PublicKey, store.AccountRecord{
		Frontier: account.Frontier,
		Height:   account.Height,
		Balance:  account.Balance,
		Version:  account.Version,
	})

	return nil
}

func (account *Account) Receive(sourceHash types.Hash) error {
	source := store.GetBlock(sourceHash)

	if source == nil {
		fmt.Printf("account %s is missing souce: %s\n", account.PublicKey.GoString(), sourceHash.GoString())
		return &MissingDependency{Dependency: sourceHash}
	}

	if source.Receivable.IsZero() {
		return Fork
	}

	switch source := source.Block.(type) {
	case *blocks.SendBlock:
		if source.Destination != account.PublicKey {
			log.Fatalf("destination mismatch")
			return Invalid
		}
	case *blocks.StateBlock:
		if source.Link != account.PublicKey {
			log.Fatalf("link mismatch")
			return Invalid
		}
	default:
		log.Fatalf("can't receive from this kind of block")
		return Invalid
	}

	account.Balance = account.Balance.Add(source.Receivable)

	return nil
}

func (account *Account) Representative() types.PublicKey {
	record := store.GetBlock(account.Frontier)

	if record == nil {
		return types.PublicKey{}
	}

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

func (account *Account) Address() string {
	return utils.PubKeyToAddress(account.PublicKey[:], false)
}
