package main

import (
	"app"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	KEY_FILE_PREFIX = "keystore"
)

var addr = flag.String("addr","localhost:4000","http service address")
var number = flag.Int("number",20,"mobile account number")

func main() {
	f,err:= os.OpenFile(filepath.Join("etc","logfile"),os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
	if err!=nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	var wg sync.WaitGroup
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	u:= url.URL{Scheme:"wss",Host:*addr,Path:"/"}
	httpPath:= "https://"+*addr;
	log.Println("connecting to %s", u.String())

	for i:=0;i<*number;i++ {
		wg.Add(1)
		go func(id int) {
			defer func() {
				wg.Done()
			}()
			// load wallet
			var account *app.Wallet
			keystoreFile:= filepath.Join("etc",KEY_FILE_PREFIX+strconv.Itoa(id))
			if _,err:= os.Stat(keystoreFile); os.IsNotExist(err) {
				account:= app.NewWallet()
				keystore:= &app.KeyStore{
					Address:account.Address,
					PrivateKey: account.PrivateKey.D.Text(16),
				}
				data,err:=json.Marshal(keystore)
				if err!=nil {
					log.Println(err.Error())
					panic(err)
				}
				err =ioutil.WriteFile(keystoreFile,data,0666)
				if err!=nil {
					log.Println(err.Error())
					panic(err)
				}
			} else {
				account,err = app.NewWalletFromFile(keystoreFile)
				if err!=nil {
					log.Fatal("user",id,"can not create wallet")
				}
			}


			client:=app.NewClient(id,u.String(),httpPath,account)
			client.Start()
		}(i)
	}
	wg.Wait()
}
