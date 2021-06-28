package script

const (
	CodeType_NONE   int = 0
	CodeType_FT     int = 1
	CodeType_UNIQUE int = 2
	CodeType_NFT    int = 3
)

type TxoData struct {
	CodeType   int
	CodeHash   []byte
	GenesisId  []byte
	AddressPkh []byte

	MetaTxId []byte // nft metatxid
	TokenIdx uint64 // nft tokenIdx

	Name    string // ft name
	Symbol  string // ft symbol
	Amount  uint64 // ft amount
	Decimal uint64 // ft decimal
}
