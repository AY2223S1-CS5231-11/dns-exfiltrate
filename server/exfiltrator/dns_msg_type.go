package exfiltrator

type dnsMsgType byte

func (t dnsMsgType) String() string {
	return string([]byte{byte(t)})
}

const (
	DNS_FILE_START dnsMsgType = '0'
	DNS_FILE_END   dnsMsgType = '1'
	DNS_FILE_DATA  dnsMsgType = '2'
)
