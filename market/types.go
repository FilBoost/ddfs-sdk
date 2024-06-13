package market

const (
	TTDDfsSync = "DDFSSync"
	TTDDfsUrl  = "DDFSUrl"
	TTDDfsHTTP = "DDFShttp"
)

func IsDDTransferType(transferType string) bool {
	switch transferType {
	case TTDDfsSync, TTDDfsUrl, TTDDfsHTTP:
		return true
	}
	return false
}
