package interceptors

import "strings"

func isLocalService(fullMethodName string) bool {
	return strings.Contains(fullMethodName, "kenda.mesync.Cloud")
}

func isCloudService(fullMethodName string) bool {
	return strings.Contains(fullMethodName, "kenda.mesync.Mesync")
}
