package xmpp // import "gosrc.io/xmpp"

import "encoding/xml"

/*
Support for:
- XEP-0333 - Chat Markers: https://xmpp.org/extensions/xep-0333.html
*/

const NSMsgChatMarkers = "urn:xmpp:chat-markers:0"

type Markable struct {
	MsgExtension
	XMLName xml.Name `xml:"urn:xmpp:chat-markers:0 markable"`
}

type MarkReceived struct {
	MsgExtension
	XMLName xml.Name `xml:"urn:xmpp:chat-markers:0 received"`
	ID      string
}

type MarkDisplayed struct {
	MsgExtension
	XMLName xml.Name `xml:"urn:xmpp:chat-markers:0 displayed"`
	ID      string
}

type MarkAcknowledged struct {
	MsgExtension
	XMLName xml.Name `xml:"urn:xmpp:chat-markers:0 acknowledged"`
	ID      string
}

func init() {
	typeRegistry.MapExtension(PKTMessage, xml.Name{NSMsgChatMarkers, "markable"}, Markable{})
	typeRegistry.MapExtension(PKTMessage, xml.Name{NSMsgChatMarkers, "received"}, MarkReceived{})
	typeRegistry.MapExtension(PKTMessage, xml.Name{NSMsgChatMarkers, "displayed"}, MarkDisplayed{})
	typeRegistry.MapExtension(PKTMessage, xml.Name{NSMsgChatMarkers, "acknowledged"}, MarkAcknowledged{})
}
