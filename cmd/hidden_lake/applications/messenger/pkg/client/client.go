package client

var (
	_ IClient = &sClient{}
)

type sClient struct {
	fBuilder   IBuilder
	fRequester IRequester
}

func NewClient(pBuilder IBuilder, pRequester IRequester) IClient {
	return &sClient{
		fBuilder:   pBuilder,
		fRequester: pRequester,
	}
}

func (p *sClient) PushMessage(pAliasName string, pPseudonym, pRequestID string, pBody []byte) error {
	return p.fRequester.PushMessage(pAliasName, p.fBuilder.PushMessage(pPseudonym, pRequestID, pBody))
}