package script

import (
	"encoding/binary"
)

// nft Auction
func decodeNFTAuction(scriptLen int, pkScript []byte, txo *TxoData) bool {
	dataLen := 0
	if pkScript[scriptLen-232-1-1-1] == OP_RETURN &&
		pkScript[scriptLen-232-1-1] == 0x4c &&
		pkScript[scriptLen-232-1] == 232 {
		// nft auction v1
		// <nft auction data> = <rabinPubkeyHashArrayHash>(20bytes) + <timeRabinPubkeyHash>(20byte) + <endTimeStamp>(8byte)
		// + <nftID>(20byte) + <nftCodeHash>(20byte)
		// + <feeAmount(8byte)> + <feeAddress(20byte)>
		// + <startBsvPrice(8byte)> + <senderAddress(20byte)>
		// + <bidTimestamp>(8byte) + <bidBsvPrice>(8byte) + <bidderAddress(20byte)>
		// + <sensibleID(36 bytes)>
		// + <proto_version(4 bytes)> + <protoType(4 bytes)> + <'sensible'(8 bytes)>

		dataLen = 1 + 1 + 20 + 20 + 8 + 20 + 20 + 8 + 20 + 8 + 20 + 8 + 8 + 20 + 36 + 4 + 4 + 8 // 232 // 0x4c + pushdata + data
	} else {
		return false
	}

	protoVersionOffset := scriptLen - 8 - 4 - 4
	sensibleOffset := protoVersionOffset - 36
	bidderAddressOffset := sensibleOffset - 20
	bidBsvPriceOffset := bidderAddressOffset - 8
	bidTimestampOffset := bidBsvPriceOffset - 8

	senderAddressOffset := bidTimestampOffset - 20
	startBsvPriceOffset := senderAddressOffset - 8

	feeAddressOffset := startBsvPriceOffset - 20
	feeAmountOffset := feeAddressOffset - 8

	nftCodeHashOffset := feeAmountOffset - 20
	nftIdOffset := nftCodeHashOffset - 20

	endTimestampOffset := nftIdOffset - 8

	txo.CodeType = CodeType_NFT_AUCTION

	nft := &NFTAuctionData{
		FeeAmount:     binary.LittleEndian.Uint64(pkScript[feeAmountOffset : feeAmountOffset+8]),
		StartBsvPrice: binary.LittleEndian.Uint64(pkScript[startBsvPriceOffset : startBsvPriceOffset+8]),
		EndTimestamp:  binary.LittleEndian.Uint64(pkScript[endTimestampOffset : endTimestampOffset+8]),
		BidTimestamp:  binary.LittleEndian.Uint64(pkScript[bidTimestampOffset : bidTimestampOffset+8]),
		BidBsvPrice:   binary.LittleEndian.Uint64(pkScript[bidBsvPriceOffset : bidBsvPriceOffset+8]),
	}
	txo.NFTAuction = nft

	copy(nft.SensibleId[:], pkScript[sensibleOffset:sensibleOffset+36])
	copy(nft.NFTCodeHash[:], pkScript[nftCodeHashOffset:nftCodeHashOffset+20])
	copy(nft.NFTID[:], pkScript[nftIdOffset:nftIdOffset+20])
	copy(nft.FeeAddressPkh[:], pkScript[feeAddressOffset:feeAddressOffset+20])
	copy(nft.SenderAddressPkh[:], pkScript[senderAddressOffset:senderAddressOffset+20])
	copy(nft.BidderAddressPkh[:], pkScript[bidderAddressOffset:bidderAddressOffset+20])

	copy(txo.CodeHash[:], GetHash160(pkScript[:scriptLen-dataLen]))
	txo.GenesisIdLen = 20
	copy(txo.GenesisId[:], nft.NFTID[:])
	txo.HasAddress = true
	copy(txo.AddressPkh[:], nft.SenderAddressPkh[:]) // sender
	return true
}
