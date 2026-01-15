package walletservice

type HTTPWalletClient struct {
	Addr string
}

func NewHTTPWalletClient(addr string) *HTTPWalletClient {
	return &HTTPWalletClient{
		Addr: addr,
	}
}