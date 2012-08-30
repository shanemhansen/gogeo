package gogeo

import "net"
import "testing"

func TestSomething(t *testing.T) {
    db, err := Open("/usr/share/GeoIP/GeoLiteCity.dat", MemoryCache)
    if err != nil {
        t.Fatal("Coulnd't open database, do you have it installed?")
    }
    addr, err := net.ResolveIPAddr("ip4", "google.com")
    record := db.RecordByIPAddr(addr)
    if record.CountryCode != "US" {
        t.Fatal("it didn't work")
    }
}
