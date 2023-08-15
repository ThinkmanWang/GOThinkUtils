package qqwry

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/yinheli/mahonia"
)

const (
	INDEX_LEN       = 7
	REDIRECT_MODE_1 = 0x01
	REDIRECT_MODE_2 = 0x02
)

type QQwry struct {
	buff  []byte
	start uint32
	end   uint32
}
type Rq struct {
	Ip       string `json:"ip"`
	Province string `json:"province"`
	City     string `json:"city"`
	Err      int    `json:"code"`
	Msg      string `json:"msg"`
	Operator string `json:"operator"`
}

func NewQQwry(file string) (qqwry *QQwry) {
	qqwry = &QQwry{}
	f, e := os.Open(file)
	if e != nil {
		log.Println(e)
		return nil
	}
	defer f.Close()
	qqwry.buff, e = ioutil.ReadAll(f)
	if e != nil {
		log.Println(e)
		return nil
	}
	qqwry.start = binary.LittleEndian.Uint32(qqwry.buff[:4])
	qqwry.end = binary.LittleEndian.Uint32(qqwry.buff[4:8])
	return qqwry
}
func (this *Rq) String() string {

	d, _ := json.Marshal(this)
	return string(d)
}
func (this *QQwry) Find(ip string) *Rq {
	rq := &Rq{Ip: ip}
	if this.buff == nil {
		rq.Err = 3
		rq.Msg = "QQwry没有初始化"
		return rq
	}

	var country []byte
	var area []byte
	ip_1 := net.ParseIP(ip)
	if ip_1 == nil {
		rq.Err = 1
		rq.Msg = "错误的IP格式"
		return rq
	}
	offset := this.searchRecord(binary.BigEndian.Uint32(ip_1.To4()))
	if offset <= 0 {
		rq.Err = 2
		rq.Msg = "IP地址没找到归属地"
		return rq
	}
	mode := this.readMode(offset + 4)
	if mode == REDIRECT_MODE_1 {
		countryOffset := this.readUint32FromByte3(offset + 5)

		mode = this.readMode(countryOffset)
		if mode == REDIRECT_MODE_2 {
			c := this.readUint32FromByte3(countryOffset + 1)
			country = this.readString(c)
			countryOffset += 4
			area = this.readArea(countryOffset)

		} else {
			country = this.readString(countryOffset)
			countryOffset += uint32(len(country) + 1)
			area = this.readArea(countryOffset)
		}

	} else if mode == REDIRECT_MODE_2 {
		countryOffset := this.readUint32FromByte3(offset + 5)
		country = this.readString(countryOffset)
		area = this.readArea(offset + 8)
	}
	enc := mahonia.NewDecoder("gbk")

	szProvinceCity := enc.ConvertString(string(country))

	szProvinceCity = strings.ReplaceAll(szProvinceCity, "省", "省 ")
	szProvinceCity = strings.ReplaceAll(szProvinceCity, "自治区", "自治区 ")
	szProvinceCity = strings.ReplaceAll(szProvinceCity, "市", "市 ")
	szProvinceCity = strings.ReplaceAll(szProvinceCity, "县", "县 ")

	lstArea := strings.Split(szProvinceCity, " ")

	if lstArea != nil {
		if 1 == len(lstArea) {
			rq.Province = strings.Trim(lstArea[0], " ")
		} else if len(lstArea) >= 2 {
			rq.Province = strings.Trim(lstArea[0], " ")
			rq.City = strings.Trim(lstArea[1], " ")
		} else {
			rq.Province = ""
			rq.City = ""
		}
	}
	rq.Operator = enc.ConvertString(string(area))

	//do for Shanghai, Beijing, Tianjing, Chongqing
	if "北京市" == rq.Province ||
		"天津市" == rq.Province ||
		"重庆市" == rq.Province ||
		"上海市" == rq.Province {
		rq.City = rq.Province
	}

	return rq
}

func (this *QQwry) readUint32FromByte3(offset uint32) uint32 {
	return byte3ToUInt32(this.buff[offset : offset+3])
}
func (this *QQwry) readMode(offset uint32) byte {
	return this.buff[offset : offset+1][0]
}

func (this *QQwry) readString(offset uint32) []byte {

	i := 0
	for {

		if this.buff[int(offset)+i] == 0 {
			break
		} else {
			i++
		}

	}
	return this.buff[offset : int(offset)+i]
}

func (this *QQwry) readArea(offset uint32) []byte {
	mode := this.readMode(offset)
	if mode == REDIRECT_MODE_1 || mode == REDIRECT_MODE_2 {
		areaOffset := this.readUint32FromByte3(offset + 1)
		if areaOffset == 0 {
			return []byte("")
		} else {
			return this.readString(areaOffset)
		}
	} else {
		return this.readString(offset)
	}
	return []byte("")
}

func (this *QQwry) getRecord(offset uint32) []byte {
	return this.buff[offset : offset+INDEX_LEN]
}

func (this *QQwry) getIPFromRecord(buf []byte) uint32 {
	return binary.LittleEndian.Uint32(buf[:4])
}

func (this *QQwry) getAddrFromRecord(buf []byte) uint32 {
	return byte3ToUInt32(buf[4:7])
}

func (this *QQwry) searchRecord(ip uint32) uint32 {

	start := this.start
	end := this.end

	// log.Printf("len info %v, %v ---- %v, %v", start, end, hex.EncodeToString(header[:4]), hex.EncodeToString(header[4:]))
	for {
		mid := this.getMiddleOffset(start, end)
		buf := this.getRecord(mid)
		_ip := this.getIPFromRecord(buf)

		// log.Printf(">> %v, %v, %v -- %v", start, mid, end, hex.EncodeToString(buf[:4]))

		if end-start == INDEX_LEN {
			//log.Printf(">> %v, %v, %v -- %v", start, mid, end, hex.EncodeToString(buf[:4]))
			offset := this.getAddrFromRecord(buf)
			buf = this.getRecord(mid + INDEX_LEN)
			if ip < this.getIPFromRecord(buf) {
				return offset
			} else {
				return 0
			}
		}

		// 找到的比较大，向前移
		if _ip > ip {
			end = mid
		} else if _ip < ip { // 找到的比较小，向后移
			start = mid
		} else if _ip == ip {
			return byte3ToUInt32(buf[4:7])
		}

	}
	return 0
}

func (this *QQwry) getMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / INDEX_LEN) >> 1
	return start + records*INDEX_LEN
}

func byte3ToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
