package xmpp // import "gosrc.io/xmpp"

import (
	"encoding/xml"
	"reflect"
)

// ============================================================================
// Message Packet

type Message struct {
	XMLName xml.Name `xml:"message"`
	PacketAttrs
	Subject    string         `xml:"subject,omitempty"`
	Body       string         `xml:"body,omitempty"`
	Thread     string         `xml:"thread,omitempty"`
	Error      Err            `xml:"error,omitempty"`
	Extensions []MsgExtension `xml:",omitempty"`
	X          *MsgXOOB       `xml:",omitempty"`
}

func (Message) Name() string {
	return "message"
}

func NewMessage(msgtype, from, to, id, lang string) Message {
	return Message{
		XMLName: xml.Name{Local: "message"},
		PacketAttrs: PacketAttrs{
			Id:   id,
			From: from,
			To:   to,
			Type: msgtype,
			Lang: lang,
		},
	}
}

type messageDecoder struct{}

var message messageDecoder

func (messageDecoder) decode(p *xml.Decoder, se xml.StartElement) (Message, error) {
	var packet Message
	err := p.DecodeElement(&packet, &se)
	return packet, err
}

// TODO: Support missing element (thread, extensions) by using proper marshaller
func (msg *Message) XMPPFormat() string {
	out, _ := xml.MarshalIndent(msg, "", "")
	return string(out)
}

// UnmarshalXML implements custom parsing for IQs
func (msg *Message) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	msg.XMLName = start.Name

	// Extract packet attributes
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			msg.Id = attr.Value
		}
		if attr.Name.Local == "type" {
			msg.Type = attr.Value
		}
		if attr.Name.Local == "to" {
			msg.To = attr.Value
		}
		if attr.Name.Local == "from" {
			msg.From = attr.Value
		}
		if attr.Name.Local == "lang" {
			msg.Lang = attr.Value
		}
	}

	// decode inner elements
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}

		switch tt := t.(type) {

		case xml.StartElement:
			var elt interface{}
			elementType := tt.Name.Space

			if extensionType := msgTypeRegistry[elementType][tt.Name.Local]; extensionType != nil {
				val := reflect.New(extensionType)
				elt = val.Interface()
				if msgExt, ok := elt.(MsgExtension); ok {
					err = d.DecodeElement(elt, &tt)
					if err != nil {
						return err
					}
					msg.Extensions = append(msg.Extensions, msgExt)
				}
			} else {
				// Decode default message elements
				var err error
				switch tt.Name.Local {
				case "body":
					err = d.DecodeElement(&msg.Body, &tt)
				case "thread":
					err = d.DecodeElement(&msg.Thread, &tt)
				case "subject":
					err = d.DecodeElement(&msg.Subject, &tt)
				case "error":
					err = d.DecodeElement(&msg.Error, &tt)
				}
				if err != nil {
					return err
				}
			}

		case xml.EndElement:
			if tt == start.End() {
				return nil
			}
		}
	}
}

// ============================================================================
// Message extensions
// Provide ability to add support to XMPP extension tags on messages

type MsgExtension interface {
}

// XEP-0184
const NSSpaceXEP0184Receipt = "urn:xmpp:receipts"

type ReceiptRequest struct {
	MsgExtension
	XMLName xml.Name `xml:"urn:xmpp:receipts request"`
}

type ReceiptReceived struct {
	MsgExtension
	XMLName xml.Name `xml:"urn:xmpp:receipts received"`
	Id      string   `xml:"id,attr"`
}

// XEP-0333
const NSSpaceXEP0333ChatMarkers = "urn:xmpp:chat-markers:0"

type ChatMarkerMarkable struct {
	MsgExtension
	XMLName xml.Name `xml:"urn:xmpp:chat-markers:0 markable"`
}

type ChatMarkerReceived struct {
	MsgExtension
	XMLName xml.Name `xml:"urn:xmpp:chat-markers:0 received"`
	Id      string   `xml:"id,attr"`
}

type ChatMarkerDisplayed struct {
	MsgExtension
	XMLName xml.Name `xml:"urn:xmpp:chat-markers:0 displayed"`
	Id      string   `xml:"id,attr"`
}

type ChatMarkerAcknowledged struct {
	MsgExtension
	XMLName xml.Name `xml:"urn:xmpp:chat-markers:0 acknowledged"`
	Id      string   `xml:"id,attr"`
}

// XEP-0066
type MsgXOOB struct {
	XMLName xml.Name `xml:"jabber:x:oob x"`
	URL     string   `xml:"url"`
	Desc    string   `xml:"desc,omitempty"`
}

// ============================================================================
// TODO: Make it configurable at to be able to easily add new XMPP extensions
//    in separate modules

var msgTypeRegistry = make(map[string]map[string]reflect.Type)

func init() {
	msgTypeRegistry[NSSpaceXEP0184Receipt] = make(map[string]reflect.Type)
	msgTypeRegistry[NSSpaceXEP0184Receipt]["request"] = reflect.TypeOf(ReceiptRequest{})
	msgTypeRegistry[NSSpaceXEP0184Receipt]["received"] = reflect.TypeOf(ReceiptReceived{})

	msgTypeRegistry[NSSpaceXEP0333ChatMarkers] = make(map[string]reflect.Type)
	msgTypeRegistry[NSSpaceXEP0333ChatMarkers]["markable"] = reflect.TypeOf(ChatMarkerMarkable{})
	msgTypeRegistry[NSSpaceXEP0333ChatMarkers]["received"] = reflect.TypeOf(ChatMarkerReceived{})
	msgTypeRegistry[NSSpaceXEP0333ChatMarkers]["displayed"] = reflect.TypeOf(ChatMarkerDisplayed{})
	msgTypeRegistry[NSSpaceXEP0333ChatMarkers]["acknowledged"] = reflect.TypeOf(ChatMarkerAcknowledged{})
}
