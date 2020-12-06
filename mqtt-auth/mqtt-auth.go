package mqttauth

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func JWTHandler(tlsConfig *tls.Config, projectId string, privateKey string) (string, error) {
	log.Println("creating jwt handler")
	claims := jwt.StandardClaims{
		Audience:  projectId,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	log.Println("Load Private Key")
	keyBytes, err := ioutil.ReadFile(privateKey)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	log.Println("Parse Private Key")
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	log.Println("Sign String")
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return tokenString, err
}

func CreateTLSConfig(caCertPath string) (*tls.Config, error) {
	certpool := x509.NewCertPool()
	log.Println(caCertPath)
	pemCerts, err := ioutil.ReadFile(caCertPath)
	if err == nil {
		//log.Println(pemCerts)
		//log.Println(string(pemCerts))
		certpool.AppendCertsFromPEM(pemCerts)
		//log.Println(certpool)
		//time.Sleep(10 * time.Second)
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
	//log.Println(config)
	//log.Println(config.RootCAs)
	return config, nil
}

func Dummy() {

}
