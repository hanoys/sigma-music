package ports

import "errors"

var (
	ErrDecodePaymentKeyFailed = errors.New("failed to decode payment payload")
)
