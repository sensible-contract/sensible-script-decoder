package script

import (
	"encoding/hex"
	"encoding/json"
)

const (
	CodeType_NONE   uint32 = 0
	CodeType_FT     uint32 = 1
	CodeType_UNIQUE uint32 = 2
	CodeType_NFT    uint32 = 3

	CodeType_SENSIBLE uint32 = 65536
	CodeType_NFT_SELL uint32 = 65536 + 1
)

var CodeTypeName []string = []string{
	"NONE",
	"FT",
	"UNIQUE",
	"NFT",
}

// 64/84 bytes
type SwapData struct {
	// fetchTokenContractHash + lpTokenID + lpTokenScriptCodeHash + Token1Amount + Token2Amount + lpAmount
	Token1Amount uint64
	Token2Amount uint64
	LpAmount     uint64
}

type FTData struct {
	CodeHash   [20]byte
	GenesisId  []byte // for search: codehash + genesis
	SensibleId []byte // GenesisTx outpoint

	Name    string // ft name
	Symbol  string // ft symbol
	Amount  uint64 // ft amount
	Decimal uint8  // ft decimal
}

func (u *FTData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		CodeHash   string
		GenesisId  string // for search: codehash + genesis
		SensibleId string // GenesisTx outpoint
		Name       string // ft name
		Symbol     string // ft symbol
		Amount     uint64 // ft amount
		Decimal    uint8  // ft decimal
	}{
		CodeHash:   hex.EncodeToString(u.CodeHash[:]),
		GenesisId:  hex.EncodeToString(u.GenesisId[:]),
		SensibleId: hex.EncodeToString(u.SensibleId[:]),
		Name:       u.Name,
		Symbol:     u.Symbol,
		Amount:     u.Amount,
		Decimal:    u.Decimal,
	})
}

type NFTData struct {
	CodeHash   [20]byte
	GenesisId  []byte // for search: codehash + genesis
	SensibleId []byte // GenesisTx outpoint

	MetaTxId        [32]byte // nft metatxid
	MetaOutputIndex uint32
	TokenIndex      uint64 // nft tokenIndex
	TokenSupply     uint64 // nft tokenSupply
}

func (u *NFTData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		CodeHash        string
		GenesisId       string // for search: codehash + genesis
		SensibleId      string // GenesisTx outpoint
		MetaTxId        string // nft metatxid
		MetaOutputIndex uint32
		TokenIndex      uint64 // nft tokenIndex
		TokenSupply     uint64 // nft tokenSupply

	}{
		CodeHash:        hex.EncodeToString(u.CodeHash[:]),
		GenesisId:       hex.EncodeToString(u.GenesisId),
		SensibleId:      hex.EncodeToString(u.SensibleId),
		MetaTxId:        hex.EncodeToString(u.MetaTxId[:]),
		MetaOutputIndex: u.MetaOutputIndex,
		TokenIndex:      u.TokenIndex,
		TokenSupply:     u.TokenSupply,
	})
}

type NFTSellData struct {
	CodeHash   [20]byte
	GenesisId  [20]byte // for search: codehash + genesis
	TokenIndex uint64   // nft tokenIndex
	Price      uint64   // nft price
}

func (u *NFTSellData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		CodeHash   string
		GenesisId  string // for search: codehash + genesis
		TokenIndex uint64
		Price      uint64
	}{
		CodeHash:   hex.EncodeToString(u.CodeHash[:]),
		GenesisId:  hex.EncodeToString(u.GenesisId[:]),
		TokenIndex: u.TokenIndex,
		Price:      u.Price,
	})
}

type UniqueData struct {
	CodeHash   [20]byte
	GenesisId  []byte // for search: codehash + genesis
	SensibleId []byte // GenesisTx outpoint
	CustomData []byte // unique data
	Swap       *SwapData
}

func (u *UniqueData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		CodeHash   string
		GenesisId  string // for search: codehash + genesis
		SensibleId string // GenesisTx outpoint
		CustomData string // unique data
		Swap       *SwapData
	}{
		CodeHash:   hex.EncodeToString(u.CodeHash[:]),
		GenesisId:  hex.EncodeToString(u.GenesisId),
		SensibleId: hex.EncodeToString(u.SensibleId),
		CustomData: hex.EncodeToString(u.CustomData),
		Swap:       u.Swap,
	})
}

type TxoData struct {
	CodeType   uint32
	HasAddress bool
	AddressPkh [20]byte
	NFT        *NFTData
	FT         *FTData
	Uniq       *UniqueData
	NFTSell    *NFTSellData
}

func (u *TxoData) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		CodeType   uint32
		HasAddress bool
		AddressPkh string
		NFT        *NFTData
		FT         *FTData
		Uniq       *UniqueData
		NFTSell    *NFTSellData
	}{
		CodeType:   u.CodeType,
		HasAddress: u.HasAddress,
		AddressPkh: hex.EncodeToString(u.AddressPkh[:]),
		NFT:        u.NFT,
		FT:         u.FT,
		Uniq:       u.Uniq,
		NFTSell:    u.NFTSell,
	})
}
