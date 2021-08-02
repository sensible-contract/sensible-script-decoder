package script

import (
	"encoding/binary"
)

func decodeNFTIssue(scriptLen int, scriptPk []byte, txo *TxoData) bool {
	// nft issue
	txo.CodeType = CodeType_NFT
	genesisIdLen := 40
	genesisOffset := scriptLen - 37 - 1 - genesisIdLen
	tokenSupplyOffset := scriptLen - 1 - 8
	tokenIndexOffset := tokenSupplyOffset - 8
	addressOffset := tokenIndexOffset - 20

	dataLen := 1 + 1 + genesisIdLen + 1 + 37 // opreturn + pushdata + pushdata + data
	txo.CodeHash = GetHash160(scriptPk[:scriptLen-dataLen])
	txo.GenesisId = make([]byte, genesisIdLen)
	copy(txo.GenesisId, scriptPk[genesisOffset:genesisOffset+genesisIdLen])

	txo.TokenSupply = binary.LittleEndian.Uint64(scriptPk[tokenSupplyOffset : tokenSupplyOffset+8])
	txo.TokenIndex = binary.LittleEndian.Uint64(scriptPk[tokenIndexOffset : tokenIndexOffset+8])

	txo.TokenIndex = txo.TokenSupply

	txo.AddressPkh = make([]byte, 20)
	copy(txo.AddressPkh, scriptPk[addressOffset:addressOffset+20])
	return true
}

func decodeNFTTransfer(scriptLen int, scriptPk []byte, txo *TxoData) bool {
	// nft transfer
	txo.CodeType = CodeType_NFT
	genesisIdLen := 40
	genesisOffset := scriptLen - 61 - 1 - genesisIdLen
	metaTxIdOffset := scriptLen - 1 - 32
	tokenIndexOffset := metaTxIdOffset - 8
	addressOffset := tokenIndexOffset - 20

	dataLen := 1 + 1 + genesisIdLen + 1 + 61 // opreturn + pushdata + pushdata + data
	txo.CodeHash = GetHash160(scriptPk[:scriptLen-dataLen])
	txo.GenesisId = make([]byte, genesisIdLen)
	copy(txo.GenesisId, scriptPk[genesisOffset:genesisOffset+genesisIdLen])

	txo.MetaTxId = make([]byte, 32)
	copy(txo.MetaTxId, scriptPk[metaTxIdOffset:metaTxIdOffset+32])

	txo.TokenIndex = binary.LittleEndian.Uint64(scriptPk[tokenIndexOffset : tokenIndexOffset+8])

	txo.AddressPkh = make([]byte, 20)
	copy(txo.AddressPkh, scriptPk[addressOffset:addressOffset+20])
	return true
}
