package script

const (
	CodeType_NONE   uint32 = 0
	CodeType_FT     uint32 = 1
	CodeType_UNIQUE uint32 = 2
	CodeType_NFT    uint32 = 3
)

var empty = make([]byte, 1)

var CodeTypeName []string = []string{
	"NONE",
	"FT",
	"UNIQUE",
	"NFT",
}

type TxoData struct {
	CodeType   uint32
	CodeHash   []byte
	GenesisId  []byte // for search: codehash + genesis
	SensibleId []byte // GenesisTx outpoint
	AddressPkh []byte

	MetaTxId        []byte // nft metatxid
	MetaOutputIndex uint32
	TokenIndex      uint64 // nft tokenIndex
	TokenSupply     uint64 // nft tokenSupply

	Name    string // ft name
	Symbol  string // ft symbol
	Amount  uint64 // ft amount
	Decimal uint64 // ft decimal
}
