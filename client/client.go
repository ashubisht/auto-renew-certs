package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"
)

var (
	oldCer = "/Users/utkarshbisht/block8/projects/auto-renew-certs/localhost/cert.pem"
	oldKey = "/Users/utkarshbisht/block8/projects/auto-renew-certs/localhost/privkey.pem"
	cer = "/Users/utkarshbisht/block8/projects/auto-renew-certs/preethd.local/cert.pem"
	key = "/Users/utkarshbisht/block8/projects/auto-renew-certs/preethd.local/privkey.pem"
)

func main() {

	hold := make(chan uint8)

	conf := &tls.Config{
		MinVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
		InsecureSkipVerify: true,
	}
	
	// conn, err := net.Dial("tcp", "localhost:1888")
	// if err != nil {
	// 	fmt.Errorf("Error while initiating connection before tls", err)
	// 	os.Exit(1)
	// }
	// fmt.Println("Intiating tls connection")
	// tlsConn := tls.Client(conn, conf)
	// defer conn.Close()

	go func ()  {
		for true {

			// Dial start

			// conn, err := net.Dial("tcp", "localhost:1888")
			// if err != nil {
			// 	fmt.Errorf("Error while initiating connection before tls", err)
			// 	os.Exit(1)
			// }
			// fmt.Println("Intiating tls connection")
			// tlsConn := tls.Client(conn, conf)
			// defer conn.Close()

			tlsConn, err := tls.Dial("tcp", "localhost:1888", conf)
			if err != nil {
				fmt.Errorf("Error while initiating connection before tls", err)
				os.Exit(1)
			}

			defer tlsConn.Close()

			// Dial work done

			certSize := len(tlsConn.ConnectionState().PeerCertificates)
			fmt.Println("Cert length size = ", certSize)
			if certSize == 0 {
				fmt.Println("Size is 0. Sleeping")
				time.Sleep(2 * time.Second)
				continue
			}
			cert := tlsConn.ConnectionState().PeerCertificates[0]
			fmt.Println("=======Certificate information=======")
			fmt.Printf("Subject: %s\n", cert.Subject.CommonName)
			fmt.Printf("Alt Names: %s\n", strings.Join(cert.DNSNames, ", "))
			fmt.Printf("Issuer: %s\n", cert.Issuer.CommonName)
			fmt.Printf("Valid from: %s\n", cert.NotBefore.String())
			fmt.Printf("Valid until: %s\n", cert.NotAfter.String())
			fmt.Println("=======Certificate information end=======")
			time.Sleep(5 * time.Second)
		}
	}()

	<- hold

}
