package script

import (
	"encoding/binary"
)

func decodeUnique(scriptLen int, pkScript []byte, txo *TxoData) bool {
	// <unique data> = <unique custom data> + <custom data length(4 bytes)> + <genesisFlag(1 bytes)> + <rabinPubKeyHashArrayHash(20 bytes)> + <sensibleID(36 bytes)> + <protoType(4 bytes)> + <protoFlag(8 bytes)>
	genesisIdLen := 56 // ft unique
	sensibleIdLen := 36

	protoTypeOffset := scriptLen - 8 - 4
	sensibleOffset := protoTypeOffset - sensibleIdLen

	genesisOffset := protoTypeOffset - genesisIdLen
	customDataSizeOffset := genesisOffset - 1 - 4
	customDataSize := binary.LittleEndian.Uint32(pkScript[customDataSizeOffset : customDataSizeOffset+4])
	varint := getVarIntLen(int(customDataSize) + 17 + genesisIdLen)
	dataLen := 1 + 1 + varint + int(customDataSize) + 17 + genesisIdLen // opreturn + 0x.. + pushdata + data

	if dataLen >= scriptLen || pkScript[scriptLen-dataLen] != OP_RETURN {
		dataLen = 0
		return false
	}
	txo.CodeType = CodeType_UNIQUE
	txo.AddressPkh = make([]byte, 20)
	txo.CodeHash = GetHash160(pkScript[:scriptLen-dataLen])

	// GenesisId is tokenIdHash
	txo.GenesisId = GetHash160(pkScript[genesisOffset : genesisOffset+genesisIdLen])

	txo.SensibleId = make([]byte, sensibleIdLen)
	copy(txo.SensibleId, pkScript[sensibleOffset:sensibleOffset+sensibleIdLen])
	return true
}

func decodeUniqueV2(scriptLen int, pkScript []byte, txo *TxoData) bool {
	// <unique data> = <unique custom data> + <custom data length(4 bytes)> + <genesisFlag(1 bytes)> + <rabinPubKeyHashArrayHash(20 bytes)> + <sensibleID(36 bytes)> + <protoVersion(4 bytes)> + <protoType(4 bytes)> + <protoFlag(8 bytes)>
	protoVersionLen := 4
	genesisIdLen := 56 // ft unique
	sensibleIdLen := 36

	protoTypeOffset := scriptLen - 8 - 4
	sensibleOffset := protoTypeOffset - protoVersionLen - sensibleIdLen

	genesisOffset := protoTypeOffset - protoVersionLen - genesisIdLen
	customDataSizeOffset := genesisOffset - 1 - 4
	customDataSize := int(binary.LittleEndian.Uint32(pkScript[customDataSizeOffset : customDataSizeOffset+4]))
	varint := getVarIntLen(customDataSize + 21 + genesisIdLen)
	dataLen := 1 + varint + customDataSize + 21 + genesisIdLen // 0x.. + pushdata + data

	if dataLen+1 >= scriptLen || pkScript[scriptLen-dataLen-1] != OP_RETURN {
		dataLen = 0
		return false
	}
	txo.CodeType = CodeType_UNIQUE
	txo.AddressPkh = make([]byte, 20)
	txo.CodeHash = GetHash160(pkScript[:scriptLen-dataLen])

	// GenesisId is tokenIdHash
	txo.GenesisId = GetHash160(pkScript[genesisOffset : genesisOffset+genesisIdLen])

	txo.CustomData = make([]byte, customDataSize)
	copy(txo.CustomData, pkScript[customDataSizeOffset-customDataSize:customDataSizeOffset])

	txo.SensibleId = make([]byte, sensibleIdLen)
	copy(txo.SensibleId, pkScript[sensibleOffset:sensibleOffset+sensibleIdLen])
	return true
}
