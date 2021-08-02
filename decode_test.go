package script

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"strings"
	"testing"
)

type TxoDataTest struct {
	CodeType      uint32
	CodeHashHex   string
	GenesisIdHex  string // for search: codehash + genesis
	SensibleIdHex string // GenesisTx outpoint
	AddressPkhHex string

	MetaTxIdHex     string // nft metatxid
	MetaOutputIndex uint32
	TokenIndex      uint64 // nft tokenIndex
	TokenSupply     uint64 // nft tokenSupply

	Name    string // ft name
	Symbol  string // ft symbol
	Amount  uint64 // ft amount
	Decimal uint64 // ft decimal
}

var scripts []string

func init() {
	dat, err := ioutil.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}
	scripts = strings.Split(string(dat), "\n")
}

func TestDecode(t *testing.T) {
	for line, scriptHex := range scripts {
		if len(scriptHex) == 0 {
			continue
		}
		script, err := hex.DecodeString(scriptHex)
		if err != nil {
			t.Logf("ignore line: %d", line)
			continue
		}

		txo := &TxoData{}

		DecodeSensibleTxo(script, txo)

		data, _ := json.Marshal(TxoDataTest{
			CodeType:      txo.CodeType,
			CodeHashHex:   hex.EncodeToString(txo.CodeHash),
			GenesisIdHex:  hex.EncodeToString(txo.GenesisId),
			SensibleIdHex: hex.EncodeToString(txo.SensibleId),
			AddressPkhHex: hex.EncodeToString(txo.AddressPkh),

			MetaTxIdHex:     hex.EncodeToString(txo.MetaTxId),
			MetaOutputIndex: txo.MetaOutputIndex,
			TokenIndex:      txo.TokenIndex,
			TokenSupply:     txo.TokenSupply,

			Name:    txo.Name,
			Symbol:  txo.Symbol,
			Amount:  txo.Amount,
			Decimal: txo.Decimal,
		})
		t.Logf("scriptLen: %d, txo: %s", len(script), strings.ReplaceAll(string(data), ",", "\n"))
	}
}
