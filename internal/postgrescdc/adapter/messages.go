package postgresadapter

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func Startup(database, user string) []byte {
	payload := bytes.NewBuffer(nil)
	binary.Write(payload, binary.BigEndian, uint32(196608))
	writeCString(payload, "user")
	writeCString(payload, user)
	writeCString(payload, "database")
	writeCString(payload, database)
	payload.WriteByte(0)
	return frame(payload.Bytes())
}
func Query(sql string) []byte {
	b := bytes.NewBuffer(nil)
	writeCString(b, sql)
	return tagged('Q', b.Bytes())
}
func Parse(name, sql string) []byte {
	b := bytes.NewBuffer(nil)
	writeCString(b, name)
	writeCString(b, sql)
	binary.Write(b, binary.BigEndian, uint16(0))
	return tagged('P', b.Bytes())
}
func Bind(portal, statement string) []byte {
	b := bytes.NewBuffer(nil)
	writeCString(b, portal)
	writeCString(b, statement)
	binary.Write(b, binary.BigEndian, uint16(0))
	binary.Write(b, binary.BigEndian, uint16(0))
	binary.Write(b, binary.BigEndian, uint16(0))
	return tagged('B', b.Bytes())
}
func writeCString(b *bytes.Buffer, s string) { b.WriteString(s); b.WriteByte(0) }
func frame(p []byte) []byte {
	b := bytes.NewBuffer(nil)
	binary.Write(b, binary.BigEndian, uint32(len(p)+4))
	b.Write(p)
	return b.Bytes()
}
func tagged(tag byte, p []byte) []byte { return append([]byte{tag}, frame(p)...) }
func ValidateFrame(b []byte, max int) error {
	if len(b) < 4 {
		return fmt.Errorf("postgres frame too short")
	}
	n := int(binary.BigEndian.Uint32(b[1:5]))
	if n < 4 || n-4 > max || n-4 != len(b)-5 {
		return fmt.Errorf("postgres frame length mismatch")
	}
	return nil
}
