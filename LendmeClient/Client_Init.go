package LendmeClient

import (
	"net/url"
	"time"
)

type Lendme_Client struct {
	Protocol        string //http or https
	Hostname        string //name or IP
	Port            string
	Module          string
	Version         string
	S2S_AccessToken string        //system to system access token
	Timeout         time.Duration //timeout if no reply after X seconds
}

type LendMe struct {
	LendmeClient *Lendme_Client
}

func InitHostConfig(Protocol,
	Hostname,
	Port,
	Module,
	Version,
	S2S_AccessToken string,
	Timeout time.Duration) (HostConfig Lendme_Client) {
	HostConfig.Protocol = Protocol
	HostConfig.Hostname = Hostname
	HostConfig.Port = Port
	HostConfig.Module = Module
	HostConfig.Version = Version
	HostConfig.S2S_AccessToken = S2S_AccessToken
	HostConfig.Timeout = Timeout * time.Second
	return
}

func NewLendmeClient(Config Lendme_Client) (conn *LendMe) {
	if Config.Hostname != "" && Config.Port != "" {
		conn = new(LendMe)
		conn.createLendmeClient(Config)
	}
	// if !conn.LendmeClient.IsAuthenticated {
	// 	log.Println("Failed to authenticate access: " + conn.LendmeClient.Auth_Err_Description)
	// } else {
	// 	log.Println("Authenticated Successfuly")
	// }
	return
}

func (conn *LendMe) createLendmeClient(Config Lendme_Client) (err error) {
	if conn.LendmeClient == nil {
		conn.LendmeClient = new(Lendme_Client)
		conn.LendmeClient.Protocol = Config.Protocol
		conn.LendmeClient.Hostname = Config.Hostname
		conn.LendmeClient.Port = Config.Port
		conn.LendmeClient.Module = Config.Module
		conn.LendmeClient.Version = Config.Version
		conn.LendmeClient.S2S_AccessToken = Config.S2S_AccessToken
		conn.LendmeClient.Timeout = Config.Timeout
	}
	// if conn.LendmeClient.S2S_AccessToken != "" {
	// 	_, err := conn.LendmeClient.ValidateToken(conn.AUCClient.S2S_AccessToken)
	// 	if err != nil {
	// 		conn.AUCClient.IsAuthenticated = false
	// 		conn.AUCClient.Auth_Err_Description = err.Error()
	// 		return err
	// 	}
	// 	conn.AUCClient.IsAuthenticated = true
	// 	conn.AUCClient.Auth_Err_Description = ""
	// 	return nil
	// } else {
	// 	if conn.AUCClient.S2S_Username != "" && conn.AUCClient.S2S_Password != "" {
	// 		sr, err := conn.AUCClient.Login(conn.AUCClient.S2S_Username, conn.AUCClient.S2S_Password)
	// 		if err != nil {
	// 			conn.AUCClient.IsAuthenticated = false
	// 			conn.AUCClient.Auth_Err_Description = err.Error()
	// 			log.Println("Failed to authenticate (Login): " + err.Error())
	// 			return err
	// 		} else {
	// 			rt, err := conn.AUCClient.RefreshToken(sr.Token)
	// 			if err != nil {
	// 				conn.AUCClient.IsAuthenticated = false
	// 				conn.AUCClient.Auth_Err_Description = err.Error()
	// 				log.Println("Failed to authenticate (RefreshToken): " + err.Error())
	// 				return err
	// 			} else {
	// 				conn.AUCClient.S2S_AccessToken = rt.Token // save token global variable to re-use
	// 				conn.AUCClient.IsAuthenticated = true
	// 				conn.AUCClient.Auth_Err_Description = ""
	// 				return nil
	// 			}
	// 		}
	// 	} else {
	// 		err = errors.New("login or password information not provided")
	// 		conn.AUCClient.IsAuthenticated = false
	// 		conn.AUCClient.Auth_Err_Description = err.Error()
	// 		return err
	// 	}
	// }
	return nil
}

func UrlQueryEscape(text string) string {
	return url.QueryEscape(text)
}
