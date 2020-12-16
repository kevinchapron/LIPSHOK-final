package database

type ProtocolActivated struct {
	Name      string `json:"name"`
	Activated bool   `json:"activated"`
}
type Protocols struct {
	List []ProtocolActivated `json:"protocols"`
}
