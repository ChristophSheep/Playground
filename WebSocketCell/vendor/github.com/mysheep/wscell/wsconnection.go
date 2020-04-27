package wscell

type connection struct {
	srcAttrName  string
	destAddress  string
	destAttrName string
}

// Connection TODO
type Connection interface {
	SrcAttrName() string
	DestAddress() string
	DestAttrName() string
}

// CreateConnection TODO
func CreateConnection(srcAttrName string, destAddress string, destAttrName string) Connection {
	return connection{
		srcAttrName:  srcAttrName,
		destAddress:  destAddress,
		destAttrName: destAttrName,
	}
}

func (c connection) SrcAttrName() string {
	return c.srcAttrName
}

func (c connection) DestAddress() string {
	return c.destAddress
}

func (c connection) DestAttrName() string {
	return c.destAttrName
}
