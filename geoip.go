/*
 Copyright (c) 2012, Shane Hansen
All rights reserved.

 Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package gogeo

/*
#cgo LDFLAGS: -lGeoIP
#include <stdio.h>
#include <errno.h>
#include <GeoIPCity.h>
*/
import "C"
import "unsafe"
import "errors"
import "net"

type GeoIP struct {
    GeoIP *C.GeoIP
}

type GeoIPRecord struct {
    CountryCode       string
    CountryCode3      string
    CountryName       string
    Region            string
    City              string
    PostalCode        string
    Latitude          float64
    Longitude         float64
    AreaCode          int
    CharSet           int
    ContinentCode     string
    CountryConfidence byte
    RegionConfidence  byte
    CityConfidence    byte
    PostalConfidence  byte
    AccuracyRadius    int
}

func parseGeoIPRecord(c_record *C.GeoIPRecord) *GeoIPRecord {
    if c_record == nil {
        return nil
    }
    record := new(GeoIPRecord)
    record.CountryCode = C.GoString(c_record.country_code)
    record.CountryCode3 = C.GoString(c_record.country_code3)
    record.CountryName = C.GoString(c_record.country_name)
    record.Region = C.GoString(c_record.region)
    record.PostalCode = C.GoString(c_record.postal_code)
    record.Latitude = float64(c_record.latitude)
    record.Longitude = float64(c_record.longitude)
    record.AreaCode = int(c_record.area_code)
    record.CharSet = int(c_record.charset)
    record.ContinentCode = C.GoString(c_record.continent_code)
    record.CountryConfidence = byte(c_record.country_conf)
    record.RegionConfidence = byte(c_record.region_conf)
    record.CityConfidence = byte(c_record.city_conf)
    record.PostalConfidence = byte(c_record.postal_conf)
    record.AccuracyRadius = int(c_record.accuracy_radius)
    return record
}

//Open a database located a filename using flags
func Open(filename string, flags GeoIPOptions) (*GeoIP, error) {
    base := C.CString(filename)
    defer C.free(unsafe.Pointer(base))
    db := C.GeoIP_open(base, C.int(flags))
    if db == nil {
        return nil, errors.New("Cannot create GeoIP object")
    }
    geoIP := &GeoIP{db}
    return geoIP, nil
}

//return a string description of the database
func (self *GeoIP) Info() string {
    return C.GoString(C.GeoIP_database_info(self.GeoIP))
}

//Gracefully close the database
func (self *GeoIP) Close() {
    if self.GeoIP != nil {
        C.GeoIP_delete(self.GeoIP)
        self.GeoIP = nil
    }
}

func (self *GeoIP) DatabaseEdition() GeoIPDBType {
    return GeoIPDBType(C.GeoIP_database_edition(self.GeoIP))
}

func IPv4ToInt(ip []byte) uint32 {
    var ipaddr uint32
    ipaddr |= uint32(ip[0]) << 24
    ipaddr |= uint32(ip[1]) << 16
    ipaddr |= uint32(ip[2]) << 8
    ipaddr |= uint32(ip[3])
    return ipaddr
}

func (self *GeoIP) RecordByIPAddr(addr *net.IPAddr) *GeoIPRecord {
    return self.RecordByIP(addr.IP)
}

func (self *GeoIP) RecordByIP(ip net.IP) *GeoIPRecord {
    if len(ip) == 4 {
        ipaddr := IPv4ToInt(ip)
        record := C.GeoIP_record_by_ipnum(self.GeoIP, C.ulong(ipaddr))
        defer C.GeoIPRecord_delete(record)
        return parseGeoIPRecord(record)
    } else if len(ip) == 16 {
        //having weird problems with cgo and in6_addr, maybe the
        //typedef is messing it up?
        addr_buf := C.CString(ip.String())
        defer C.free(unsafe.Pointer(addr_buf))
        record := C.GeoIP_record_by_name_v6(self.GeoIP, addr_buf)
        if record != nil {
            defer C.GeoIPRecord_delete(record)
            return parseGeoIPRecord(record)
        }
        return nil

    }
    return nil
}
