# stun


## STUN server 

 In the Binding request/response transaction, a Binding request is sent from a STUN client to a STUN server. When the Binding request arrives at the STUN server, it may have passed through one or more NATs between the STUN client and the STUN server.  As the Binding request message passes through a NAT, the NAT will modify the source transport address (that is, the source IP address and the source port) of the packet.  As a result, the source transport address of the request received by the server will be the public IP address and port created by the NAT closest to the server.  This is called a reflexive transport address.  The STUN server copies that source transport address into an XOR-MAPPED- ADDRESS attribute in the STUN Binding response and sends the Binding response back to the STUN client.  As this packet passes back through a NAT, the NAT will modify the destination transport address in the IP header, but the transport address in the XOR-MAPPED-ADDRESS attribute within the body of the STUN response will remain untouched. In this way, the client can learn its reflexive transport address allocated by the outermost NAT with respect to the STUN server.

## STUN message structure

 All STUN messages MUST start with a 20-byte header followed by zero
   or more Attributes.  The STUN header contains a STUN message type,
   magic cookie, transaction ID, and message length.

       0                   1                   2                   3
       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |0 0|     STUN Message Type     |         Message Length        |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                         Magic Cookie                          |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                                                               |
      |                     Transaction ID (96 bits)                  |
      |                                                               |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+




0                 1
                        2  3  4 5 6 7 8 9 0 1 2 3 4 5

                       +--+--+-+-+-+-+-+-+-+-+-+-+-+-+
                       |M |M |M|M|M|C|M|M|M|C|M|M|M|M|
                       |11|10|9|8|7|1|6|5|4|0|3|2|1|0|
                       +--+--+-+-+-+-+-+-+-+-+-+-+-+-+





The magic cookie field MUST contain the fixed value 0x2112A442 in
   network byte order. 


## Reference

- [https://tools.ietf.org/html/rfc5389](https://tools.ietf.org/html/rfc5389)
- [http://www.can.mydns.jp/~ryo/rfc/rfc5389-ja.txt](http://www.can.mydns.jp/~ryo/rfc/rfc5389-ja.txt)
- [https://qiita.com/okyk/items/a405f827e23cb9ef3bde#_reference-7b63e074c87bb20ea01f](https://qiita.com/okyk/items/a405f827e23cb9ef3bde#_reference-7b63e074c87bb20ea01f)