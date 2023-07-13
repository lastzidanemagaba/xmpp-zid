package main

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/gin-gonic/gin"
	"mellium.im/sasl"
	"mellium.im/xmpp"
	"mellium.im/xmpp/jid"
)

const (
	xmppServer = "xm.jtisrv.com"
	xmppPort   = 5223
	login      = "zidane@xm.jtisrv.com"
	jwtToken   = "zidane"

	xmppServer2 = "xm.jtisrv.com"
	xmppPort2   = 5223
	login2      = "zidane@xm.jtisrv.com"
	pass        = "zidane"
)

var xmppClient *xmpp.Session

func main() {
	r := gin.Default()

	r.GET("/connect", func(c *gin.Context) {
		j := jid.MustParse(login2)
		cfg := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         xmppServer2,
		}

		xmppClient, err := xmpp.DialClientSession(
			context.TODO(),
			j,
			xmpp.BindResource(),
			xmpp.StartTLS(cfg),
			xmpp.SASL("", pass, sasl.ScramSha1, sasl.Plain),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		defer xmppClient.Close()

		c.JSON(http.StatusOK, gin.H{
			"message": "connected",
		})
	})

	r.POST("/connect_jwt", func(c *gin.Context) {
		xmppservernya := c.PostForm("xmpp_server")
		user := c.PostForm("user")
		passnya := c.PostForm("pass")

		j := jid.MustParse(user)

		options := []xmpp.StreamFeature{
			xmpp.SASL("", passnya, sasl.ScramSha256Plus, sasl.ScramSha256, sasl.ScramSha1Plus, sasl.ScramSha1, sasl.Plain),
			xmpp.StartTLS(&tls.Config{
				InsecureSkipVerify: true,
				ServerName:         xmppservernya,
			}),
			xmpp.BindResource(),
		}
		var err error
		xmppClient, err = xmpp.DialClientSession(
			context.TODO(),
			j,
			options...,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "connected",
		})
	})

	r.Run(":6666")
}
