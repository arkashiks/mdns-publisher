package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grandcat/zeroconf"
)

func mustGetInterface(name string) net.Interface {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		log.Fatal(err)
	}
	return *iface
}

func main() {
    iface := mustGetInterface("veth3")
    
    services := []struct {
        name     string
        service  string
        port     int
        ip       string
        hostname string
    }{
        {"SOME_NAME", "_smb._tcp", 445, "SOME_IP", "SOME_HOSTNAME"},
        {"DO CHR",  "_smb._tcp", 445, "10.132.0.254", "SOME_HOSTNAME"},
    }

    var servers []*zeroconf.Server
    for _, svc := range services {
        server, err := zeroconf.RegisterProxy(
            svc.name, svc.service, "local.",
            svc.port, svc.hostname,  // ← unique per service
            []string{svc.ip},
            []string{}, []net.Interface{iface},
        )
        if err != nil {
            log.Fatal(err)
        }
        servers = append(servers, server)
    }
    defer func() {
        for _, s := range servers {
            s.Shutdown()
        }
    }()

    log.Println("mDNS publisher running")
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    <-sig
}
