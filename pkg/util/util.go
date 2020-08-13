package util

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
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

func BootConfigFile(src io.Reader, dst io.Writer, kickstartUrl string) {
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "kernelopt=") {
			re := regexp.MustCompile(`^(kernelopt=).*$`)
			t := re.ReplaceAllString(line, fmt.Sprintf("${1}%s", kickstartUrl))
			_, err := dst.Write([]byte(fmt.Sprintf("%s\n", t)))
			if err != nil {
				fmt.Println(err)
			}

		} else if strings.Contains(line, "prefix=") {
			re := regexp.MustCompile(`^(prefix=).*$`)
			t := re.ReplaceAllString(line, fmt.Sprintf("${1}%s", SYMLINK_PER_SERVER_DIR))
			//fmt.Println(t)
			dst.Write([]byte(fmt.Sprintf("%s\n", t)))
		} else if strings.Contains(line, "kernel=") {
			re := regexp.MustCompile(`^(kernel=).*$`)
			t := re.ReplaceAllString(line, fmt.Sprintf("${1}%s", KERNEL_FILE))
			//fmt.Println(t)
			dst.Write([]byte(fmt.Sprintf("%s\n", t)))
		} else {
			dst.Write([]byte(fmt.Sprintf("%s\n", line)))
		}

	}
}
