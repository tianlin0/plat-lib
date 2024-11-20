package udpdnslb

/*
 * tgLog库文件
 */

import (
	"fmt"
	"github.com/tianlin0/plat-lib/goroutines"
	"net"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// UDPDnsLb udp dns loadbalancer
type UDPDnsLb struct {
	nsLookupInterval time.Duration
	domain           string
	port             int
	prt              uint64
	addrs            []*net.UDPAddr
	// ipsComp use to check whether ips be changed after a dns request
	ipsComp string
	mu      sync.RWMutex
	ronce   sync.Once
	inited  bool
}

var dnslbmu = sync.RWMutex{}
var dnslbs = make(map[string]*UDPDnsLb)

// NewUDPDnsLb create udp dns loadbalancer
func NewUDPDnsLb(domain string, port int, lookupInterval time.Duration) *UDPDnsLb {
	dnslbKey := fmt.Sprintf("%s|%d|%d", domain, port, lookupInterval)
	dnslbmu.Lock()
	defer dnslbmu.Unlock()
	dnslb, ok := dnslbs[dnslbKey]
	if ok {
		return dnslb
	}
	dnslb = &UDPDnsLb{
		nsLookupInterval: lookupInterval,
		domain:           domain,
		port:             port,
	}
	dnslb.lookup()
	dnslb.registerLookup()
	dnslbs[dnslbKey] = dnslb
	return dnslb
}

// registerLookup register cyclically looking up
func (u *UDPDnsLb) registerLookup() {
	u.ronce.Do(func() {
		goroutines.GoAsyncHandler(func(params ...interface{}) {
			ticker := time.NewTicker(u.nsLookupInterval)
			for {
				select {
				case <-ticker.C:
					u.lookup()
				default:
					time.Sleep(500 * time.Millisecond)
				}
			}
		}, nil)
	})
}

// lookup do a dns lookup
func (u *UDPDnsLb) lookup() {
	ips := make(map[string]struct{}, 16)
	for i := 1; i <= 32; i++ {
		laddrs, err := net.LookupHost(u.domain)
		if err != nil {
			// leave the previous ip list remained and exit this lookup
			//log.Errorf("lookup host failed: %v", err)
			break
		}
		for _, addr := range laddrs {
			ips[addr] = struct{}{}
		}
		if len(laddrs) > 1 {
			// dns server return a list of ips
			break
		} else {
			// dns server return only one ip per request
			if i > len(ips)*2 && i > 12 {
				break
			}
		}
	}
	// will not update if ip list is empty
	if len(ips) == 0 {
		//log.Warnf("received empty ip list from dns of %s, will not update local cache", u.domain)
		return
	}

	naddrs := make([]*net.UDPAddr, 0, len(ips))
	sipArr := make([]string, 0, len(ips))
	for sip := range ips {
		ip := net.ParseIP(sip)
		udpAddr := net.UDPAddr{IP: ip, Port: u.port}
		naddrs = append(naddrs, &udpAddr)
		sipArr = append(sipArr, sip)
	}
	u.mu.Lock()
	u.addrs = naddrs
	u.inited = true
	u.mu.Unlock()
	// show this log only in case of ip list change
	sort.Strings(sipArr)
	sipsComp := strings.Join(sipArr, ";")
	if sipsComp != u.ipsComp {
		u.ipsComp = sipsComp
		//log.Infof("refresh udp addrs of %s:%d: %v", u.domain, u.port, naddrs)
	}
}

// GetUDPAddr get a udp addr from the pool. return nil if the poll has not been initialized
func (u *UDPDnsLb) GetUDPAddr() *net.UDPAddr {
	u.mu.RLock()
	if !u.inited || len(u.addrs) == 0 {
		return nil
	}
	addr := u.addrs[atomic.AddUint64(&u.prt, 1)%uint64(len(u.addrs))]
	u.mu.RUnlock()
	return addr
}
