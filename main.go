package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func exit(err *error) {
	if *err != nil {
		log.Println("exited with error:", (*err).Error())
		os.Exit(1)
	} else {
		log.Println("exited")
	}
}

func decodeModifies(env []string, pfx string) map[string]map[string]string {
	modifies := map[string]map[string]string{}
	for _, item := range env {
		skvRaw := strings.SplitN(item, "=", 2)
		if len(skvRaw) != 2 {
			continue
		}
		s, kvRaw := skvRaw[0], skvRaw[1]
		if !strings.HasPrefix(s, pfx) {
			continue
		}
		s = strings.TrimPrefix(s, pfx)
		if li := strings.LastIndex(s, "__"); li >= 0 {
			s = s[0:li]
		}
		if modifies[s] == nil {
			modifies[s] = map[string]string{}
		}
		kv := strings.SplitN(kvRaw, "=", 2)
		if len(kv) != 2 {
			continue
		}
		k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		modifies[s][k] = v
	}
	return modifies
}

func applyModifies(modifies map[string]map[string]string, lines [][]byte) [][]byte {
	var (
		i    int
		line []byte

		s        []byte
		replaced = map[string]bool{}

		isComment bool
	)

	for {
		if i >= len(lines) {
			break
		}

		line = lines[i]

		line = bytes.TrimSpace(line)

		// 分段
		if bytes.HasPrefix(line, []byte("[")) && bytes.HasSuffix(line, []byte("]")) {
			// 结束当前分段，把所有未处理的值写入分段内
			kv := modifies[string(s)]
			delete(modifies, string(s))
			for k, v := range kv {
				log.Printf("Append: [%s] %s = %s", s, k, v)
				lines = append(append(lines[0:i], []byte(fmt.Sprintf("%s = %s", k, v))), lines[i:]...)
				i++
			}
			// 开始新分段
			s = bytes.TrimSpace(bytes.TrimSuffix(bytes.TrimPrefix(line, []byte("[")), []byte("]")))
			replaced = map[string]bool{}
			i++
			continue
		}

		// 分段不作处理
		if len(modifies[string(s)]) == 0 && len(replaced) == 0 {
			i++
			continue
		}

		// 移除注释
		if bytes.HasPrefix(line, []byte(";")) {
			line = bytes.TrimLeft(line, ";")
			isComment = true
		} else {
			isComment = false
		}
		currentSplits := bytes.SplitN(line, []byte("="), 2)
		// 当前行不包含 k = v
		if len(currentSplits) != 2 {
			i++
			continue
		}
		currentK := bytes.TrimSpace(currentSplits[0])
		// 检查键值是否已经写入，则注释掉该行
		if replaced[string(currentK)] {
			if !isComment {
				log.Printf("Comment: [%s] #%d: %s", s, i, line)
				lines[i] = append([]byte("; "), line...)
			}
			i++
			continue
		}
		// 尝试获取 kv
		if v, ok := modifies[string(s)][string(currentK)]; ok {
			// 找到了 k 相同的值，替换当前行
			log.Printf("Replace: [%s] %s = %s", s, currentK, v)
			lines[i] = []byte(fmt.Sprintf("%s = %s", currentK, v))
			replaced[string(currentK)] = true
			delete(modifies[string(s)], string(currentK))
			i++
			continue
		}

		i++
	}

	// 剩余的 k v
	kv := modifies[string(s)]
	delete(modifies, string(s))
	for k, v := range kv {
		log.Printf("Append: [%s] %s = %s", s, k, v)
		lines = append(lines, []byte(fmt.Sprintf("%s = %s", k, v)))
	}

	// 额外的分区
	for s, kv := range modifies {
		log.Printf("Append Section: [%s]", s)
		if s == "" {
			for k, v := range kv {
				log.Printf("Append: %s = %s", k, v)
				lines = append([][]byte{[]byte(fmt.Sprintf("%s = %s", k, v))}, lines...)
			}
		} else {
			lines = append(lines, []byte(fmt.Sprintf("[%s]", s)))
			for k, v := range kv {
				log.Printf("Append: %s = %s", k, v)
				lines = append(lines, []byte(fmt.Sprintf("%s = %s", k, v)))
			}
		}
	}

	return lines
}

func main() {
	var err error
	defer exit(&err)

	var (
		optFrom string
		optTo   string
	)

	flag.StringVar(&optFrom, "from", "", "environment variable prefix")
	flag.StringVar(&optTo, "to", "", "file to modify")
	flag.Parse()

	if optFrom == "" {
		err = errors.New("missing --from")
		return
	}

	if optTo == "" {
		err = errors.New("missing --to")
		return
	}

	modifies := decodeModifies(os.Environ(), optFrom)

	var buf []byte
	if buf, err = os.ReadFile(optTo); err != nil {
		return
	}

	lines := bytes.Split(buf, []byte("\n"))

	lines = applyModifies(modifies, lines)

	buf = bytes.Join(lines, []byte("\n"))

	if err = os.WriteFile(optTo, buf, 0755); err != nil {
		return
	}
}
