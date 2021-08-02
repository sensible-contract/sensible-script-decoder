package script

import (
	"encoding/binary"
)

// nft sell
func decodeNFTSell(scriptLen int, scriptPk []byte, txo *TxoData) bool {
	// dataLen := 0
	if scriptPk[scriptLen-112-1-1-1] == OP_RETURN &&
		scriptPk[scriptLen-112-1-1] == 0x4c &&
		scriptPk[scriptLen-112-1] == 112 {
		// nft sell v1
		// NFT 的codehash + NFT的genesis + NFT index + 卖家地址 + 价格 + nftID + 1 + 0x00010001 即 65537
		// <nft sell data> = <codehash(20 bytes)> + <genesis(20 bytes)> + <tokenIndex(8 bytes)> + <sellerAddress(20 bytes)> + <satoshisPrice(8 bytes)> + <nftID(20 bytes)> + <proto_version(4 bytes)> + <proto_type(4 bytes)> + <'sensible'(8 bytes)>

		// dataLen = 1 + 1 + 20 + 20 + 8 + 20 + 8 + 20 + 4 + 4 + 8 // 0x4c + pushdata + data
	} else {
		return false
	}

	protoVersionOffset := scriptLen - 8 - 4 - 4
	nftIdOffset := protoVersionOffset - 20
	priceOffset := nftIdOffset - 8
	addressOffset := priceOffset - 20
	tokenIndexOffset := addressOffset - 8
	genesisOffset := tokenIndexOffset - 20
	codehashOffset := genesisOffset - 20

	txo.CodeType = CodeType_NFT_SELL

	// txo.CodeHash = GetHash160(scriptPk[:scriptLen-dataLen])
	txo.CodeHash = make([]byte, 20)
	copy(txo.CodeHash, scriptPk[codehashOffset:codehashOffset+20])

	txo.GenesisId = make([]byte, 20)
	copy(txo.GenesisId, scriptPk[genesisOffset:genesisOffset+20])

	txo.TokenIndex = binary.LittleEndian.Uint64(scriptPk[tokenIndexOffset : tokenIndexOffset+8])

	txo.AddressPkh = make([]byte, 20)
	copy(txo.AddressPkh, scriptPk[addressOffset:addressOffset+20])

	// price
	txo.Amount = binary.LittleEndian.Uint64(scriptPk[priceOffset : priceOffset+8])
	return true
}
