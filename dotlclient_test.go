package proto9

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"os/user"
	"reflect"
	"sync"
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

	rpipe, wpipe, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	logWg := &sync.WaitGroup{}
	logWg.Add(1)
	go func() {
		defer logWg.Done()
		brdr := bufio.NewReader(rpipe)
		for {
			line, err := brdr.ReadString('\n')
			if err != nil {
				return
			}
			if len(line) == 0 {
				continue
			}
			t.Log(line[:len(line)-1])
		}
	}()

	t.Cleanup(func() {
		logWg.Wait()
	})

	diod.Stderr = wpipe
	diod.Stdout = wpipe

	err = diod.Start()
	if err != nil {
		t.Fatal(err)
	}
	_ = wpipe.Close()

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
	client, err := NewClient(server.Dial(), "9P2000.L", 4096)
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

func TestDotLAttach(t *testing.T) {
	client, server := NewTestDotLClient(t)
	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()
}

func TestDotLEmptyWalk(t *testing.T) {
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

func TestDotLWalkOne(t *testing.T) {
	client, server := NewTestDotLClient(t)

	err := os.MkdirAll(server.ServeDir+"/1/2/3", 0o777)
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

func TestDotLWalkMulti(t *testing.T) {
	client, server := NewTestDotLClient(t)

	err := os.MkdirAll(server.ServeDir+"/1/2/3/4/5/6/7/8/9/10/11/12/13/14", 0o777)
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

func TestDotLShortWalk(t *testing.T) {
	client, server := NewTestDotLClient(t)

	err := os.MkdirAll(server.ServeDir+"/1/2/3/4", 0o777)
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

func TestDotLRemove(t *testing.T) {
	client, server := NewTestDotLClient(t)

	err := os.Mkdir(server.ServeDir+"/x", 0o777)
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
	_, err = os.Stat(server.ServeDir + "/x")
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatal(err)
	}
}

func TestDotLRead(t *testing.T) {
	client, server := NewTestDotLClient(t)

	expected, err := io.ReadAll(
		&io.LimitedReader{R: rand.Reader, N: int64(2 * client.Msize())},
	)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(server.ServeDir+"/x", expected, 0o777)
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
	defer wf.Clunk()

	err = wf.Open(L_O_RDONLY)
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, len(expected), len(expected))
	n, err := wf.Read(0, buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != (client.Msize() - IOHDRSZ) {
		t.Fatalf("unexpected read count %d", n)
	}

	if !reflect.DeepEqual(buf[:n], expected[:n]) {
		t.Fatalf("%v\n!=\n%v", buf[:n], expected[:n])
	}
}

func TestDotLWrite(t *testing.T) {
	client, server := NewTestDotLClient(t)

	err := os.WriteFile(server.ServeDir+"/x", []byte("hello"), 0o777)
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
	defer wf.Clunk()

	err = wf.Open(L_O_TRUNC | L_O_WRONLY)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := io.ReadAll(
		&io.LimitedReader{R: rand.Reader, N: int64(2 * client.Msize())},
	)
	if err != nil {
		t.Fatal(err)
	}

	n, err := wf.Write(0, expected)
	if err != nil {
		t.Fatal(err)
	}
	if n != (client.Msize() - IOHDRSZ) {
		t.Fatalf("unexpected write count %d", n)
	}

	buf, err := os.ReadFile(server.ServeDir + "/x")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(buf[:n], expected[:n]) {
		t.Fatalf("%v\n!=\n%v", buf[:n], expected[:n])
	}
}

func TestDotLCreate(t *testing.T) {
	client, server := NewTestDotLClient(t)

	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()
	_, _, err = f.Create("x", 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(server.ServeDir + "/x")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDotLGetAttr(t *testing.T) {
	client, server := NewTestDotLClient(t)

	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()

	_, err = f.GetAttr(L_GETATTR_ALL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDotLSetAttr(t *testing.T) {
	client, server := NewTestDotLClient(t)

	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()

	err = f.SetAttr(LSetAttr{
		Valid: L_SETATTR_MODE,
		Mode:  0o777,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDotLRename(t *testing.T) {
	client, server := NewTestDotLClient(t)

	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()

	xf, _, err := f.Walk([]string{})
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = xf.Create("x", 0, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	// XXX
	// https://github.com/chaos/diod/issues/93
	// When this issue is addressed the enclosed section can be deleted.
	err = xf.Clunk()
	if err != nil {
		t.Fatal(err)
	}

	xf, _, err = f.Walk([]string{"x"})
	if err != nil {
		t.Fatal(err)
	}
	// END

	err = xf.Rename(f, "y")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(server.ServeDir + "/y")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDotLMkdir(t *testing.T) {
	client, server := NewTestDotLClient(t)

	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()

	_, err = f.Mkdir("x", 0o777, 0)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(server.ServeDir + "/x")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDotLStatfs(t *testing.T) {
	client, server := NewTestDotLClient(t)

	f, err := AttachDotL(client, server.Aname, server.Uname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Clunk()

	_, err = f.Statfs()
	if err != nil {
		t.Fatal(err)
	}
}
