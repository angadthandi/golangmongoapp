package genuuid

import (
	"math/rand"
	"strings"
	"time"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

func GenUUID() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Debugf("Unable to create uuid: %s", err)

		// fallback
		uuid := genFromTimeStamp()
		return uuid
	}

	return uuid.String()
}

func genFromTimeStamp() string {
	rand.Seed(time.Now().UTC().UnixNano())
	return randomString(32)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	retStr := string(bytes)

	return strings.ToLower(retStr)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
