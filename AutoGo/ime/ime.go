package ime

import (
	"bufio"
	_ "embed"
	"encoding/binary"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Dasongzi1366/AutoGo/core"
	"github.com/Dasongzi1366/AutoGo/motion"
)

var mutex sync.Mutex
var conn net.Conn

func init() {
	go func() {
		for {
			if conn == nil {
				con, err := net.Dial("unix", "@autogo_ime")
				if err == nil {
					mutex.Lock()
					conn = con
					mutex.Unlock()
				}
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

// GetClipText 获取剪切板内容
func GetClipText() string {
	mutex.Lock()
	defer mutex.Unlock()
	if conn != nil {
		// 发送请求
		err := sendMessage(conn, uint64(time.Now().UnixNano()), 3, []byte{})
		if err != nil {
			return ""
		}

		// 接收回复
		return receiveMessage(conn)
	} else {
		return core.ImeGetClipText()
	}
}

// SetClipText 设置剪切板内容
func SetClipText(text string) bool {
	mutex.Lock()
	defer mutex.Unlock()
	if conn != nil {
		return sendMessage(conn, 0, 2, []byte(text)) == nil
	} else {
		return core.ImeSetClipText(text)
	}
}

// KeyAction 模拟按键,依赖魔改的QQ拼音输入法
func KeyAction(code int) {
	mutex.Lock()
	defer mutex.Unlock()
	if conn != nil {
		sendMessage(conn, 0, 4, []byte(core.I2s(code)))
	} else {
		//core.ImeSetClipText(text)
	}
}

// InputText 输入文本
func InputText(text string, displayId int) {
	mutex.Lock()
	defer mutex.Unlock()
	if conn != nil {
		sendMessage(conn, 0, 1, []byte(text))
	} else {
		core.ImeSetClipText(text)
		motion.KeyAction(279, displayId)
		time.Sleep(200 * time.Millisecond)
		core.ImeSetClipText("")
	}
}

// GetIMEList 获取输入法列表
func GetIMEList() []string {
	shellOutput := core.Shell("ime list -a | grep mId")

	var imeList []string

	scanner := bufio.NewScanner(strings.NewReader(shellOutput))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "mId=") {
			mId := strings.SplitN(line, " ", 2)[0]
			mId = strings.TrimPrefix(mId, "mId=")
			imeList = append(imeList, mId)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil
	}

	return imeList
}

// SetCurrentIME 设置当前输入法
func SetCurrentIME(packageName string) {
	mutex.Lock()
	defer mutex.Unlock()
	output := core.Shell("dumpsys package " + packageName)
	lines := strings.Split(output, "\n")
	re := regexp.MustCompile(packageName + `/\S+`)
	for i, line := range lines {
		if strings.Contains(lines[i], "android.permission.BIND_INPUT_METHOD") && strings.Contains(lines[i], packageName+"/") {
			imeID := re.FindString(line)
			if imeID != "" {
				core.Shell("ime enable " + imeID)
				core.Shell("ime set " + imeID)
				break
			}
		}
	}
	if conn == nil {
		conn, _ = net.Dial("unix", "@autogo_ime")
	}
}

// sendMessage 发送消息到IME服务
// 格式: [8字节长度][8字节消息ID][1字节类型][消息本体]
func sendMessage(conn net.Conn, msgID uint64, msgType byte, data []byte) error {
	// 1. 发送8字节长度（消息本体的长度）
	lengthBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(lengthBuffer, uint64(len(data)))
	_, err := conn.Write(lengthBuffer)
	if err != nil {
		closeConn(conn)
		return err
	}

	// 2. 发送8字节消息ID
	msgIDBuffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(msgIDBuffer, msgID)
	_, err = conn.Write(msgIDBuffer)
	if err != nil {
		closeConn(conn)
		return err
	}

	// 3. 发送1字节消息类型
	_, err = conn.Write([]byte{msgType})
	if err != nil {
		closeConn(conn)
		return err
	}

	// 4. 发送消息本体
	_, err = conn.Write(data)
	if err != nil {
		closeConn(conn)
		return err
	}

	return nil
}

// receiveMessage 接收IME服务的回复消息
func receiveMessage(conn net.Conn) string {
	// 设置读取超时
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	defer conn.SetReadDeadline(time.Time{}) // 清除超时

	lengthBuffer := make([]byte, 8)
	msgIDBuffer := make([]byte, 8)
	msgTypeBuffer := make([]byte, 1)

	// 1. 读取8字节长度
	_, err := conn.Read(lengthBuffer)
	if err != nil {
		closeConn(conn)
		return ""
	}
	msgLen := binary.LittleEndian.Uint64(lengthBuffer)

	// 2. 读取8字节消息ID
	_, err = conn.Read(msgIDBuffer)
	if err != nil {
		closeConn(conn)
		return ""
	}

	// 3. 读取1字节消息类型
	_, err = conn.Read(msgTypeBuffer)
	if err != nil {
		closeConn(conn)
		return ""
	}

	// 4. 读取消息本体
	msgBuffer := make([]byte, msgLen)
	_, err = conn.Read(msgBuffer)
	if err != nil {
		closeConn(conn)
		return ""
	}

	return string(msgBuffer)
}

func closeConn(c net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	if c != nil {
		c.Close()
	}
	conn = nil
}
