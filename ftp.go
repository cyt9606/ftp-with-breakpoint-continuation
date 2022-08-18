package main

import (
	"bytes"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"os"
	"strconv"
	"time"
)

func Link(host string, port string, user string, pwd string) (*ftp.ServerConn, error) {
	c, err := ftp.Dial(host+":"+port, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return nil, err
	}
	c.Login(user, pwd)
	if err != nil {
		return nil, err
	}
	return c, err
}
func CloseLink(c *ftp.ServerConn) error {
	err := c.Quit()
	return err
}

func download(c *ftp.ServerConn, src string, dst string, buffer int) {
	temp := dst + ".temp"
	var index int64
	if _, err := os.Stat(temp); err != nil {
		index = 0
	} else if err == nil {
		data, err := os.ReadFile(temp)
		if err != nil {
			panic(err)
		}
		index, _ = strconv.ParseInt(string(data), 10, 64)
	}
	reader, err := c.RetrFrom(src, uint64(index))
	if err != nil {
		panic(err)
	}
	dstFile, _ := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer dstFile.Close()
	dstFile.Seek(index, 0)
	buf := make([]byte, buffer, buffer)
	var n2 int
	total := int(index)
	for {
		n1, err := reader.Read(buf)
		if err == io.EOF {
			fmt.Println("done")
			os.Remove(temp)
			break
		}
		n2, _ = dstFile.Write(buf[:n1])
		total += n2
		err = os.WriteFile(temp, []byte(strconv.Itoa(total)), 0666)
		if err != nil {
			panic(err)
		}
	}
}
func upload(c *ftp.ServerConn, src string, dst string) {
	index := uint64(0)
	dsttmp := dst + ".tmp"
	if target, _ := c.List(dsttmp); target != nil {
		index = target[0].Size
	}
	srcFile, err := os.Open(src)
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()
	srcFile.Seek(int64(index), 0)
	buf := make([]byte, 4096, 4096)
	for {
		n, err := srcFile.Read(buf)
		if err == io.EOF {
			fmt.Println("done")
			c.Rename(dsttmp, dst)
			break
		}
		data := bytes.NewReader(buf)
		c.StorFrom(dsttmp, data, index)
		index += uint64(n)
	}
}

//func upload(c *ftp.ServerConn, src string, dst string) {
//
//}
