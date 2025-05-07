package LendmeClient

import (
	"afr_auth_center/AuthCenterClient"
	"net/url"
	"time"
)

type Lendme_Client struct {
	Protocol        string //http or https
	Hostname        string //name or IP
	LendmePort      string
	LoyaltyPort     string
	LendMeModule    string
	LendMeVersion   string
	LoyaltyModule   string
	LoyaltyVersion  string
	S2S_AccessToken string        //system to system access token
	Timeout         time.Duration //timeout if no reply after X seconds
	AUC_client      *AuthCenterClient.AUC_Client
}

type LendMe struct {
	LendmeClient *Lendme_Client
}

func InitHostConfig(Protocol,
	Hostname,
	LendmePort,
	LoyaltyPort,
	LendMeModule,
	LendMeVersion,
	LoyaltyModule,
	LoyaltyVersion,
	S2S_AccessToken string,
	Timeout time.Duration) (HostConfig Lendme_Client) {
	HostConfig.Protocol = Protocol
	HostConfig.Hostname = Hostname
	HostConfig.LendmePort = LendmePort
	HostConfig.LoyaltyPort = LoyaltyPort
	HostConfig.LendMeModule = LendMeModule
	HostConfig.LendMeVersion = LendMeVersion
	HostConfig.LoyaltyModule = LoyaltyModule
	HostConfig.LoyaltyVersion = LoyaltyVersion
	HostConfig.S2S_AccessToken = S2S_AccessToken
	HostConfig.Timeout = Timeout * time.Second
	return
}

func NewLendmeClient(Config Lendme_Client) (conn *LendMe) {
	if Config.Hostname != "" && Config.LendmePort != "" {
		conn = new(LendMe)
		conn.createLendmeClient(Config)
		auc := AuthCenterClient.NewAUCClient(*Config.AUC_client)
		conn.LendmeClient.AUC_client = auc.AUCClient
		conn.LendmeClient.S2S_AccessToken = auc.AUCClient.S2S_AccessToken
	}
	return
}

func (conn *LendMe) createLendmeClient(Config Lendme_Client) (err error) {
	if conn.LendmeClient == nil {
		conn.LendmeClient = new(Lendme_Client)
		conn.LendmeClient.Protocol = Config.Protocol
		conn.LendmeClient.Hostname = Config.Hostname
		conn.LendmeClient.LendmePort = Config.LendmePort
		conn.LendmeClient.LoyaltyPort = Config.LoyaltyPort
		conn.LendmeClient.LendMeModule = Config.LendMeModule
		conn.LendmeClient.LendMeVersion = Config.LendMeVersion
		conn.LendmeClient.LoyaltyModule = Config.LoyaltyModule
		conn.LendmeClient.LoyaltyVersion = Config.LoyaltyVersion
		conn.LendmeClient.S2S_AccessToken = Config.S2S_AccessToken
		conn.LendmeClient.Timeout = Config.Timeout
	}
	return nil
}

func UrlQueryEscape(text string) string {
	return url.QueryEscape(text)
}
