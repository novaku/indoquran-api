package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/indoquran-api/src/config"
	"github.com/jbrodriguez/mlog"
)

func encode(txt, enType string) string {
	var (
		secret  []byte = []byte(config.Config.Secret.HMAC)
		message        = []byte(txt)
		result  string
	)

	hash := hmac.New(sha256.New, secret)
	hash.Write(message)

	if enType == "hexis" {
		// to lowercase hexits
		result = hex.EncodeToString(hash.Sum(nil))
		mlog.Info("hexis result for %s = %s", txt, result)
		return result
	} else if enType == "base64" {
		// to base64
		result = base64.StdEncoding.EncodeToString(hash.Sum(nil))
		mlog.Info("base64 result for %s = %s", txt, result)
		return result
	}

	return ""
}

// HMACValidation : to validate HMAC
func HMACValidation(secret, timestamp, encType string) (bool, error) {
	var (
		now        time.Time
		stamp      int64
		err        error
		stampTime  time.Time
		expireTime time.Time
		isValid    bool
	)

	now = time.Now()
	stamp, err = strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		mlog.Error(err)
		return false, err
	}
	stampTime = time.Unix(stamp, 0)

	expireTime = now.Add(12 * time.Hour) // HMAC expire in 1 hours
	if now.Before(stampTime) {
		err = fmt.Errorf("HMAC expire at %s", expireTime.Format("2006-01-02 15:04:05"))
		mlog.Error(err)
		return false, err
	}

	enc := encode(timestamp, encType)
	isValid = enc == secret
	mlog.Info("FE = %s, BE = %s, valid = %+v", secret, enc, isValid)
	return isValid, nil
}