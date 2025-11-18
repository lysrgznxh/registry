package util

import (
	"agent-network-protocol/registry/core"
	"agent-network-protocol/registry/util/log"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

// 记录api访问
func ApiAccess(r *http.Request) {
	log.GetModuleLogger(core.LmServApi).WithFields(logrus.Fields{
		"url":    fmt.Sprintf("%s%s", r.Host, r.RequestURI),
		"method": r.Method,
	}).Info("api")
}
