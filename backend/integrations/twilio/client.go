package twilio

import (
	"strings"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type TwilioClient struct {
	client         *twilio.RestClient
	whatsappNumber string
}

func NewTwilioClient(twiloAccountSid, twiloAuthToken, whatsappNumber string) TwilioClient {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: twiloAccountSid,
		Password: twiloAuthToken,
	})

	return TwilioClient{
		client:         client,
		whatsappNumber: whatsappNumber,
	}
}

func (t *TwilioClient) SendWhatsAppMessage(to string, body string) (*twilioApi.ApiV2010Message, error) {
	to_normalized := strings.TrimPrefix(to, "whatsapp:")
	params := &twilioApi.CreateMessageParams{}
	params.SetTo("whatsapp:" + to_normalized)
	params.SetFrom("whatsapp:" + t.whatsappNumber)
	params.SetBody(body)

	resp, err := t.client.Api.CreateMessage(params)
	if err != nil {
		return &twilioApi.ApiV2010Message{}, err
	}

	return resp, nil
}
