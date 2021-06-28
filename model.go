package script

type TxoData struct {
	IsNFT      bool
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
