# push
app push client 

use golang develop

> apns: apple push client   
> mipush: xiaomi push client     
> fcm： google fcm push client    


必备命令行工具 OpenSSL

1.打开终端 输入以下命令,导出cert文件,在命令执行过程中需要输入证书密码 

```
$ openssl pkcs12 -clcerts -nokeys -out iphone-dev-cert.pem -in iphone-dev.p12
```

2.导出加密过的key文件,在命令执行过程中需要输入证书密码 

```
$ openssl pkcs12 -nocerts -out iphone-dev-key-pwd.pem -in iphone-dev.p12
```

3.去掉key文件中的密码,在命令执行过程中需要输入证书密码 

```
$ openssl rsa -in iphone-dev-key-pwd.pem -out iphone-dev-key-noenc.pem
$ mv iphone-dev-key-noenc.pem iphone-dev-key.pem
```

4.验证pem有效性

测试环境   

```
$ openssl s_client -connect gateway.sandbox.push.apple.com:2195 -cert iphone-dev-cert.pem -key iphone-dev-key.pem
```
生产环境

```
$ openssl s_client -connect gateway.push.apple.com:2195 -cert iphone-cert.pem -key iphone-key.pem
```
