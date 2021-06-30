package script

func ExtractPkScriptForTxo(Pkscript, scriptType []byte) (txo *TxoData) {
	txo = &TxoData{
		CodeType:   CodeType_NONE,
		CodeHash:   empty,
		GenesisId:  empty,
		SensibleId: empty,
		AddressPkh: empty,

		MetaTxId: empty,
		TokenIdx: 0,

		Name:    "",
		Symbol:  "",
		Amount:  0,
		Decimal: 0,
	}

	if isPubkeyHash(scriptType) {
		txo.AddressPkh = make([]byte, 20)
		copy(txo.AddressPkh, Pkscript[3:23])
		return txo
	}

	if isPayToScriptHash(scriptType) {
		txo.AddressPkh = GetHash160(Pkscript[2 : len(Pkscript)-1])
		return txo
	}

	if isPubkey(scriptType) {
		txo.AddressPkh = GetHash160(Pkscript[1 : len(Pkscript)-1])
		return txo
	}

	// if isMultiSig(scriptType) {
	// 	return Pkscript[:]
	// }
	DecodeSensibleTxo(Pkscript, txo)

	return txo
}

func GetLockingScriptType(pkscript []byte) (scriptType []byte) {
	length := len(pkscript)
	if length == 0 {
		return
	}
	scriptType = make([]byte, 0)

	lenType := 0
	p := uint(0)
	e := uint(length)

	for p < e && lenType < 32 {
		c := pkscript[p]
		if 0 < c && c < 0x4f {
			cnt, cntsize := SafeDecodeVarIntForScript(pkscript[p:])
			p += cnt + cntsize
		} else {
			p += 1
		}
		scriptType = append(scriptType, c)
		lenType += 1
	}
	return
}

func IsLockingScriptOnlyEqual(pkscript []byte) bool {
	// test locking script
	// "0b 3c4b616e7965323032303e 87"

	length := len(pkscript)
	if length == 0 {
		return true
	}
	if pkscript[length-1] != 0x87 {
		return false
	}
	cnt, cntsize := SafeDecodeVarIntForScript(pkscript)
	if length == int(cnt+cntsize+1) {
		return true
	}
	return false
}
