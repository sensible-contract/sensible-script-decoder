package script

import (
	"bytes"
	"encoding/binary"
)

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
		protoType := binary.LittleEndian.Uint32(scriptPk[protoTypeOffset : protoTypeOffset+4])

		switch protoType {
		case CodeType_FT:
			ret = decodeFT(scriptLen, scriptPk, txo)

		case CodeType_UNIQUE:
			ret = decodeUniqueV2(scriptLen, scriptPk, txo)
			if !ret {
				ret = decodeUnique(scriptLen, scriptPk, txo)
			}

		case CodeType_NFT:
			ret = decodeNFT(scriptLen, scriptPk, txo)

		case CodeType_NFT_SELL:
			ret = decodeNFTSell(scriptLen, scriptPk, txo)

		default:
			ret = false
		}
		return ret
	}

	if scriptPk[scriptLen-1] < 2 &&
		scriptPk[scriptLen-37-1] == 37 &&
		scriptPk[scriptLen-37-1-40-1] == 40 &&
		scriptPk[scriptLen-37-1-40-1-1] == OP_RETURN {
		ret = decodeNFTIssue(scriptLen, scriptPk, txo)

	} else if scriptPk[scriptLen-1] == 1 &&
		scriptPk[scriptLen-61-1] == 61 &&
		scriptPk[scriptLen-61-1-40-1] == 40 &&
		scriptPk[scriptLen-61-1-40-1-1] == OP_RETURN {
		ret = decodeNFTTransfer(scriptLen, scriptPk, txo)
	}

	return ret
}

func ExtractPkScriptForTxo(pkScript, scriptType []byte) (txo *TxoData) {
	txo = &TxoData{
		CodeHash:   empty,
		GenesisId:  empty,
		SensibleId: empty,
		AddressPkh: empty,
		MetaTxId:   empty,
		CustomData: empty,
	}

	if len(pkScript) == 0 {
		return txo
	}

	if isPubkeyHash(scriptType) {
		txo.AddressPkh = make([]byte, 20)
		copy(txo.AddressPkh, pkScript[3:23])
		return txo
	}

	if isPayToScriptHash(scriptType) {
		txo.AddressPkh = GetHash160(pkScript[2 : len(pkScript)-1])
		return txo
	}

	if isPubkey(scriptType) {
		txo.AddressPkh = GetHash160(pkScript[1 : len(pkScript)-1])
		return txo
	}

	// if isMultiSig(scriptType) {
	// 	return pkScript[:]
	// }

	if IsOpreturn(scriptType) {
		if hasSensibleFlag(pkScript) {
			txo.CodeType = CodeType_SENSIBLE
		}
		return txo
	}

	DecodeSensibleTxo(pkScript, txo)

	return txo
}

func GetLockingScriptType(pkScript []byte) (scriptType []byte) {
	length := len(pkScript)
	if length == 0 {
		return
	}
	scriptType = make([]byte, 0)

	lenType := 0
	p := uint(0)
	e := uint(length)

	for p < e && lenType < 32 {
		c := pkScript[p]
		if 0 < c && c < 0x4f {
			cnt, cntsize := SafeDecodeVarIntForScript(pkScript[p:])
			p += cnt + cntsize
		} else {
			p += 1
		}
		scriptType = append(scriptType, c)
		lenType += 1
	}
	return
}
