package market

const (
	TTDDfsSync = "DDFSSync"
	TTDDfsUrl  = "DDFSUrl"
)

func IsDDTransferType(transferType string) bool {
	switch transferType {
	case TTDDfsSync, TTDDfsUrl:
		return true
	}
	return false
}
