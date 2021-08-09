# sensible-script-decoder

Locking script decoder for Sensible Contract

解码锁定脚本，获得Sensible相关的字段数据。目前支持识别4种脚本类型（`CodeType`）：

0. `NONE` 普通脚本
1. `FT` sensible FT合约脚本
2. `UNIQUE` sensible unique合约脚本
3. `NFT` sensible NFT合约脚本


# 使用方法

参见 [satoblock/task/serial/tx.go](https://github.com/sensible-contract/satoblock/blob/8138d70eeef8bb7c726b5482090dc9191cc53aa2/task/serial/tx.go#L114)

```golang
	import (
		scriptDecoder "github.com/sensible-contract/sensible-script-decoder"
	)

	d.ScriptType = scriptDecoder.GetLockingScriptType(d.Script)
	txo := scriptDecoder.ExtractPkScriptForTxo(d.Script, d.ScriptType)

	d.CodeType = txo.CodeType
	d.CodeHash = txo.CodeHash
	d.GenesisId = txo.GenesisId
	...
```
