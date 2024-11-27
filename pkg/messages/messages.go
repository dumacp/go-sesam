package messages

type MsgApdu struct {
	Data []byte
}
type MsgApduResponse struct {
	Data []byte
}
type MsgClose struct{}
type MsgOpen struct{}
type MsgAck struct {
	Error string
}
type MsgEncryptRequest struct {
	MsgID    string
	KeySlot  int
	DevInput []byte
	Data     []byte
	IV       []byte
}
type MsgEncryptResponse struct {
	MsgID  string
	Cipher []byte
	SamUID string
}
type MsgDecryptRequest struct {
	MsgID    string
	KeySlot  int
	DevInput []byte
	Data     []byte
	IV       []byte
}
type MsgDecryptResponse struct {
	MsgID  string
	Plain  []byte
	SamUID string
}
type MsgDumpSecretKeyRequest struct {
	KeySlot int
}
type MsgDumpSecretKeyResponse struct {
	Data []byte
}
type MsgCreateEntryKeyRequest struct {
	ContainKeysID []string
	EntryKeyID    string
	ChangeKeySlot int
	Keys          []byte
	AuthHost      bool
	OfflineKey    bool
	PICCKey       bool
	Alg           string
}
type MsgEntryKeyRequest struct {
	EntryKeyID string
	PersoKeyID string
	TargetSlot int
}
type MsgEntryKeyResponse struct {
	Error []byte
	Data  []byte
}
type MsgCreateKeyRequest struct {
	KeySlot int
	Alg     string
}
type MsgImportKeyRequest struct {
	KeySlot int
	Data    []byte
	Alg     string
}
type MsgEnableKeysRequest struct{}
type MsgEnableKeysResponse struct {
	Data []int
}
type MsgAuth struct {
	Key     string
	Version int
	Slot    int
}

type MsgGetUid struct {
}

type MsgUid struct {
	UID string
}
