package network

import (
   "sync"
   "strconv"
   "net"
   "strings"
   "BlockChainTest/bean"
   "github.com/golang/protobuf/proto"
   "BlockChainTest/crypto"
   "errors"
   "fmt"
)

var onceSender sync.Once
var SenderQueue chan *Message

type IP string

func (ip IP) String() string {
   return string(ip)
}

type Port int

func (port Port) String() string {
   return strconv.Itoa(int(port))
}

type Peer struct {
   IP      IP
   Port    Port
   PubKey  *bean.Point
   Address bean.Address
}

func (peer *Peer) ToString() string {
   return peer.IP.String() + ":" + peer.Port.String()
}

type Message struct {
   Peer *Peer
   Msg  interface{}
}

func identifyMessage(message *Message) (int, interface{}) {
   msg := message.Msg
   switch msg.(type) {
   case *bean.Setup:
      return bean.MsgTypeSetUp, msg.(*bean.Setup)
   case *bean.Commitment:
      return bean.MsgTypeCommitment, msg.(*bean.Commitment)
   case *bean.Challenge:
      return bean.MsgTypeChallenge, msg.(*bean.Challenge)
   case *bean.Response:
      return bean.MsgTypeResponse, msg.(*bean.Response)
   default:
      return -1, nil
   }
}

func GetSenderQueue() chan *Message {
   onceSender.Do(func() {
      SenderQueue = make(chan *Message,  10)
   })
   return SenderQueue
}

func (m *Message) Cipher() ([]byte, error) {
   serializable, err := bean.Serialize(m.Msg)
   if err != nil {
      return nil, err
   }
   sig, err := crypto.Sign(serializable.Body)
   if err != nil {
      return nil, err
   }
   serializable.Sig = sig
   pubKey, err := crypto.GetPubKey()
   if err != nil {
      return nil, err
   }
   serializable.PubKey = pubKey
   plaintext, err := proto.Marshal(serializable)
   if err != nil {
      return nil, err
   }
   cipher, err := crypto.Encrypt(m.Peer.PubKey, plaintext)
   if err != nil {
      return nil, err
   }
   return cipher, nil
}

func (m *Message) Send() error {
   cipher, err := m.Cipher()
   if err != nil {
      return err
   }
   addr, err := net.ResolveTCPAddr("tcp", m.Peer.ToString())
   if err != nil {
     return err
   }
   conn, err := net.DialTCP("tcp", nil, addr)
   if err != nil {
     return err
   }
   defer conn.Close()
   fmt.Println("Send msg:", cipher)
   if _, err := conn.Write(cipher); err != nil {
      return err
   }
   return nil
}

func SendMessage(peers []*Peer, msg interface{}) {
   queue := GetSenderQueue()
   for _, peer := range peers {
      message := &Message{peer, msg}
      queue <- message
   }
}

func DecryptIntoMessage(cipher []byte) (*Message, error) {
   plaintext, err := crypto.Decrypt(cipher)
   if err != nil {
      return nil, err
   }
   serializable, msg, err := bean.Deserialize(plaintext)
   if err != nil {
      return nil, err
   }
   if !crypto.Verify(serializable.Sig, serializable.PubKey, serializable.Body) {
      return nil, errors.New("decrypt fail")
   }
   peer := &Peer{PubKey: serializable.PubKey}
   message := &Message{Peer: peer, Msg: msg}
   return message, nil
}

func startListen(process func(int, interface{})) {
  go func() {
     //room for modification addr := &net.TCPAddr{IP: net.ParseIP("x.x.x.x"), Port: receiver.listeningPort()}
     addr := &net.TCPAddr{Port: listeningPort}
     listener, err := net.ListenTCP("tcp", addr)
     if err != nil {
        fmt.Println("error", err)
        return
     }
     for {
        conn, err := listener.AcceptTCP()
        fmt.Println("listen from ", conn.RemoteAddr())
        if err != nil {
           continue
        }
        b := make([]byte, bufferSize)
        cipher := b
        offset := 0
        for {
           n, err := conn.Read(cipher)
           if err != nil {
              break
           } else {
              cipher = b[n:]
              offset += n
           }
           fmt.Println("Receive before decrypt", cipher)
           message, err := DecryptIntoMessage(cipher)
           fmt.Println("Receive after decrypt", cipher)
           if err != nil {
              return
           }
           fromAddr := conn.RemoteAddr().String()
           ip := fromAddr[:strings.LastIndex(fromAddr, ":")]
           message.Peer.IP = IP(ip)
           //queue := GetReceiverQueue()
           //queue <- message
           //p := processor.GetInstance()
           t, msg := identifyMessage(message)
           if msg != nil {
              process(t, msg)
           }
        }
     }
  }()
}

func startSend() {
   go func() {
      sender := GetSenderQueue()
      for {
         if message, ok := <-sender; ok {
            message.Send()
         }
      }
   }()
}

func Start(process func(int, interface{})) {
   startListen(process)
   startSend()
}