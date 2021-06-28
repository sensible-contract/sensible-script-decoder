package script

import (
	"bytes"
	"encoding/binary"
)

var empty = make([]byte, 1)

func hasSensibleFlag(scriptPk []byte) bool {
	return bytes.HasSuffix(scriptPk, []byte("sensible")) || bytes.HasSuffix(scriptPk, []byte("oraclesv"))
}

func DecodeSensibleTxo(scriptPk []byte, txo *TxoData) bool {
	scriptLen := len(scriptPk)
	if scriptLen < 1024 {
		return false
	}

	ret := false
	if hasSensibleFlag(scriptPk) {
		protoTypeOffset := scriptLen - 8 - 4
		if scriptPk[protoTypeOffset] == 1 { // PROTO_TYPE == 1
			ret = decodeFT(scriptLen, scriptPk, txo)
		} else if scriptPk[protoTypeOffset] == 2 { // PROTO_TYPE == 2
			ret = decodeUnique(scriptLen, scriptPk, txo)
		}
	} else if scriptPk[scriptLen-1] < 2 && scriptPk[scriptLen-37-1] == 37 && scriptPk[scriptLen-37-1-40-1] == 40 && scriptPk[scriptLen-37-1-40-1-1] == OP_RETURN {
		ret = decodeNFTIssue(scriptLen, scriptPk, txo)
	} else if scriptPk[scriptLen-1] == 1 && scriptPk[scriptLen-61-1] == 61 && scriptPk[scriptLen-61-1-40-1] == 40 && scriptPk[scriptLen-61-1-40-1-1] == OP_RETURN {
		ret = decodeNFTTransfer(scriptLen, scriptPk, txo)
	}

	return ret
}

func decodeFT(scriptLen int, scriptPk []byte, txo *TxoData) bool {
	dataLen := 0
	genesisIdLen := 0
	if scriptPk[scriptLen-72-36-1-1] == 0x4c && scriptPk[scriptLen-72-36-1] == 108 {
		genesisIdLen = 36                       // new ft
		dataLen = 1 + 1 + 1 + 72 + genesisIdLen // opreturn + 0x4c + pushdata + data
	} else if scriptPk[scriptLen-72-20-1-1] == 0x4c && scriptPk[scriptLen-72-20-1] == 92 {
		genesisIdLen = 20                       // old ft
		dataLen = 1 + 1 + 1 + 72 + genesisIdLen // opreturn + 0x4c + pushdata + data
	} else if scriptPk[scriptLen-50-36-1-1] == 0x4c && scriptPk[scriptLen-50-36-1] == 86 {
		genesisIdLen = 36                       // old ft
		dataLen = 1 + 1 + 1 + 50 + genesisIdLen // opreturn + 0x4c + pushdata + data
	} else if scriptPk[scriptLen-92-20-1-1] == 0x4c && scriptPk[scriptLen-92-20-1] == 112 {
		genesisIdLen = 20                       // old ft
		dataLen = 1 + 1 + 1 + 92 + genesisIdLen // opreturn + 0x4c + pushdata + data
	} else {
		genesisIdLen = 40                       // error ft
		dataLen = 1 + 1 + 1 + 72 + genesisIdLen // opreturn + 0x4c + pushdata + data
		return false
	}

	protoTypeOffset := scriptLen - 8 - 4
	genesisOffset := protoTypeOffset - genesisIdLen
	amountOffset := genesisOffset - 8
	addressOffset := amountOffset - 20
	decimalOffset := addressOffset - 1
	symbolOffset := decimalOffset - 1 - 10
	nameOffset := symbolOffset - 20

	txo.CodeType = CodeType_FT
	txo.Decimal = uint64(scriptPk[decimalOffset])
	txo.Symbol = string(bytes.TrimRight(scriptPk[symbolOffset:symbolOffset+10], "\x00"))
	txo.Name = string(bytes.TrimRight(scriptPk[nameOffset:nameOffset+20], "\x00"))

	txo.Amount = binary.LittleEndian.Uint64(scriptPk[amountOffset : amountOffset+8])

	txo.AddressPkh = make([]byte, 20)
	copy(txo.AddressPkh, scriptPk[addressOffset:addressOffset+20])

	txo.CodeHash = GetHash160(scriptPk[:scriptLen-dataLen])
	txo.GenesisId = make([]byte, genesisIdLen)
	copy(txo.GenesisId, scriptPk[genesisOffset:genesisOffset+genesisIdLen])
	return true
}

func decodeUnique(scriptLen int, scriptPk []byte, txo *TxoData) bool {
	txo.CodeType = CodeType_UNIQUE
	genesisIdLen := 36 // ft unique
	protoTypeOffset := scriptLen - 8 - 4
	genesisOffset := protoTypeOffset - genesisIdLen
	customDataSizeOffset := genesisOffset - 1 - 4

	customDataSize := binary.LittleEndian.Uint64(scriptPk[customDataSizeOffset : customDataSizeOffset+4])
	varint := getVarIntLen(customDataSize)
	dataLen := 1 + 1 + varint + int(customDataSize) + 17 + genesisIdLen // opreturn + 0x.. + pushdata + data

	if dataLen >= scriptLen || scriptPk[scriptLen-dataLen] != OP_RETURN {
		dataLen = 0
		return false
	}

	txo.AddressPkh = make([]byte, 20)
	txo.CodeHash = GetHash160(scriptPk[:scriptLen-dataLen])
	txo.GenesisId = make([]byte, genesisIdLen)
	copy(txo.GenesisId, scriptPk[genesisOffset:genesisOffset+genesisIdLen])

	return true
}

func decodeNFTIssue(scriptLen int, scriptPk []byte, txo *TxoData) bool {
	// nft issue
	txo.CodeType = CodeType_NFT
	genesisIdLen := 40
	genesisOffset := scriptLen - 37 - 1 - genesisIdLen
	tokenIdxOffset := scriptLen - 1 - 8
	addressOffset := tokenIdxOffset - 8 - 20

	dataLen := 1 + 1 + genesisIdLen + 1 + 37 // opreturn + pushdata + pushdata + data
	txo.CodeHash = GetHash160(scriptPk[:scriptLen-dataLen])
	txo.GenesisId = make([]byte, genesisIdLen)
	copy(txo.GenesisId, scriptPk[genesisOffset:genesisOffset+genesisIdLen])

	txo.TokenIdx = binary.LittleEndian.Uint64(scriptPk[tokenIdxOffset : tokenIdxOffset+8])

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
	tokenIdxOffset := metaTxIdOffset - 8
	addressOffset := tokenIdxOffset - 20

	dataLen := 1 + 1 + genesisIdLen + 1 + 61 // opreturn + pushdata + pushdata + data
	txo.CodeHash = GetHash160(scriptPk[:scriptLen-dataLen])
	txo.GenesisId = make([]byte, genesisIdLen)
	copy(txo.GenesisId, scriptPk[genesisOffset:genesisOffset+genesisIdLen])

	txo.MetaTxId = make([]byte, 32)
	copy(txo.MetaTxId, scriptPk[metaTxIdOffset:metaTxIdOffset+32])

	txo.TokenIdx = binary.LittleEndian.Uint64(scriptPk[tokenIdxOffset : tokenIdxOffset+8])

	txo.AddressPkh = make([]byte, 20)
	copy(txo.AddressPkh, scriptPk[addressOffset:addressOffset+20])

	return true
}
