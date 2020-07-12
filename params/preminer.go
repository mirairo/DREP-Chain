package params

import (
	"github.com/drep-project/DREP-Chain/crypto"
	"math/big"
)

const (
	HoleAddressStr = "0x0000000000000000000000000000000000000000"

	DefaultGenesisParamMainnet = `{"Preminer":[{"Addr":"0xfb6711cbafbd5e75c612d69db9025f7eb5096d46","Value":10000000000000000000000000000}],"Miners":[{"Pubkey":"0x0305bfc35e079ca7ae7ad6fcf246f8ed0247d0fdbda0b7daa10f2b2f3a88e9fd8a","Node":"enode://f57881c48aaccf97485c2b65b421bfeda22cc3b427c44be7607b122fc1688abb@172.104.123.143:10086"},{"Pubkey":"0x0354d5d560039693ec4dedb416354bdaa7aa70823fd5fabf56ceeb76386efc5670","Node":"enode://9d25d161ae4b676e2df55accca93c3137df3166326d04420ffbdf66e887bd494@172.104.116.219:10086"},{"Pubkey":"0x031283a05887e8b291fd909a6302cdd84b0b50bbc0ae5da8ce186435a1c6b8d10a","Node":"enode://bc7ca1b57175f2d5c85da73d367408529468a034b97d083aaecf88196090e245@172.105.103.59:10086"},{"Pubkey":"0x0349028fc4ce0ed5b6d496fcbd12755779954c1a17e74dfb283190e14ee5ea3170","Node":"enode://0ebd0422ca32d70292be128342f9e5ca32ab3cef28dc32cc332169e578e7b4f5@109.74.203.199:10086"},{"Pubkey":"0x03902b01673cc0ffad466f68a7a1494c9da1a80fb63dae8c9d18dfc8f9aeff1eba","Node":"enode://c01ff36a9914f781a058c77e98f09eedd5ad7e2e575c7e6c2d6cc86d859693f2@139.162.161.8:10086"},{"Pubkey":"0x028b2388f1cd0b8056a3ef41a7559952a0abe40ce8d93f586c7aa23b02f841ed23","Node":"enode://bb7925d62126c1058f7ad951e862aa4fc0aaaa4dc7cae2fcac3b86aee316bc02@172.104.178.138:10086"},{"Pubkey":"0x0284454961eb9c389bf292e1d9a1eb957e4b22249771ce942a9eda1768fe607b32","Node":"enode://9863364f265843bbb2d7810c9495e83b04cc541bedc8ae44ca63a2c5ec9d1b75@172.105.182.178:10086"},{"Pubkey":"0x02e42e247ae1f9a737aaa974f83dc8ca2e96c9e9a635e7e4d124df633a886d5e0e","Node":"enode://7523f570a35d58841440a792b51ae0608a1c7f387d7f0899c59720a5f851604a@172.105.42.46:10086"},{"Pubkey":"0x036b07d786bd104b83069049be0f918c9b1505f5a90ab645837794a20c63f2b6f1","Node":"enode://179838dc7bdb7285f82ec8e15515e09208a7d7452898dd4b978df88b694402d3@96.126.122.90:10086"},{"Pubkey":"0x0299e438ccbca1a492dfa04ab3189b687f6311fc99cd9f6112b23b4433e4c758ec","Node":"enode://ab789802b04057655247b51fb58f48643a7a4ae883a786ee4871b67c6e57250f@45.79.252.45:10086"},{"Pubkey":"0x025d404222e2f67a77963ae7cd846bc95174e0150d47c61665b5a280fbef989c64","Node":"enode://2a87c3c98b416f22d11e8952b36c220942461ccd69ecdbc54aef1b0d90238da4@45.56.95.177:10086"},{"Pubkey":"0x02ea32b18f6a37b23ebccca1ed6ae4f513da764f85f3a37234c18cad627e820dc9","Node":"enode://8cec1501ea4bb26b6404809e4e3eef53149b490ec7d378085dd7a24b9e6d1211@45.33.89.236:10086"},{"Pubkey":"0x0398586c5d012e677ef9eb74785de6aa69e38c13e17e45071ec91c9e32bdc64848","Node":"enode://5a764d497a55e42ea33a833770f63e2ba8f6200da824db32ef2bcc937af88c44@176.58.112.109:10086"},{"Pubkey":"0x03ce8088cd8ab8fa7cfd7ecd71015bf525d9b071b0090889e399bc021e10496770","Node":"enode://956f3f79214db8065298e8447f8aa14ca611c02bb886a1f1e0be92eb74ad8984@172.104.179.82:10086"},{"Pubkey":"0x022d06d1131b740ee6324d35c4096ce5dc54a90c0b962b795fd46fa371bc734e82","Node":"enode://8ecedb07fcb27642c1560e000c79fa80a4e185001a40030fb50e698ff008a98d@139.162.188.229:10086"},{"Pubkey":"0x0306c3c718e9f9fad21a63c27e1527c599ec54611e17d6b90fa9a81dc7b6588624","Node":"enode://f2055090bb60b42c91272ce00f671c0048c46e065d489910341c101d4bb94af3@23.239.16.64:10086"},{"Pubkey":"0x0373d68e621f5063875c9b4647cb7b990d96f59fb3769cbb892e462975937acb48","Node":"enode://772c23d1de1ac72451c17da61889b5c211127393916f93a8671c202aad2ef9df@45.56.124.251:10086"},{"Pubkey":"0x02b6d93dddb0cb702266f55b7982067d4b6f5587af0273b660db8ac590d405733f","Node":"enode://d8eb92d787384442e40476fa1d4c564facca3667298566b3950554ce4cc0036f@172.104.106.112:10086"},{"Pubkey":"0x02da50f6372568a6f8dc9d05a65659dc7d5db7db5c1c9bdc4369c70ceb19200b28","Node":"enode://6890bf120d755a498a301e9cbbcbabba0c4c6628133c35e5624eb85b5e1f09ae@172.105.175.246:10086"},{"Pubkey":"0x03d94c1a566bf73783dbe21bfb4fb7b3e673de16a94f5b616c32bd4d44a2bf30c9","Node":"enode://f0f207d9199a246871759b29424358bc23bb9d2e3730dc3b151c9f9abdddd4be@172.105.206.6:10086"},{"Pubkey":"0x03c1586860d81359a51655dfe03663beb4db51ec5373f7b39a525561498dfe31b6","Node":"enode://5969d83137afd2ea7ed06419009ea52f86213094b46c6444468880f45659fedd@139.162.30.122:10086"}]}`

	DefaultGenesisParamTestnet = `{
	      "Preminer": [
	            {
	                  "Addr": "0x7d17376a5a611c768970f7ce99fbe309450bff6f",
	                  "Value": 10000000000000000000000000000
	            }
	      ],
	      "Miners": [
	            {
	                  "Pubkey": "0x0328378210fd26ac195c4880b5cf8a68e5477d5f2f409e4526ed6b49681091a391",
	                  "Node": "enode://548c58daf6dc65d463c155027fce3a909d555683543d1dca34cff1d68868c54f@39.100.111.74:44444"
	            },
	            {
	                  "Pubkey": "0x03efe1cad6eb9e161a9d4809eb0d40e9d9392d70e877cc2b41cd7a7526628ee007",
	                  "Node": "enode://385c49f05a235115515d5581485be6cd66bbcaf2dbace93d641b5e4c87c20255@39.98.39.224:44444"
	            },
	            {
	                  "Pubkey": "0x031afda919527b8c55997e8a2c2cdf33fc025708a73085b4f7f5c500a6c68ddb08",
	                  "Node": "enode://9296c4f6e4ceaaea24d0416f49bf7624e920d1f71f7a51877a5d0ed156e35ac5@39.99.44.60:44444"
	            }
	      ]
	}`
)

var (
	HoleAddress = crypto.HexToAddress(HoleAddressStr)
)

func CoinFromNumer(number int64) *big.Int {
	return big.NewInt(0).Mul(big.NewInt(Coin), big.NewInt(number))
}
