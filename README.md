## Go bindings to the GeoIP city database

gogeo allows you you to access a GeoIP city database to determine
Geographical information from an IP address.

### Features

- IPv4 support
- GeoRecord support
- Unlike other libraries, doesn't read the whole database into memory
- Fine grained control over database open modes, caching, etc.
- Works with go net.Addr's, which simplifies integration with go network programs.

### Limitations

- No IPv6 support yet
- cgo based (not suitable for app engine deployments)


### Setup

- Install libgeoip1 libgeoip-dev (on ubuntu)
- go get -u github.com/shanemhansen/gogeo

### Example usage

- Look at the unit tests

    package main
    import "github.com/shanemhansen/gogeo"
    import "net"
    import "fmt"
    
    func main() {
        db, err := gogeo.Open("/usr/share/GeoIP/GeoLiteCity.dat", gogeo.MemoryCache)
        if err != nil {
            panic(err)
        }
        addr, err := net.ResolveIPAddr("ip4", "google.com")
        record := db.RecordByIPAddr(addr)
        fmt.Printf("hello, from %s", record.CountryCode)
   }
