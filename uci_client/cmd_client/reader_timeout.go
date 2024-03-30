package cmd_client

type ReaderTimeout struct {
	Msg string
}

func NewReaderTimeout(msg string) *ReaderTimeout {
	return &ReaderTimeout{msg}
}

func (rt *ReaderTimeout) Error() string {
	return rt.Msg
}
