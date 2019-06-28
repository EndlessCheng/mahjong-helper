package tenhou

import (
	"encoding/xml"
)

// 需要注意的是，牌谱并未记录舍牌是手切还是摸切，
// 这里认为在摸牌后，只要切出的牌和摸的牌相同就认为是摸切，否则认为是手切
type RecordAction struct {
	XMLName xml.Name
	message
}

func (a *RecordAction) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	a.Tag = start.Name.Local
	type action RecordAction // 防止无限递归
	return d.DecodeElement((*action)(a), &start)
}

type Record struct {
	XMLName xml.Name        `xml:"mjloggm"`
	Actions []*RecordAction `xml:",any"`
}
