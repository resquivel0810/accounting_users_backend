package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

func main() {
	host := "server29.hostfactory.ch"
	port := "465"
	address := net.JoinHostPort(host, port)

	fmt.Printf("=== Probando conectividad SMTP ===\n")
	fmt.Printf("Servidor: %s\n\n", address)

	// Método 1: Conexión TCP simple
	fmt.Println("1. Probando conexión TCP...")
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		fmt.Printf("✗ Error de conexión TCP: %v\n", err)
		fmt.Println("  Posibles causas:")
		fmt.Println("  - Firewall bloqueando el puerto")
		fmt.Println("  - Red del contenedor sin acceso externo")
		fmt.Println("  - Servidor SMTP no disponible")
		return
	}
	defer conn.Close()
	fmt.Println("✓ Conexión TCP exitosa")

	// Método 2: Conexión TLS
	fmt.Println("\n2. Probando conexión TLS...")
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: true,
	})
	
	err = tlsConn.Handshake()
	if err != nil {
		fmt.Printf("✗ Error en handshake TLS: %v\n", err)
		return
	}
	defer tlsConn.Close()
	fmt.Println("✓ Conexión TLS exitosa")
	fmt.Println("✓ Servidor SMTP es accesible desde este contenedor")
}
