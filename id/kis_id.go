package id

import (
	"kis-flow/common"
	"strings"

	"github.com/google/uuid"
)

func KisId(prefix ...string) (kisId string) {
	idStr := strings.Replace(uuid.New().String(), "-", "", -1)
	kisId = formatKisId(idStr, prefix...)
	return
}

func formatKisId(idStr string, prefix ...string) string {
	var kisId string
	for _, fix := range prefix {
		kisId += fix
		kisId += common.KisIdJoinChar
	}
	kisId += idStr

	return kisId
}
