package connector

type commandHandler func(conn *Conn, packet *Packet)

var commandHandlerMap = map[CmdType]commandHandler{
	CmdTypeUnknown:   UnknownHandler,
	CmdTypeHeartbeat: HeartbeatHandler,
}

func UnknownHandler(_ *Conn, _ *Packet) {
	panic("unknown command type in packet")
}

func HeartbeatHandler(conn *Conn, _ *Packet) {
	// TODO: 处理心跳
}
