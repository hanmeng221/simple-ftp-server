# simple ftp server
use golang to build a ftp server

This is a interview test of Sensetime. It will:

不能抄代码，可以参考开源软件的实现
开发机: Linux/MacOSX
语言版本: Golang 1.8以上
用纯Golang实现，基础功能只能使用标准库提供的功能，扩展功能可视需求使用第三方库
尽量保证你的实现的安全性和稳定性，以产品代码的质量要求自己
评分标准
遵循RFC 765标准 使用标准客户端ftp/Filezilla命令连接你的server能正常使用
鼓励实现加分项，具体加分数视实现难度及实现成熟程度而定

Usage

1.Open the terminal

3.run :

   sudo go run myftp.go
       
  
 4. the shell have some optional parameters
 
    SYNOPSIS
        sudo go run myftp.go [OPTION] ... 
        
    DESCRIPTION
    
      -p=PORT listening port 
      -h=HOST binding address 
      -d=DIR change current directory 
      -h print this help message
      
Note

    1. the ftp account is :
          name:"hanmeng"
          password:“abc"
    
    2. there are several default setting:
      port:                     21
      binding address:          127.0.0.1
      change current directory: /Users
      
    
    3.Supported commands
      USER
      PASS
      LIST
      CWD
      PWD
      STOR
      RETR
      QUIT
      PASV
       
    4. When you run a ftp-client,please run it at passive model：
      in os
        run:   
           ftp -p ftpserver-ip
      when you sign in,you can use command
        1. ls 
           View files under the current directory
        2. cd path
           change directory path
        3. pwd
           View  current directory path
        4.get filename
           download file from server
        5.put filename
           upload file to server
        6. bye
           disconnect from the server
           
    
