package LendmeClient

import (
	"afr_auth_center/AuthCenterClient"
	"net/url"
	"time"
)

type Loyalty_Client struct {
	Protocol        string //http or https
	Hostname        string //name or IP
	LoyaltyPort     string
	LoyaltyModule   string
	LoyaltyVersion  string
	S2S_AccessToken string        //system to system access token
	Timeout         time.Duration //timeout if no reply after X seconds
	AUC_client      *AuthCenterClient.AUC_Client
}

type LendMe struct {
	LoyaltyClient *Loyalty_Client
}

func InitHostConfig(Protocol,
	Hostname,
	LoyaltyPort,
	LoyaltyModule,
	LoyaltyVersion,
	S2S_AccessToken string,
	Timeout time.Duration) (HostConfig Loyalty_Client) {
	HostConfig.Protocol = Protocol
	HostConfig.Hostname = Hostname
	HostConfig.LoyaltyPort = LoyaltyPort
	HostConfig.LoyaltyModule = LoyaltyModule
	HostConfig.LoyaltyVersion = LoyaltyVersion
	HostConfig.S2S_AccessToken = S2S_AccessToken
	HostConfig.Timeout = Timeout * time.Second
	return
}

func NewLoyaltyClient(Config Loyalty_Client) (conn *LendMe) {
	if Config.Hostname != "" && Config.LoyaltyPort != "" {
		conn = new(LendMe)
		conn.createLoyaltyClient(Config)
		auc := AuthCenterClient.NewAUCClient(*Config.AUC_client)
		conn.LoyaltyClient.AUC_client = auc.AUCClient
		conn.LoyaltyClient.S2S_AccessToken = auc.AUCClient.S2S_AccessToken
	}
	return
}

func (conn *LendMe) createLoyaltyClient(Config Loyalty_Client) (err error) {
	if conn.LoyaltyClient == nil {
		conn.LoyaltyClient = new(Loyalty_Client)
		conn.LoyaltyClient.Protocol = Config.Protocol
		conn.LoyaltyClient.Hostname = Config.Hostname
		conn.LoyaltyClient.LoyaltyPort = Config.LoyaltyPort
		conn.LoyaltyClient.LoyaltyModule = Config.LoyaltyModule
		conn.LoyaltyClient.LoyaltyVersion = Config.LoyaltyVersion
		conn.LoyaltyClient.S2S_AccessToken = Config.S2S_AccessToken
		conn.LoyaltyClient.Timeout = Config.Timeout
	}
	return nil
}

func UrlQueryEscape(text string) string {
	return url.QueryEscape(text)
}
