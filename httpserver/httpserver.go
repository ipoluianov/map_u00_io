package httpserver

import (
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/ipoluianov/gomisc/logger"
	"github.com/ipoluianov/map_u00_io/utils"
)

var Instance *HttpServer

type HttpServer struct {
	srv    *http.Server
	srvTLS *http.Server
}

func NewHttpServer() *HttpServer {
	var c HttpServer
	return &c
}

func init() {
	Instance = NewHttpServer()
}

func (c *HttpServer) Start() {
	go c.thListen()
	go c.thListenTLS()
	go c.thTest()
}

func (c *HttpServer) thTest() {
	logger.Println("HttpServer thTest begin")
	privateKey, publicKey := utils.GenerateKeyPair()
	logger.Println("HttpServer thTest privateKey:", privateKey)
	logger.Println("HttpServer thTest publicKey:", publicKey)
	logger.Println("HttpServer thTest publicKey hex:", "0x"+hex.EncodeToString(publicKey))
	for {
		/*var item Item
		item.Address = "0x" + hex.EncodeToString(publicKey)
		item.DT = time.Now().Format("2006-01-02 15:04:05")
		item.Value = "test value " + item.DT
		item.Signature = utils.GenerateSignature(privateKey, []byte(item.Value))
		err := SetData(item)
		if err != nil {
			logger.Println("HttpServer thTest error:", err)
		}*/

		var item Item
		item.Address = "0x" + hex.EncodeToString(publicKey)
		item.DisplayName = "Test Data"
		item.DT = time.Now().Format("2006-01-02 15:04:05")
		item.Value = "test value " + item.DT
		item.Signature = utils.GenerateSignature(privateKey, []byte(item.Value))
		bs, _ := json.Marshal(item)
		hexData := hex.EncodeToString(bs)
		//http.Get("https://test.u00.io:8443/set-json-hex/" + hexData)
		http.Get("https://map.u00.io/set-json-hex/" + hexData)
		time.Sleep(1 * time.Second)
	}
}

func (c *HttpServer) portHttp() string {
	if utils.IsRoot() {
		return ":80"
	}
	return ":8080"
}

func (c *HttpServer) portHttps() string {
	if utils.IsRoot() {
		return ":443"
	}
	return ":8443"
}

func (c *HttpServer) thListen() {
	c.srv = &http.Server{
		Addr: c.portHttp(),
	}

	c.srv.Handler = c

	logger.Println("HttpServer thListen begin")
	err := c.srv.ListenAndServe()
	if err != nil {
		logger.Println("HttpServer thListen error: ", err)
	}
	logger.Println("HttpServer thListen end")
}

func (c *HttpServer) thListenTLS() {
	logger.Println("HttpServer::thListenTLS begin")
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = make([]tls.Certificate, 0)
	pathToBundle := logger.CurrentExePath() + "/bundle.crt"
	pathToPrivate := logger.CurrentExePath() + "/private.key"
	logger.Println("HttpServer::thListenTLS bundle.crt path:", pathToBundle)
	logger.Println("HttpServer::thListenTLS private.key path:", pathToPrivate)
	logger.Println("HttpServer::thListenTLS loading certificates ...")
	cert, err := tls.LoadX509KeyPair(pathToBundle, pathToPrivate)
	if err == nil {
		logger.Println("HttpServer::thListenTLS certificates is loaded SUCCESS")
		tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
	} else {
		logger.Println("HttpServer::thListenTLS loading certificates ERROR", err)
		return
	}

	serverAddress := c.portHttps()
	c.srvTLS = &http.Server{
		Addr:      serverAddress,
		TLSConfig: tlsConfig,
	}
	c.srvTLS.Handler = c

	logger.Println("HttpServer::thListenTLS starting server at", serverAddress)
	listener, err := tls.Listen("tcp", serverAddress, tlsConfig)
	if err != nil {
		logger.Println("HttpServer::thListenTLS starting server ERROR", err)
		return
	}

	logger.Println("HttpServer::thListenTLS starting server SUCCESS")
	err = c.srvTLS.Serve(listener)
	if err != nil {
		logger.Println("HttpServerTLS thListen error: ", err)
		return
	}
	logger.Println("HttpServer::thListenTLS end")
}

func (c *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		logger.Println("ProcessHTTP host: ", r.Host)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Request-Method", "GET")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			return
		}
		redirectUrl := ""
		if utils.IsRoot() {
			host := strings.ReplaceAll(r.Host, c.portHttp(), "")
			redirectUrl = "https://" + host + r.RequestURI

		} else {
			host := strings.ReplaceAll(r.Host, c.portHttp(), "")
			redirectUrl = "https://" + host + c.portHttps() + r.RequestURI
		}
		logger.Println("Redirect to HTTPS:", redirectUrl)
		http.Redirect(w, r, redirectUrl, http.StatusMovedPermanently)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Request-Method", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}

	parts := strings.FieldsFunc(r.RequestURI, func(r rune) bool {
		return r == '/'
	})

	var result []byte

	reqType := ""

	if len(parts) > 0 {
		reqType = parts[0]
	}

	if reqType == "get" {
		if len(parts) < 2 {
			w.WriteHeader(500)
			w.Write([]byte("wrong request: api - missing argument"))
			return
		}
		pageCode := parts[1]
		result := GetData(pageCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
		return
	}

	if reqType == "get-addresses" {
		var addresses []string
		for address := range storage.items {
			addresses = append(addresses, address)
		}
		slices.Sort(addresses)
		result, _ = json.Marshal(addresses)
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
		return
	}

	if reqType == "set-json-hex" {
		if len(parts) < 2 {
			w.WriteHeader(500)
			w.Write([]byte("wrong request: api - missing argument"))
			return
		}
		hexData := parts[1]
		data, err := hex.DecodeString(hexData)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("wrong request: api - hex decode error"))
			return
		}
		var item Item
		err = json.Unmarshal(data, &item)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("wrong request: api - json decode error"))
			return
		}

		SetData(item)
		return
	}

	// STATIC HTML
	html := string("123")
	html = strings.ReplaceAll(html, "%CONTENT%", "UNKNOWN QUERY")
	result = []byte(html)
	w.Write(result)
}
