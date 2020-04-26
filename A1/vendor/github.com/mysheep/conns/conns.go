package conns

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const json_url = "http://10.0.0.138/?connections/ajax_reset_conn"

// Example Json from Roouter

/*
{
	"response": [
	 {"id":"7","protocal":"udp(17)","lan":"0.0.0.0:68","wan_box":"0.0.0.0:68","wan":"127.0.0.1:67","wan_status":"0","timeout":"127","tx_bytes":"0/0","tx_packets":"0/0","bind":"DHCP","intf":"vIPTV","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18545","protocal":"udp(17)","lan":"*.*.*.*:68","wan_box":"*.*.*.*:68","wan":"*.*.*.*:*","wan_status":"0","timeout":"52","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"loop","route_mode":"","direction":"Eingehend","flags":".WT..F..V....","ctype":"N"}
		,{"id":"3733","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:49253","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:49253","wan":"[2a00:1450:400c:c00::bc]:5228","wan_status":"HERGESTELLT","timeout":"406","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L.........","ctype":"N"}
		,{"id":"15636","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:47311","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:47311","wan":"[2a00:1450:400d:809::200a]:443","wan_status":"HERGESTELLT","timeout":"0","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"15657","protocal":"tcp(6)","lan":"10.0.0.2:58074","wan_box":"100.72.228.161:60372","wan":"17.242.176.14:443","wan_status":"HERGESTELLT","timeout":"795","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"16145","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:49367","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:49367","wan":"[2a03:2880:f007:1:face:b00c:0:1]:443","wan_status":"HERGESTELLT","timeout":"2508","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"S"}
		,{"id":"16171","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:54090","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:54090","wan":"[2a03:2880:f007:1:face:b00c:0:1]:443","wan_status":"HERGESTELLT","timeout":"2517","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"S"}
		,{"id":"16180","protocal":"tcp(6)","lan":"10.0.0.6:37751","wan_box":"100.72.228.161:60409","wan":"172.217.19.110:443","wan_status":"HERGESTELLT","timeout":"2572","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"S"}
		,{"id":"16221","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:55819","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:55819","wan":"[2a00:1450:400d:805::2001]:443","wan_status":"HERGESTELLT","timeout":"2576","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"S"}
		,{"id":"16227","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:55820","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:55820","wan":"[2a00:1450:400d:805::2001]:443","wan_status":"HERGESTELLT","timeout":"2576","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"S"}
		,{"id":"16239","protocal":"tcp(6)","lan":"10.0.0.6:37660","wan_box":"100.72.228.161:60395","wan":"199.232.18.133:443","wan_status":"HERGESTELLT","timeout":"816","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"16517","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:40321","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:40321","wan":"[2a00:1450:400d:809::200a]:443","wan_status":"HERGESTELLT","timeout":"2858","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"S"}
		,{"id":"16735","protocal":"tcp(6)","lan":"10.0.0.9:57035","wan_box":"100.72.228.161:60521","wan":"17.57.146.52:5223","wan_status":"HERGESTELLT","timeout":"695","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L.........","ctype":"N"}
		,{"id":"16834","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57054","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57054","wan":"[2a00:1450:400d:805::200e]:443","wan_status":"ZEIT_WARTEN","timeout":"225","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"16857","protocal":"tcp(6)","lan":"10.0.0.6:58842","wan_box":"100.72.228.161:60638","wan":"128.116.112.44:443","wan_status":"HERGESTELLT","timeout":"25","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"16988","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:46797","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:46797","wan":"[2a00:1450:400d:808::2003]:443","wan_status":"HERGESTELLT","timeout":"3160","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"S"}
		,{"id":"17071","protocal":"tcp(6)","lan":"10.0.0.9:57036","wan_box":"100.72.228.161:60524","wan":"52.43.222.181:443","wan_status":"HERGESTELLT","timeout":"691","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17092","protocal":"tcp(6)","lan":"10.0.0.9:57047","wan_box":"100.72.228.161:60531","wan":"17.242.57.246:443","wan_status":"HERGESTELLT","timeout":"3397","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"S"}
		,{"id":"17097","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57063","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57063","wan":"[2a00:1450:400d:804::200a]:443","wan_status":"ZEIT_WARTEN","timeout":"225","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"17134","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57064","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57064","wan":"[2a00:1450:400d:803::200a]:443","wan_status":"HERGESTELLT","timeout":"882","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17136","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57062","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57062","wan":"[2a00:1450:400d:802::200e]:443","wan_status":"ZEIT_WARTEN","timeout":"225","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"17142","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57074","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57074","wan":"[2a00:1450:400d:809::200e]:443","wan_status":"ZEIT_WARTEN","timeout":"225","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"17150","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57071","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57071","wan":"[2a00:1450:400d:809::200e]:443","wan_status":"ZEIT_WARTEN","timeout":"225","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"17170","protocal":"tcp(6)","lan":"10.0.0.9:57066","wan_box":"100.72.228.161:60543","wan":"148.251.127.85:80","wan_status":"HERGESTELLT","timeout":"890","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17176","protocal":"tcp(6)","lan":"10.0.0.6:55740","wan_box":"100.72.228.161:60636","wan":"128.116.112.44:443","wan_status":"HERGESTELLT","timeout":"23","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"17195","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57061","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57061","wan":"[2a00:1450:400d:804::200a]:443","wan_status":"HERGESTELLT","timeout":"881","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17238","protocal":"tcp(6)","lan":"10.0.0.9:57068","wan_box":"100.72.228.161:60545","wan":"129.27.124.196:443","wan_status":"HERGESTELLT","timeout":"889","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17242","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57069","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57069","wan":"[2a00:1450:400d:808::2003]:443","wan_status":"HERGESTELLT","timeout":"862","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17261","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:53721","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:53721","wan":"[2a00:1450:400d:808::200a]:443","wan_status":"HERGESTELLT","timeout":"3389","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"S"}
		,{"id":"17278","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57091","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57091","wan":"[2a00:1450:400d:809::200e]:443","wan_status":"HERGESTELLT","timeout":"865","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17279","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57092","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57092","wan":"[2a00:1450:400d:803::200a]:443","wan_status":"HERGESTELLT","timeout":"860","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17316","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:58935","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:58935","wan":"[2a00:1450:400d:802::200a]:443","wan_status":"HERGESTELLT","timeout":"3156","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"S"}
		,{"id":"17335","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57078","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57078","wan":"[2a00:1450:400d:803::200a]:443","wan_status":"HERGESTELLT","timeout":"882","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17342","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57090","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57090","wan":"[2a00:1450:400d:805::200e]:443","wan_status":"HERGESTELLT","timeout":"865","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17427","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57057","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57057","wan":"[2a00:1450:400d:808::2003]:443","wan_status":"HERGESTELLT","timeout":"817","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17496","protocal":"tcp(6)","lan":"10.0.0.9:57038","wan_box":"100.72.228.161:60526","wan":"52.43.222.181:443","wan_status":"HERGESTELLT","timeout":"688","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17590","protocal":"tcp(6)","lan":"10.0.0.6:48114","wan_box":"100.72.228.161:60645","wan":"128.116.112.44:443","wan_status":"HERGESTELLT","timeout":"26","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"17727","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57162","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57162","wan":"[2a00:1450:400d:808::200e]:443","wan_status":"HERGESTELLT","timeout":"880","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17759","protocal":"tcp(6)","lan":"10.0.0.6:41556","wan_box":"100.72.228.161:60644","wan":"128.116.112.44:443","wan_status":"HERGESTELLT","timeout":"25","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"17776","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57161","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57161","wan":"[2a00:1450:400d:808::200e]:443","wan_status":"ZEIT_WARTEN","timeout":"28","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"17800","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57134","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57134","wan":"[2a00:1450:400d:803::200a]:443","wan_status":"HERGESTELLT","timeout":"882","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17801","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57132","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57132","wan":"[2a00:1450:400d:808::200e]:443","wan_status":"HERGESTELLT","timeout":"863","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"17830","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:56129","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:56129","wan":"[2a00:1450:400d:802::200e]:443","wan_status":"HERGESTELLT","timeout":"260","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"17841","protocal":"tcp(6)","lan":"10.0.0.6:46887","wan_box":"100.72.228.161:60635","wan":"172.217.16.110:443","wan_status":"HERGESTELLT","timeout":"252","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"17955","protocal":"tcp(6)","lan":"10.0.0.6:57003","wan_box":"100.72.228.161:60659","wan":"34.248.112.237:443","wan_status":"HERGESTELLT","timeout":"162","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"17998","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57186","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57186","wan":"[2a00:1450:400d:808::200e]:443","wan_status":"ZEIT_WARTEN","timeout":"225","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18003","protocal":"tcp(6)","lan":"10.0.0.6:56035","wan_box":"100.72.228.161:60660","wan":"52.31.72.97:443","wan_status":"HERGESTELLT","timeout":"98","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"18022","protocal":"tcp(6)","lan":"10.0.0.6:40069","wan_box":"100.72.228.161:60654","wan":"128.116.112.44:443","wan_status":"HERGESTELLT","timeout":"84","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"18028","protocal":"udp(17)","lan":"100.72.228.161:28395","wan_box":"100.72.228.161:28395","wan":"195.3.96.67:53","wan_status":"0","timeout":"114","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18051","protocal":"tcp(6)","lan":"10.0.0.6:42683","wan_box":"100.72.228.161:60657","wan":"128.116.112.44:443","wan_status":"HERGESTELLT","timeout":"85","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"18059","protocal":"tcp(6)","lan":"10.0.0.9:57307","wan_box":"100.72.228.161:60724","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"170","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18114","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:50509","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:50509","wan":"[2a00:1450:400d:804::200a]:443","wan_status":"HERGESTELLT","timeout":"648","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"18117","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:37286","wan_box":"[2001:871:10:ad4e:8929:82b3:5cec:be70]:37286","wan":"[2a00:1450:400d:804::200a]:443","wan_status":"HERGESTELLT","timeout":"648","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"18146","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:d8d0:d3bb:b71e:3599]:58246","wan_box":"[2001:871:10:ad4e:d8d0:d3bb:b71e:3599]:58246","wan":"[2a01:b740:a10:f100::4]:443","wan_status":"ZEIT_WARTEN","timeout":"146","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18147","protocal":"udp(17)","lan":"[2001:870:10:ad4e::a:1]:23413","wan_box":"[2001:870:10:ad4e::a:1]:23413","wan":"[2001:850:3040:1:0:a1:4:a11]:53","wan_status":"0","timeout":"129","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18159","protocal":"tcp(6)","lan":"10.0.0.6:55501","wan_box":"100.72.228.161:60689","wan":"39.103.7.176:80","wan_status":"HERGESTELLT","timeout":"419","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"...L......R..","ctype":"S"}
		,{"id":"18163","protocal":"tcp(6)","lan":"10.0.0.9:57253","wan_box":"100.72.228.161:60709","wan":"129.27.124.196:443","wan_status":"LETZTE_ACK","timeout":"4","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18214","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57205","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57205","wan":"[2a00:1450:400c:c02::bd]:443","wan_status":"ZEIT_WARTEN","timeout":"105","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18220","protocal":"tcp(6)","lan":"10.0.0.9:57308","wan_box":"100.72.228.161:60725","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"178","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18232","protocal":"tcp(6)","lan":"10.0.0.9:57311","wan_box":"100.72.228.161:60726","wan":"129.27.124.196:443","wan_status":"LETZTE_ACK","timeout":"199","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18244","protocal":"udp(17)","lan":"100.72.228.161:15714","wan_box":"100.72.228.161:15714","wan":"195.3.96.67:53","wan_status":"0","timeout":"29","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18287","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57281","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57281","wan":"[2a00:1450:400d:805::200a]:443","wan_status":"ZEIT_WARTEN","timeout":"240","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18312","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57251","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57251","wan":"[2a00:1450:400c:c06::bd]:443","wan_status":"ZEIT_WARTEN","timeout":"225","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18317","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57248","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57248","wan":"[2a00:1450:400d:809::200e]:443","wan_status":"ZEIT_WARTEN","timeout":"119","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18329","protocal":"udp(17)","lan":"[2001:870:10:ad4e::a:1]:31272","wan_box":"[2001:870:10:ad4e::a:1]:31272","wan":"[2001:850:3040:1:0:a1:4:a11]:53","wan_status":"0","timeout":"117","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18335","protocal":"udp(17)","lan":"100.72.228.161:32231","wan_box":"100.72.228.161:32231","wan":"195.3.96.67:53","wan_status":"0","timeout":"2","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18339","protocal":"udp(17)","lan":"100.72.228.161:12372","wan_box":"100.72.228.161:12372","wan":"195.3.96.67:53","wan_status":"0","timeout":"29","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18352","protocal":"udp(17)","lan":"10.0.0.2:123","wan_box":"100.72.228.161:60719","wan":"17.253.52.253:123","wan_status":"0","timeout":"29","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L.........","ctype":"N"}
		,{"id":"18356","protocal":"tcp(6)","lan":"10.0.0.9:57306","wan_box":"100.72.228.161:60723","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"163","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18363","protocal":"tcp(6)","lan":"10.0.0.9:57300","wan_box":"100.72.228.161:60721","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"150","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18367","protocal":"udp(17)","lan":"10.0.0.2:123","wan_box":"100.72.228.161:60720","wan":"17.253.52.125:123","wan_status":"0","timeout":"31","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L.........","ctype":"N"}
		,{"id":"18368","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57301","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57301","wan":"[2a00:1450:400d:805::200e]:443","wan_status":"HERGESTELLT","timeout":"857","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"18373","protocal":"udp(17)","lan":"[2001:870:10:ad4e::a:1]:31218","wan_box":"[2001:870:10:ad4e::a:1]:31218","wan":"[2001:850:3040:1:0:a1:4:a11]:53","wan_status":"0","timeout":"116","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18377","protocal":"udp(17)","lan":"[2001:870:10:ad4e::a:1]:31218","wan_box":"[2001:870:10:ad4e::a:1]:31218","wan":"[2001:850:3040::a1]:53","wan_status":"0","timeout":"116","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18379","protocal":"udp(17)","lan":"100.72.228.161:15714","wan_box":"100.72.228.161:15714","wan":"213.33.98.136:53","wan_status":"0","timeout":"29","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18382","protocal":"udp(17)","lan":"100.72.228.161:27492","wan_box":"100.72.228.161:27492","wan":"213.33.98.136:53","wan_status":"0","timeout":"116","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18389","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57252","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57252","wan":"[2a00:1450:400c:c00::bd]:443","wan_status":"HERGESTELLT","timeout":"880","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"18404","protocal":"tcp(6)","lan":"10.0.0.9:57344","wan_box":"100.72.228.161:60729","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"246","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18407","protocal":"udp(17)","lan":"100.72.228.161:5934","wan_box":"100.72.228.161:5934","wan":"195.3.96.67:53","wan_status":"0","timeout":"130","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18413","protocal":"udp(17)","lan":"100.72.228.161:5934","wan_box":"100.72.228.161:5934","wan":"213.33.98.136:53","wan_status":"0","timeout":"130","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18415","protocal":"udp(17)","lan":"100.72.228.161:23613","wan_box":"100.72.228.161:23613","wan":"213.33.98.136:53","wan_status":"0","timeout":"129","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18421","protocal":"tcp(6)","lan":"10.0.0.9:57337","wan_box":"100.72.228.161:60728","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"205","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18427","protocal":"udp(17)","lan":"[2001:870:10:ad4e::a:1]:31272","wan_box":"[2001:870:10:ad4e::a:1]:31272","wan":"[2001:850:3040::a1]:53","wan_status":"0","timeout":"117","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18439","protocal":"tcp(6)","lan":"10.0.0.9:57298","wan_box":"100.72.228.161:60716","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"134","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18443","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57341","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57341","wan":"[2a00:1450:400d:808::200a]:443","wan_status":"HERGESTELLT","timeout":"882","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"18445","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57269","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57269","wan":"[2a00:1450:400c:c08::bd]:443","wan_status":"HERGESTELLT","timeout":"857","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"18459","protocal":"udp(17)","lan":"[2001:870:10:ad4e::a:1]:25133","wan_box":"[2001:870:10:ad4e::a:1]:25133","wan":"[2001:850:3040::a1]:53","wan_status":"0","timeout":"117","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18479","protocal":"udp(17)","lan":"[2001:870:10:ad4e::a:1]:25133","wan_box":"[2001:870:10:ad4e::a:1]:25133","wan":"[2001:850:3040:1:0:a1:4:a11]:53","wan_status":"0","timeout":"117","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18480","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57343","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57343","wan":"[2a00:1450:400d:809::200e]:443","wan_status":"HERGESTELLT","timeout":"891","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"18492","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57249","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57249","wan":"[2a00:1450:400d:802::200a]:443","wan_status":"ZEIT_WARTEN","timeout":"119","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18536","protocal":"udp(17)","lan":"100.72.228.161:27130","wan_box":"100.72.228.161:27130","wan":"213.33.98.136:53","wan_status":"0","timeout":"88","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18538","protocal":"tcp(6)","lan":"10.0.0.9:57336","wan_box":"100.72.228.161:60727","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"204","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18540","protocal":"udp(17)","lan":"100.72.228.161:16747","wan_box":"100.72.228.161:16747","wan":"195.3.96.67:53","wan_status":"0","timeout":"99","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18544","protocal":"udp(17)","lan":"[2001:870:10:ad4e::a:1]:23413","wan_box":"[2001:870:10:ad4e::a:1]:23413","wan":"[2001:850:3040::a1]:53","wan_status":"0","timeout":"129","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18555","protocal":"tcp(6)","lan":"10.0.0.9:57271","wan_box":"100.72.228.161:60713","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"98","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18562","protocal":"udp(17)","lan":"100.72.228.161:27130","wan_box":"100.72.228.161:27130","wan":"195.3.96.67:53","wan_status":"0","timeout":"88","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18565","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57233","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57233","wan":"[2a00:1450:400c:c00::bd]:443","wan_status":"HERGESTELLT","timeout":"880","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"18566","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57234","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57234","wan":"[2a00:1450:400c:c09::bd]:443","wan_status":"ZEIT_WARTEN","timeout":"105","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18576","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57342","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57342","wan":"[2a00:1450:400c:c0a::bd]:443","wan_status":"HERGESTELLT","timeout":"891","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"18608","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57339","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57339","wan":"[2a00:1450:400d:805::200e]:443","wan_status":"HERGESTELLT","timeout":"881","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"18624","protocal":"udp(17)","lan":"100.72.228.161:28395","wan_box":"100.72.228.161:28395","wan":"213.33.98.136:53","wan_status":"0","timeout":"113","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18628","protocal":"udp(17)","lan":"100.72.228.161:16747","wan_box":"100.72.228.161:16747","wan":"213.33.98.136:53","wan_status":"0","timeout":"98","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18646","protocal":"tcp(6)","lan":"10.0.0.9:57268","wan_box":"100.72.228.161:60711","wan":"129.27.124.196:443","wan_status":"ZEIT_WARTEN","timeout":"96","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18659","protocal":"udp(17)","lan":"10.0.0.2:123","wan_box":"100.72.228.161:60718","wan":"17.253.54.125:123","wan_status":"0","timeout":"28","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L.........","ctype":"N"}
		,{"id":"18662","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57254","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57254","wan":"[2a00:1450:400d:804::200a]:443","wan_status":"ZEIT_WARTEN","timeout":"148","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18669","protocal":"udp(17)","lan":"100.72.228.161:13236","wan_box":"100.72.228.161:13236","wan":"195.3.96.67:53","wan_status":"0","timeout":"28","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18719","protocal":"tcp(6)","lan":"10.0.0.9:57270","wan_box":"100.72.228.161:60712","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"92","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18722","protocal":"udp(17)","lan":"100.72.228.161:27492","wan_box":"100.72.228.161:27492","wan":"195.3.96.67:53","wan_status":"0","timeout":"115","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18731","protocal":"udp(17)","lan":"100.72.228.161:12372","wan_box":"100.72.228.161:12372","wan":"213.33.98.136:53","wan_status":"0","timeout":"28","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18736","protocal":"tcp(6)","lan":"10.0.0.9:57283","wan_box":"100.72.228.161:60715","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"117","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18741","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57340","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57340","wan":"[2a00:1450:400d:808::200e]:443","wan_status":"HERGESTELLT","timeout":"891","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"18763","protocal":"udp(17)","lan":"100.72.228.161:13236","wan_box":"100.72.228.161:13236","wan":"213.33.98.136:53","wan_status":"0","timeout":"28","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18765","protocal":"udp(17)","lan":"100.72.228.161:15341","wan_box":"100.72.228.161:15341","wan":"195.3.96.67:53","wan_status":"0","timeout":"28","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18773","protocal":"tcp(6)","lan":"10.0.0.9:57278","wan_box":"100.72.228.161:60714","wan":"40.77.226.250:443","wan_status":"LETZTE_ACK","timeout":"110","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18779","protocal":"udp(17)","lan":"100.72.228.161:32231","wan_box":"100.72.228.161:32231","wan":"213.33.98.136:53","wan_status":"0","timeout":"1","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18795","protocal":"udp(17)","lan":"100.72.228.161:23613","wan_box":"100.72.228.161:23613","wan":"195.3.96.67:53","wan_status":"0","timeout":"128","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
		,{"id":"18798","protocal":"tcp(6)","lan":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57282","wan_box":"[2001:871:10:ad4e:a950:4521:fcb8:f5e5]:57282","wan":"[2a00:1450:400d:805::200a]:443","wan_status":"HERGESTELLT","timeout":"767","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"..........R..","ctype":"S"}
		,{"id":"18799","protocal":"udp(17)","lan":"100.72.228.161:15341","wan_box":"100.72.228.161:15341","wan":"213.33.98.136:53","wan_status":"0","timeout":"28","tx_bytes":"0/0","tx_packets":"0/0","bind":"","intf":"vINTERNET","route_mode":"","direction":"Outgoing","flags":"I..L......R..","ctype":"N"}
	]
}
*/

type response struct {
	Response []Connection
}

type Connection struct {
	Id         string
	Protocol   string
	Lan        string
	Wan_box    string
	Wan        string
	Wan_status string
	Timeout    string
	Tx_bytes   string
	Tx_packets string
	bind       string
	Intf       string
	Route_mode string
	Direction  string
	Flags      string
	Ctype      string
}

func Get() ([]Connection, error) {

	jsonBlob, err := getJson(json_url)

	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	connections, err := getConnections(jsonBlob)

	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	return connections, err

}

func getConnections(jsonBlob []byte) ([]Connection, error) {

	var resp response

	err := json.Unmarshal(jsonBlob, &resp)

	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	connections := resp.Response

	return connections, nil
}

func getJson(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func printConnTable(conns []Connection) {
	for _, conn := range conns {
		fmt.Printf("%40s %40s \n", conn.Lan, conn.Wan)
	}
}
