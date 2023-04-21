package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/rpc"
	"time"
)

var (
	oldCer = "/Users/utkarshbisht/block8/projects/auto-renew-certs/certs/localhost/cert.pem"
	oldKey = "/Users/utkarshbisht/block8/projects/auto-renew-certs/certs/localhost/privkey.pem"
	cer = "/Users/utkarshbisht/block8/projects/auto-renew-certs/certs/preethd.local/cert.pem"
	key = "/Users/utkarshbisht/block8/projects/auto-renew-certs/certs/preethd.local/privkey.pem"
)

// Handler
type Server struct {
	certType uint8 // 0 = local, 1 = preethd
	hold chan uint8
}

func (srvr *Server) DeliverValue(msg string, _ *interface{}) error {
	fmt.Println("Registered an incoming messsage")
	return nil
}
// Handler ends

func main() {

	srvr := &Server {
		certType: 0,
		hold: make(chan uint8),
	}
	handler := rpc.NewServer()
	handler.Register(srvr)


	config := &tls.Config{
		GetCertificate: srvr.LoadCertificate(),
		MinVersion:   tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	l, err := tls.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 1888), config)
	if err != nil {
		log.Fatalf("Error in starting tls tcp listener", err)
	}

	go func() {
		for {
			cxn, err := l.Accept()
			if err != nil {
				fmt.Errorf("Error Accept Request: %s\n", err)
				return
			}
			go handler.ServeConn(cxn)
		}
	}()
	fmt.Println("Server started at port", 1888)

	go func() {
    for true {
        time.Sleep(10 * time.Second)
				fmt.Println("Swapping certificate")
				srvr.SwapCerts()
				config.GetCertificate = srvr.LoadCertificate()
    }
}()
	<- srvr.hold
}

// Directly copied and then modified
func (srvr *Server) LoadCertificate() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	fmt.Println("Loading certificate")
	var certPath, keyPath string
		if(srvr.certType == 0){
			fmt.Println("load 0")
			certPath = oldCer
			keyPath = oldKey
		}else {
			fmt.Println("load 1")
			certPath = cer
			keyPath = key
		}
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		cer, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			fmt.Errorf("Error in loading oldCer and oldKey")
		}
		return &cer, err
	}
}

func (srvr *Server) SwapCerts()  {
	if(srvr.certType == 0){
		srvr.certType = 1
	}else {
		srvr.certType = 0
	}
}