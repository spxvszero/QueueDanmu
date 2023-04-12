package bili

const (
	OP_HEARTBEAT        = int32(2)
	OP_HEARTBEAT_REPLY  = int32(3)
	OP_SEND_SMS_REPLY   = int32(5)
	OP_AUTH             = int32(7)
	OP_AUTH_REPLY       = int32(8)
)

const (
	MaxBodySize = int32(1 << 11)
	// size
	CmdSize       = 4
	PackSize      = 4
	HeaderSize    = 2
	VerSize       = 2
	OperationSize = 4
	SeqIdSize     = 4
	HeartbeatSize = 4
	RawHeaderSize = PackSize + HeaderSize + VerSize + OperationSize + SeqIdSize
	MaxPackSize   = MaxBodySize + int32(RawHeaderSize)
	// offset
	PackOffset      = 0
	HeaderOffset    = PackOffset + PackSize
	VerOffset       = HeaderOffset + HeaderSize
	OperationOffset = VerOffset + VerSize
	SeqIdOffset     = OperationOffset + OperationSize
	HeartbeatOffset = SeqIdOffset + SeqIdSize
	OP_RAW          = int32(11)
)

const (
	ProtoVersion0 = iota
	ProtoVersion1
	ProtoVersion2
	ProtoVersion3
)

type Proto struct {
	PacketLength int32
	HeaderLength int16
	Version      int16
	Operation    int32
	SequenceId   int32
	Body         []byte
	BodyMuti     [][]byte
	// 解析时的错误消息,非消息内容
	ErrMsg 		 error
}


type AuthRespParam struct {
	Code int64 `json:"code,omitempty"`
}
