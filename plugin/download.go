package plugin

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type DownloadProgress struct {
	Length  int64
	Current int64
	Status  chan string
}

func (d *DownloadProgress) Write(p []byte) (int, error) {
	n := len(p)
	d.Current += int64(n)
	d.Status <- fmt.Sprintf("正在下载：%d%%", d.Current*100/d.Length)

	return n, nil
}

func DownloadPlugin(name string, status chan string) error {
	name += ".so"
	file := "plugins/" + name

	status <- "正在连接"

	resp, err := http.Get("https://xiayesuifeng.gitlab.io/goblog-plugins/" + name)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("连接失败")
	}

	length, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		return err
	}

	status <- "正在下载"

	out, err := os.Create(file + ".temp")
	if err != nil {
		return err
	}
	defer out.Close()

	dp := &DownloadProgress{Length: length, Status: status}
	_, err = io.Copy(out, io.TeeReader(resp.Body, dp))
	if err != nil {
		return err
	}

	err = os.Rename(file+".temp", file)
	if err != nil {
		return err
	}

	status <- "下载完成"

	return nil
}
