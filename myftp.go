package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type ServerOpt struct {
	port int
	addr string
	auth string
	pass string
	filepath string

}

type client struct {
	auth string
	pass string
	filepath string
	passiveport int
	newportlistener net.Listener
	newportconn net.Conn
	portinuse bool
	quitflag bool
}

func Default_init(opt *ServerOpt,args []string) (*ServerOpt,error){
	if opt == nil {
		return nil,errors.New("server input error")
	}else {
		if len(args) > 1{
			for i:= 1; i < len(args) ; i++{
				operation := args[i]
				switch operation[0:2] {
				case "-h":
					if len(operation) > 2{
						opt.addr = operation[3:]
					}else{
						help()
						return nil,nil
					}
					break
				case "-p":
					var err error
					opt.port,err = strconv.Atoi(operation[3:])
					if err!=nil{
						fmt.Println("args error in -p")
						help()
						return nil,errors.New("args error")
					}
					break
				case "-d":
					if accessable(operation[3:]){
						opt.filepath = operation[3:]
					}else{
						fmt.Println("args error in -d")
						help()
						return nil,errors.New("args error")
					}
					break
				}
			}
		}

		var newserver ServerOpt

		if opt.addr == ""{
			newserver.addr = "127.0.0.1"
		}else {
			newserver.addr	= opt.addr
		}

		if opt.auth == ""{
			newserver.auth = "hanmeng"
		}else{
			newserver.auth = opt.auth
		}

		if opt.filepath == ""{
			newserver.filepath = "/Users"
		}else {
			newserver.filepath = opt.filepath
		}

		if opt.pass == ""{
			newserver.pass = "abc"
		}else{
			newserver.pass = opt.pass
		}

		if opt.port == 0{
			newserver.port = 21
		}else{
			newserver.port = opt.port
		}
		return &newserver,nil
	}
}

func send(con net.Conn,commond string){
	nc := commond + "\r\n"
	b :=[]byte(nc)
	con.Write(b)
	log.Printf("send <%s>\n",commond)
}
func receive(con net.Conn,buf []byte) []string{
	num,err := con.Read(buf)
	if err != nil{
		fmt.Println(err)
		return nil
	}
	s :=string(buf[0:num-2])
	log.Println("receive <" + s + ">")
	return strings.Fields(s)

}

func access(aclt client) bool{
	if aclt.auth == "hanmeng" && aclt.pass == "abc"{
		return true
	}else{
		return false
	}
}

func accessable(path string)bool{
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}


func substr(s string,pos, length int) string{
	runes := []rune(s)
	l :=pos + length
	if l > len(runes){
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(directory string) string{
	return substr(directory,0,strings.LastIndex(directory,"/"))
}

func readdir(path string)string{
	files,_ :=ioutil.ReadDir(path)
	fmt.Println("path:",path)
	var information string
	for index,file := range files{
		information += file.Mode().String()+ "\t" + strconv.Itoa(int(file.Size())) +  "\t" + file.ModTime().Month().String()+"\t"+strconv.Itoa(file.ModTime().Day()) + "\t" + strconv.Itoa(file.ModTime().Hour())+":"+strconv.Itoa(file.ModTime().Minute()) + "\t" + file.Name()
		if index < (len(files) -1){
			information += "\r\n"
		}
	}
	return information
}

func randport()int{
	return rand.Intn(100)*256 + rand.Intn(100)+ 1024
}
func help(){
	fmt.Println("ftpserver [opertion]\r\n\t-p=PORT listening port\r\n\t-h=HOST binding address\r\n\t-d=DIR change current directory\r\n\t-h print help message")
}

func server(conn net.Conn,serv *ServerOpt,wg sync.WaitGroup){
	defer wg.Done()
	buf := make([]byte, 4096)
	send(conn, "220 welcome to hanmeng ftp server")
	var clt client
	var err error
	clt.filepath = serv.filepath
	for {
		if clt.quitflag == true {
			break
		}
		com := receive(conn, buf)
		if len(com) == 1 {
			switch com[0] {
			case "PWD":
				send(conn, "257"+` "`+clt.filepath+`"`+"is current directory")
				break
			case "LIST":
				send(conn, "150 Opening data channel for directory listing of "+`"`+clt.filepath+`"`)
				var dirinformation= readdir(clt.filepath)
				send(clt.newportconn, dirinformation)
				clt.newportconn.Close()
				clt.newportlistener.Close()
				clt.portinuse = false
				send(conn, "226 Successfully transferred "+`"`+clt.filepath+`"`)
				break
			case "PASV":
				var newport= randport()
				if clt.portinuse == true {
					clt.newportconn.Close()
					clt.newportlistener.Close()
					clt.portinuse = false
				}
				clt.newportlistener, err = net.Listen("tcp", net.JoinHostPort(serv.addr, strconv.Itoa(newport)))
				if err != nil {
					fmt.Println(err)
					return
				} else {
					log.Println("successful open a port")
				}
				send(conn, "227"+" Entering passive mode (127,0,0,1,"+strconv.Itoa(newport/256)+","+strconv.Itoa(newport%256)+")")
				clt.newportconn, err = clt.newportlistener.Accept()
				clt.portinuse = true
				//Accept a request()
				if err != nil {
					fmt.Println(err)
				}
				log.Println("new port:", newport, "connected")
				break
			case "QUIT":
				send(conn, "221 good bye")
				conn.Close()
				clt.quitflag = true
				//listener.Close()
				break
			}
		}
		if len(com) == 2 {
			switch com[0] {
			case "USER":
				clt.auth = com[1]
				send(conn, "331")
				break
			case "PASS":
				clt.pass = com[1]
				if access(clt) {
					send(conn, "230 login success")
				} else {
					send(conn, "530 login fail")
				}
				break
			case "CWD":
				var newfilepath = clt.filepath
				if com[1] == ".." {
					newfilepath = getParentDirectory(newfilepath)
				} else {
					newfilepath += "/" + com[1]
				}
				if accessable(newfilepath) {
					clt.filepath = newfilepath
					send(conn, "250 CWD successful. "+`"`+clt.filepath+`"`+"is current directory")
				} else {
					send(conn, "550 CWD failed. path:"+newfilepath+"is not legal path")
				}
				break
			case "RETR":
				send(conn, "150 Opening data channel for file download from server of "+`"`+com[1]+`"`)
				var filename = clt.filepath + "/" + com[1]
				var filecode *os.File
				if filecode, err = os.Open(filename); err != nil {
					log.Println("file open error:",err)
					return
				}
				if _, err = io.Copy(clt.newportconn, filecode); err != nil {
					log.Println("transferred file error:",err)
					return
				}
				filecode.Close()
				clt.newportconn.Close()
				clt.newportlistener.Close()
				clt.portinuse = false
				send(conn, "226 Successfully transferred "+`"`+com[1]+`"`)
				break
			case "STOR":
				send(conn, "150 Opening data channel for file upload to server of "+`"`+com[1]+`"`)
				//recv file name
				fo, err := os.Create(clt.filepath + "/" + com[1])
				if err != nil {
					log.Println("Create file error:" + err.Error())
					return
				}
				if _, err = io.Copy(fo, clt.newportconn); err != nil {
					log.Println("transferred file error:",err)
					return
				}
				fo.Close()
				clt.newportconn.Close()
				clt.newportlistener.Close()
				clt.portinuse = false
				send(conn, "226 Successfully transferred "+`"`+com[1]+`"`)
				break
			}
		}
	}
}
func init() {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Llongfile)
}

func main(){
	log.Println("start main")
	var seropt ServerOpt
	var args  = os.Args
	var wg sync.WaitGroup
	log.Println("args:",args)
	//read in the args
	//create the serverOpt
	serv, err := Default_init(&seropt,args)
	if serv == nil{
		log.Fatalln("end")
	}
	if err == nil {
		log.Printf("server:\n\taddr:%s\n\tport:%d\n\tauth:%s\n\tpass:%s\n\tdirpath:%s\n", serv.addr, serv.port, serv.auth, serv.pass, serv.filepath)
	}else {
		log.Fatalln("server init error:",err)
	}
	// create a socket
	log.Println("create a socket")
	listener,err := net.Listen("tcp",net.JoinHostPort(serv.addr, strconv.Itoa(serv.port)))
	if err != nil{
		log.Fatalln("create socket fail:",err)
	}else{
		log.Println("successed create a socket")
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("fail to accept a connect ",err)
			continue
		}
		wg.Add(1)
		log.Println("create another connect,total connect :",wg)
		go server(conn,serv,wg)
	}
	wg.Wait()
	log.Println("end")
}