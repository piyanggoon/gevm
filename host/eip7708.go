package host

import (
	"github.com/Giulio2002/gevm/spec"
	"github.com/Giulio2002/gevm/state"
	"github.com/Giulio2002/gevm/types"
)

var (
	eip7708TransferTopic = types.B256{
		0xdd, 0xf2, 0x52, 0xad, 0x1b, 0xe2, 0xc8, 0x9b,
		0x69, 0xc2, 0xb0, 0x68, 0xfc, 0x37, 0x8d, 0xaa,
		0x95, 0x2b, 0xa7, 0xf1, 0x63, 0xc4, 0xa1, 0x16,
		0x28, 0xf5, 0x5a, 0x4d, 0xf5, 0x23, 0xb3, 0xef,
	}
	eip7708BurnTopic = types.B256{
		0xcc, 0x16, 0xf5, 0xdb, 0xb4, 0x87, 0x32, 0x80,
		0x81, 0x5c, 0x1e, 0xe0, 0x9d, 0xbd, 0x06, 0x73,
		0x6c, 0xff, 0xcc, 0x18, 0x44, 0x12, 0xcf, 0x7a,
		0x71, 0xa0, 0xfd, 0xb7, 0x5d, 0x39, 0x7c, 0xa5,
	}
	eip7708SystemAddress = types.Address{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfe,
	}
)

func appendEIP7708TransferLog(j *state.Journal, from, to types.Address, amount types.Uint256) {
	if j == nil || !j.Cfg.Spec.IsEnabledIn(spec.Amsterdam) || amount == types.U256Zero || from == to {
		return
	}
	log := j.AllocLog()
	log.Address = eip7708SystemAddress
	log.Topics[0] = eip7708TransferTopic
	log.Topics[1] = addressTopic(from)
	log.Topics[2] = addressTopic(to)
	log.NumTopics = 3
	amount32 := amount.ToBytes32()
	log.Data = append(types.Bytes(nil), amount32[:]...)
}

func appendEIP7708BurnLog(j *state.Journal, from types.Address, amount types.Uint256) {
	if j == nil || !j.Cfg.Spec.IsEnabledIn(spec.Amsterdam) || amount == types.U256Zero {
		return
	}
	log := j.AllocLog()
	log.Address = eip7708SystemAddress
	log.Topics[0] = eip7708BurnTopic
	log.Topics[1] = addressTopic(from)
	log.NumTopics = 2
	amount32 := amount.ToBytes32()
	log.Data = append(types.Bytes(nil), amount32[:]...)
}

func addressTopic(addr types.Address) types.B256 {
	var topic types.B256
	copy(topic[12:], addr[:])
	return topic
}
