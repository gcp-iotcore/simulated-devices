package mqttauth

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func jwtHandler(tlsConfig *tls.Config, projectId string, privateKey string) (string, error) {
	log.Println("creating jwt handler")
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = jwt.StandardClaims{
		Audience:  projectId,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}
	log.Println("Load Private Key")
	keyBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	log.Println("[main] Parse Private Key")
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	log.Println("[main] Sign String")
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return tokenString, err
}

func createTLSConfig(caCertPath string) (*tls.Config, error) {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(caCertPath)
	if err == nil {
		certpool.AppendCertsFromPEM(pemCerts)
	} else {
		log.Fatal(err)
		return nil, err
	}
	log.Println("Creating TLS Config")
	config := &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{},
		MinVersion:         tls.VersionTLS12,
	}
	return config, nil
}

func Dummy() {

}
