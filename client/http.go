package client

import (
	"encoding/base64"
	"fmt"
	"github.com/jcmturner/gokrb5/GSSAPI"
	"net/http"
)

// Get service ticket and set as the SPNEGO authorization header on HTTP request object
func (cl *Client) SetSPNEGOHeader(HTTPReq *http.Request) error {
	tkt, skey, err := cl.GetServiceTicket("HTTP/" + HTTPReq.Host)
	if err != nil {
		return fmt.Errorf("Could not get service ticket: %v", err)
	}
	negTokenInit, err := GSSAPI.NewNegTokenInitKrb5(*cl.Config, cl.Credentials.CName, tkt, skey)
	if err != nil {
		return fmt.Errorf("Could not create NegTokenInit: %v", err)
	}
	SPNEGOToken := GSSAPI.SPNEGO{
		Init:         true,
		NegTokenInit: negTokenInit,
	}
	nb, err := SPNEGOToken.Marshal()
	if err != nil {
		return fmt.Errorf("Could marshal SPNEGO: %v", err)
	}

	hs := "Negotiate " + base64.StdEncoding.EncodeToString(nb)
	HTTPReq.Header.Set("Authorization", hs)
	return nil
}