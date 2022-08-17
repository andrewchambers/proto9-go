package proto9

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"os/user"
	"syscall"
	"testing"
	"time"
)

func GetFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

type DiodTestServer struct {
	t             *testing.T
	Aname         string
	Uname         string
	ServeDir      string
	ListenAddress string
	Diod          *exec.Cmd
}

func (s *DiodTestServer) Dial() net.Conn {
	c, err := net.Dial("tcp", s.ListenAddress)
	if err != nil {
		s.t.Fatal(err)
	}
	s.t.Cleanup(func() { c.Close() })
	return c
}

// Create a test server that is automatically cleaned up when
// the test finishes.
func NewDiodTestServer(t *testing.T) *DiodTestServer {

	_, err := exec.LookPath("diod")
	if err != nil {
		t.Skip("diod not found in path")
	}

	port, err := GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	listenAddress := fmt.Sprintf("127.0.0.1:%d", port)

	currentUser, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	dir := t.TempDir()

	diodOpts := []string{
		"-f",
		"-l", listenAddress,
		"-e", dir,
		"-d", "1",
		"-n",
		"-U", currentUser.Username,
	}

	diod := exec.Command(
		"diod",
		diodOpts...,
	)
	diod.Stderr = os.Stderr

	err = diod.Start()
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_ = diod.Process.Signal(syscall.SIGTERM)
		_, _ = diod.Process.Wait()
	})

	t.Logf("starting diod %v", diodOpts)

	up := false
	for i := 0; i < 2000; i++ {
		c, err := net.Dial("tcp", listenAddress)
		if err == nil {
			up = true
			_ = c.Close()
			break
		}
		time.Sleep(1 * time.Millisecond)
	}
	if !up {
		t.Fatal("diod server never came up")
	}

	return &DiodTestServer{
		t:             t,
		ListenAddress: listenAddress,
		Aname:         dir,
		Uname:         currentUser.Username,
		ServeDir:      dir,
		Diod:          diod,
	}
}

func NewTestDotLClient(t *testing.T) (*Client, *DiodTestServer) {
	server := NewDiodTestServer(t)
	client, err := NewClient(server.Dial(), "9P2000.L", 1024*1024)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = client.Close()
	})
	return client, server
}

func TestClientConnect(t *testing.T) {
	client, _ := NewTestDotLClient(t)
	err := client.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDotlAttach(t *testing.T) {
	client, server := NewTestDotLClient(t)
	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()
}

func TestDotlEmptyWalk(t *testing.T) {
	client, server := NewTestDotLClient(t)
	f1, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}

	f2, _, err := f1.Walk([]string{})
	if err != nil {
		t.Fatal()
	}

	err = f1.Clunk()
	if err != nil {
		t.Fatal()
	}

	err = f2.Clunk()
	if err != nil {
		t.Fatal()
	}
}

func TestDotlWalkOne(t *testing.T) {
	client, server := NewTestDotLClient(t)

	dir := server.ServeDir

	err := os.MkdirAll(dir+"/1/2/3", 0o777)
	if err != nil {
		t.Fatal(err)
	}

	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()

	wf, _, err := f.Walk([]string{"1", "2", "3"})
	if err != nil {
		t.Fatal(err)
	}
	err = wf.Clunk()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDotlWalkMulti(t *testing.T) {
	client, server := NewTestDotLClient(t)

	dir := server.ServeDir

	err := os.MkdirAll(dir+"/1/2/3/4/5/6/7/8/9/10/11/12/13/14", 0o777)
	if err != nil {
		t.Fatal(err)
	}

	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()

	wf, _, err := f.Walk(
		[]string{
			"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	err = wf.Clunk()
	if err != nil {
		t.Fatal(err)
	}

	err = f.Clunk()
	if err != nil {
		t.Fatal(err)
	}
}

func TestDotlShortWalk(t *testing.T) {
	client, server := NewTestDotLClient(t)

	dir := server.ServeDir

	err := os.MkdirAll(dir+"/1/2/3/4", 0o777)
	if err != nil {
		t.Fatal(err)
	}

	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()

	_, qids, err := f.Walk(
		[]string{
			"1", "2", "x", "3", "4",
		},
	)
	if err != ErrShortWalk {
		t.Fatal(err)
	}
	if len(qids) != 2 {
		t.Fatal("unexpected qid count")
	}
}

func TestDotlRemove(t *testing.T) {
	client, server := NewTestDotLClient(t)

	dir := server.ServeDir

	err := os.Mkdir(dir+"/x", 0o777)
	if err != nil {
		t.Fatal(err)
	}

	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()

	wf, _, err := f.Walk([]string{"x"})
	if err != nil {
		t.Fatal(err)
	}
	err = wf.Remove()
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(dir + "/x")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatal(err)
	}
}
