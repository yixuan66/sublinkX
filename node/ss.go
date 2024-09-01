package node

import (
 "fmt"
 "net/url"
 "strconv"
 "strings"
)

// ss匹配规则
type Ss struct {
 Param  Param
 Server string
 Port   int
 Name   string
 Type   string
}

type Param struct {
 Cipher   string
 Password string
}

func parsingSS(s string) (string, string, string) {
 // ss url编码分为三部分：加密方式、服务器地址和端口、备注
 s = strings.Replace(s, "ss://", "", 1)
 addrIndex := strings.Index(s, "@")
 nameIndex := strings.Index(s, "#")
 var addr, name string
 param := s

 if nameIndex != -1 {
  name, _ = url.QueryUnescape(s[nameIndex+1:])
  param = s[:nameIndex]
 }

 if addrIndex != -1 {
  addr = s[addrIndex+1:]
  param = s[:addrIndex]
 } else {
  addr = strings.Split(Base64Decode(param), "@")[1]
  param = strings.Split(Base64Decode(param), "@")[0]
 }
 return param, addr, name
}

// 开发者测试
func CallSSURL() {
 ss := Ss{}
 ss.Name = "测试"
 ss.Server = "rhsdrhwa.cfprefer1.xyz"
 ss.Port = 45611
 ss.Param.Cipher = "aes-256-gcm"
 ss.Param.Password = "254903bd-ac6e-4a26-9ad3-2cf0d96c426c"
 fmt.Println(EncodeSSURL(ss))
}

// ss 编码输出
func EncodeSSURL(s Ss) string {
 //编码格式 ss://base64(base64(method:password)@hostname:port)
 p := Base64Encode(s.Param.Cipher + ":" + s.Param.Password)
 if s.Name == "" {
  s.Name = s.Server + ":" + strconv.Itoa(s.Port)
 }
 param := fmt.Sprintf("%s@%s:%d#%s",
  p,
  s.Server,
  s.Port,
  s.Name,
 )
 return "ss://" + param
}

func DecodeSSURL(s string) (Ss, error) {
 param, addr, name := parsingSS(s)
 decodedParam, err := Base64Decode(param)
 if err != nil {
  return Ss{}, fmt.Errorf("failed to decode base64 param: %v", err)
 }
 param = decodedParam
 if param == "" || addr == "" {
  return Ss{}, fmt.Errorf("invalid SS URL")
 }
 parts := strings.Split(addr, ":")
 port, err := strconv.Atoi(parts[len(parts)-1])
 if err != nil {
  return Ss{}, fmt.Errorf("invalid port in SS URL: %v", err)
 }
 server := strings.Replace(ValRetIPv6Addr(addr), ":"+parts[len(parts)-1], "", -1)
 cipher := strings.Split(param, ":")[0]
 password := strings.Replace(param, cipher+":", "", 1)
 if name == "" {
  name = addr
 }
 if CheckEnvironment() {
  fmt.Println("Param:", param)
  fmt.Println("Server", server)
  fmt.Println("Port", port)
  fmt.Println("Name:", name)
  fmt.Println("Cipher:", cipher)
  fmt.Println("Password:", password)
 }
 return Ss{
  Param: Param{
   Cipher:   cipher,
   Password: password,
  },
  Server: server,
  Port:   port,
  Name:   name,
  Type:   "ss",
 }, nil
}
