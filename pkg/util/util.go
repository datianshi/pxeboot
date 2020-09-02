package util

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/datianshi/pxeboot/pkg/config"
	"github.com/datianshi/pxeboot/pkg/model"
	"github.com/oklog/ulid"
	"io"
	"regexp"
	"strings"
	"time"
)
const SYMLINK_PER_SERVER_DIR = "images"
const KERNEL_FILE = "b.b00"
func Dot_To_Dash(dot string) string{
	return strings.ReplaceAll(dot, ".", "-")
}

func Dash_To_Dot(dot string) string{
	return strings.ReplaceAll(dot, "-", ".")
}

func Colon_To_Dash(dot string) string{
	return strings.ReplaceAll(dot, ":", "-")
}

func BootConfigFile(src io.Reader, dst io.Writer, bindIP, nic string) {
	kickstartUrl := fmt.Sprintf("http://%s/kickstart/%s/ks.cfg", bindIP, nic)
	prefix := fmt.Sprintf("01-%s/%s", nic, SYMLINK_PER_SERVER_DIR)
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "kernelopt=") {
			re := regexp.MustCompile(`^(kernelopt=).*$`)
			t := re.ReplaceAllString(line, fmt.Sprintf("${1}ks=%s", kickstartUrl))
			_, err := dst.Write([]byte(fmt.Sprintf("%s\n", t)))
			if err != nil {
				fmt.Println(err)
			}

		} else if strings.Contains(line, "prefix=") {
			re := regexp.MustCompile(`^(prefix=).*$`)
			t := re.ReplaceAllString(line, fmt.Sprintf("${1}%s", prefix))
			//fmt.Println(t)
			dst.Write([]byte(fmt.Sprintf("%s\n", t)))
		} else {
			dst.Write([]byte(fmt.Sprintf("%s\n", line)))
		}

	}
}

func ConvertToServerItems(nics map[string]config.ServerConfig) []model.ServerItem {
	var items []model.ServerItem
	for k,v := range nics {
		item := model.ServerItem{
			v.Ip,
			v.Hostname,
			k,
			v.Gateway,
			v.Netmask,
		}
		items = append(items, item)
	}
	return items
}

func UUID() (string, error) {
	t := time.Unix(1000000, 0)
	entropy := make([]byte, 48)
	if _, err := io.ReadFull(rand.Reader, entropy); err != nil {
		return "", err
	}
	uuid := ulid.MustNew(ulid.Timestamp(t), bytes.NewReader(entropy[:])).String()
	return uuid, nil
}