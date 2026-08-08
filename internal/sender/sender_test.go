package sender

import (
	"bufio"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestSenderUDPDeliversMessage(t *testing.T) {
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenPacket() error: %v", err)
	}
	defer conn.Close()

	host, port := splitHostPort(t, conn.LocalAddr().String())

	s, err := NewSender(UDP, host, port)
	if err != nil {
		t.Fatalf("NewSender() error: %v", err)
	}
	defer s.Close()

	const want = "<134>test message\n"
	if err := s.Send(want); err != nil {
		t.Fatalf("Send() error: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFrom(buf)
	if err != nil {
		t.Fatalf("ReadFrom() error: %v", err)
	}

	if got := string(buf[:n]); got != want {
		t.Errorf("received %q, want %q", got, want)
	}
}

func TestSenderTCPDeliversMessage(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Listen() error: %v", err)
	}
	defer ln.Close()

	host, port := splitHostPort(t, ln.Addr().String())

	received := make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		line, _ := bufio.NewReader(conn).ReadString('\n')
		received <- line
	}()

	s, err := NewSender(TCP, host, port)
	if err != nil {
		t.Fatalf("NewSender() error: %v", err)
	}
	defer s.Close()

	const want = "<134>tcp message\n"
	if err := s.Send(want); err != nil {
		t.Fatalf("Send() error: %v", err)
	}

	select {
	case got := <-received:
		if got != want {
			t.Errorf("received %q, want %q", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for TCP listener to receive message")
	}
}

func TestSenderAppendsMissingNewline(t *testing.T) {
	conn, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("ListenPacket() error: %v", err)
	}
	defer conn.Close()

	host, port := splitHostPort(t, conn.LocalAddr().String())

	s, err := NewSender(UDP, host, port)
	if err != nil {
		t.Fatalf("NewSender() error: %v", err)
	}
	defer s.Close()

	if err := s.Send("no trailing newline"); err != nil {
		t.Fatalf("Send() error: %v", err)
	}

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	buf := make([]byte, 1024)
	n, _, err := conn.ReadFrom(buf)
	if err != nil {
		t.Fatalf("ReadFrom() error: %v", err)
	}

	if got := string(buf[:n]); got != "no trailing newline\n" {
		t.Errorf("received %q, want trailing newline appended", got)
	}
}

func TestNewSenderUnsupportedProtocol(t *testing.T) {
	if _, err := NewSender(Protocol("carrier-pigeon"), "127.0.0.1", 514); err == nil {
		t.Error("NewSender() with unsupported protocol expected error, got nil")
	}
}

func splitHostPort(t *testing.T, addr string) (string, int) {
	t.Helper()
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("SplitHostPort(%q) error: %v", addr, err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("parsing port %q: %v", portStr, err)
	}
	return host, port
}
