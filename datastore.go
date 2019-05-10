package ipcarta

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
)

type ElasticResponse struct {
	Found  bool            `json:"found"`
	Source json.RawMessage `json:"_source"`
}

func makeKey(n string) string {
	h := md5.New()
	io.WriteString(h, n)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Set(n Network) {
	payload, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	key := makeKey(n.Network)
	endpoint := fmt.Sprintf("%s/ipcarta/ips/%s", config.ElasticSearchHost, key)
	fmt.Println(endpoint)
	fmt.Println(string(payload))
	resp, err := http.Post(endpoint, "text/json", bytes.NewBuffer(payload))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
}

func Search(ip string) (bool, json.RawMessage) {
	for i := 32; i >= 0; i-- {
		_, ipv4Net, err := net.ParseCIDR(fmt.Sprintf("%s/%s", ip, strconv.Itoa(i)))
		if err != nil {
			return false, nil
		}
		res, err := http.Get(fmt.Sprintf("%s/ipcarta/ips/%s", config.ElasticSearchHost, makeKey(ipv4Net.String())))
		if err != nil {
			return false, nil
		}
		defer res.Body.Close()
		msg, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		var response ElasticResponse
		json.Unmarshal(msg, &response)
		if response.Found {
			return true, response.Source
		}
	}
	return false, nil
}
