package yoomoney

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"github.com/google/uuid"
	"github.com/hanoys/sigma-music/internal/domain"
	"github.com/hanoys/sigma-music/internal/ports"
	"net/url"
	"slices"
	"strconv"
)

type Config struct {
	Scheme string
	Host   string
	Path   string
	Wallet string
}

type PaymentYoomoneyGateway struct {
	config *Config
}

func NewPaymentGateway(config *Config) *PaymentYoomoneyGateway {
	return &PaymentYoomoneyGateway{
		config: config,
	}
}

func (g *PaymentYoomoneyGateway) GetPaymentUrl(ctx context.Context, payload domain.PaymentPayload) (url.URL, error) {
	userUUID, _ := uuid.Parse(payload.UserID.String())
	userUUIDBytes, _ := userUUID.MarshalBinary()

	paySumBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(paySumBytes, uint64(payload.PaymentSum))

	dataBytes := slices.Concat(userUUIDBytes, paySumBytes)
	encodedData := base64.StdEncoding.EncodeToString(dataBytes)
	formParams := url.Values{
		"sum":           {strconv.FormatInt(payload.PaymentSum, 10)},
		"receiver":      {g.config.Wallet},
		"quickpay-form": {"donate"},
		"label":         {encodedData},
	}

	return url.URL{
		Scheme:   g.config.Scheme,
		Host:     g.config.Host,
		Path:     g.config.Path,
		RawQuery: formParams.Encode(),
	}, nil
}

func (g *PaymentYoomoneyGateway) ProcessPayment(ctx context.Context, key string) (domain.PaymentPayload, error) {
	dataBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return domain.PaymentPayload{}, ports.ErrDecodePaymentKeyFailed
	}

	var userID, courseID uuid.UUID
	err = userID.UnmarshalBinary(dataBytes[:16])
	if err != nil {
		return domain.PaymentPayload{}, ports.ErrDecodePaymentKeyFailed
	}

	err = courseID.UnmarshalBinary(dataBytes[16:32])
	if err != nil {
		return domain.PaymentPayload{}, ports.ErrDecodePaymentKeyFailed
	}

	paySum := binary.LittleEndian.Uint64(dataBytes[32:])
	if err != nil {
		return domain.PaymentPayload{}, ports.ErrDecodePaymentKeyFailed
	}

	return domain.PaymentPayload{
		UserID:     userID,
		PaymentSum: int64(paySum),
	}, nil
}
