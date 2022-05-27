module github.com/dumacp/go-sesam

go 1.16

replace github.com/dumacp/go-logs => ../go-logs

replace github.com/dumacp/smartcard => ../smartcard

require (
	github.com/AsynkronIT/protoactor-go v0.0.0-20210810091324-c3c6e02d5d46
	github.com/dumacp/go-logs v0.0.0-00010101000000-000000000000
	github.com/dumacp/smartcard v0.0.0-20210811220221-74a3161e587b
	github.com/looplab/fsm v0.2.0
)
