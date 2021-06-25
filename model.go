package script

type TxoData struct {
	IsNFT      bool
	CodeHash   []byte
	GenesisId  []byte
	AddressPkh []byte
	MetaTxId   []byte
	Name       string // ft name
	Symbol     string // ft symbol
	DataValue  uint64 // ft amount / nft tokenIdx
	Decimal    uint64 // ft decimal
}
