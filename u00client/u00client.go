package u00client

import (
	"archive/zip"
	"bytes"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ipoluianov/gomisc/logger"
	"github.com/ipoluianov/map_u00_io/utils"
)

type U00Client struct {
	privateKey []byte
	publicKey  []byte
}

func NewClientWithKey(privateKey []byte) *U00Client {
	var c U00Client
	c.privateKey = privateKey
	if len(privateKey) >= 64 {
		c.publicKey = privateKey[32:64]
	}
	return &c
}

func NewClient() *U00Client {
	privateKey, _ := utils.GenerateKeyPair()
	return NewClientWithKey(privateKey)
}

func (c *U00Client) Address() string {
	if len(c.publicKey) != 32 {
		return ""
	}
	return "0x" + hex.EncodeToString(c.publicKey)
}

func (c *U00Client) sendPostBytes(url string, data []byte, contentType string) ([]byte, int, error) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Post(url, contentType, bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return body, resp.StatusCode, nil
}

func (c *U00Client) WriteValue(name string, dt time.Time, value string) error {
	if len(c.privateKey) != 64 || len(c.publicKey) != 32 {
		return errors.New("private key is not set or public key is empty")
	}

	domain := hex.EncodeToString(c.publicKey[:1])
	domain = domain[:1]

	var err error
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	var zipFile io.Writer
	{
		zipFile, err = zipWriter.Create("value")
		if err == nil {
			zipFile.Write([]byte(value))
		}
	}
	{
		zipFile, err = zipWriter.Create("name")
		if err == nil {
			zipFile.Write([]byte(name))
		}
	}
	{
		value = dt.Format("2006-01-02 15:04:05.000")
		zipFile, err = zipWriter.Create("time")
		if err == nil {
			zipFile.Write([]byte(value))
		}
	}
	zipWriter.Close()
	zipFileContent := buf.Bytes()

	signature := ed25519.Sign(c.privateKey, zipFileContent)

	frame := make([]byte, 32+len(zipFileContent)+64)
	copy(frame[:32], c.publicKey)
	copy(frame[32:32+64], signature)
	copy(frame[32+64:], zipFileContent)

	respBS, status, err := c.sendPostBytes("https://s"+domain+".u00.io/set", frame, "application/octet-stream")
	if err != nil {
		logger.Println("U00Client WriteValue error:", err, respBS, status)
		return err
	}

	logger.Println("U00Client WriteValue response:", string(respBS), status)

	return nil
}

func (c *U00Client) ReadValue(address string) (value []byte, err error) {
	if len(c.publicKey) != 32 {
		return nil, errors.New("public key is not set or invalid")
	}

	domain := hex.EncodeToString(c.publicKey[:1])
	domain = domain[:1]

	bs, err := http.Get("https://s" + domain + ".u00.io/get/" + address)
	if err != nil {
		logger.Println("U00Client ReadValue error:", err)
		return nil, err
	}

	data, err := io.ReadAll(bs.Body)
	bs.Body.Close()
	if err != nil {
		logger.Println("U00Client ReadValue error reading body:", err)
		return nil, err
	}

	fmt.Println("GET response:", hex.EncodeToString(data))
	return nil, nil
}
