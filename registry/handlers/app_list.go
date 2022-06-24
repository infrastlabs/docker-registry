package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"crypto/tls"
	"flag"
	"encoding/json"
	"html/template"
	"sort"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

var (
	//addr = flag.String("addr", ":8080", "ui address")
	apis = flag.String("api", "http://localhost:50001", "api address")
	user = flag.String("user", "admin", "registry user")
	pass = flag.String("pass", "xxx", "registry pass")
	size = flag.Bool("size", false, "registry size") //or: req's param
	//tout = flag.Duration("tout", time.Second, "api cache timeout")
	host string // host:port of api server
)
func parseArgs() {
	flag.Parse()
	api, err := url.Parse(*apis)
	if err != nil {
		panic(err)
	}
	host = api.Host
}

func getClient() (*http.Client) {
	//client := &http.Client{}

	//TODO: 判断registry的参数是否为https模式；
	//https://blog.csdn.net/weixin_43064185/article/details/125300410

	// 忽略 https 证书验证
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	return client
}
func (app *App) imageList(w http.ResponseWriter, r *http.Request) {
// func imageList() error {
	// config *configuration.Configuration
	// config := registry.config
	// fmt.Println("===app.Config.List.Apis:  "+app.Config.List.Apis)
	// *apis= app.Config.List.Apis //TODO tls
	*apis= "http://localhost:8143"
	*user= app.Config.List.User
	*pass= app.Config.List.Pass
	*size= app.Config.List.Size
	parseArgs()


	/* res, err := http.Get(*apis + "/v2/_catalog")	
	if err != nil {
		return err
	} */

	
	//生成client 参数为默认
	client := getClient() //&http.Client{}
	url := *apis + "/v2/_catalog"
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(*user, *pass)//设置需要认证的username和password
	if err != nil {
		panic(err)
	}

	//处理返回结果
	res, err := client.Do(req)
	if err != nil {
		// panic(err)
		w.Write([]byte("请求错误，请检查地址："+*apis + "/v2/_catalog"))
		return
	}
	// TODO: check status
	var catalog struct {
		Repos []string `json:"repositories"`
	}
	if err := json.NewDecoder(res.Body).Decode(&catalog); err != nil {
		return 
	}
	if err := res.Body.Close(); err != nil {
		return 
	}

	data := make(map[string][]string, len(catalog.Repos))
	

	for _, name := range catalog.Repos {
		// res, err := http.Get(*apis + "/v2/" + name + "/tags/list")
		// if err != nil {
		// 	return err
		// }

		//生成client 参数为默认
		client := &http.Client{}
		url := *apis + "/v2/" + name + "/tags/list"
		req, err := http.NewRequest("GET", url, nil)
		req.SetBasicAuth(*user,*pass)//设置需要认证的username和password
		if err != nil {
			panic(err)
		}
	
		//处理返回结果
		res, _ := client.Do(req)
		// defer res.Body.Close()
		// body, err := ioutil.ReadAll(res.Body)
		// fmt.Printf("%s\n",body)		


		// TODO: check status
		var tags struct {
			Name string   `json:"name"`
			Tags []string `json:"tags"`
		}
		if err := json.NewDecoder(res.Body).Decode(&tags); err != nil {
			return 
		}
		data[name] = tags.Tags
		if err := res.Body.Close(); err != nil {
			return 
		}
	}

	auth := remote.WithAuth(authn.FromConfig(authn.AuthConfig{
		Username: *user,
		Password: *pass,
	}))
	data2 := make([]string, 0, len(catalog.Repos))
	for repo, tags := range data {
		for _, tag := range tags {
			// data2 = append(data2, host+"/"+repo+":"+tag)
			//data2 = append(data2, repo+":"+tag)
			if true!=*size {
				data2 = append(data2, repo+":"+tag)
			} else {
				tagName, _ := name.NewTag(host+"/"+repo+":"+tag) //host > *apis
				img, err := remote.Image(tagName, auth)
				if err != nil {
					println("causing failure " + err.Error())
					//return nil
					data2 = append(data2, repo+":"+tag+" (errSize)")
				} else {

					//digest, _ := img.Digest()
					layers, _ := img.Layers()

					var imgSize int64
					for _, layer := range layers {
						size, _ := layer.Size()
						imgSize += size
					}

					data2 = append(data2, repo+":"+tag+" ("+ByteCountSI(imgSize)+")")
				}
			}
		}
	}
	sort.Strings(data2)

	if err := htmlTPL.Execute(w, data2); err != nil {
		println("causing failure " + err.Error())
	}
	return 
}
var htmlTPL = template.Must(template.New("").Parse(`<!DOCTYPE html>
<html lang="en" dir="ltr">
	<head>
		<meta charset="utf-8">
		<title>Registry Listing</title>
	</head>
	<body>
		<ul>
		{{- range .}}
			<li>{{.}}</li>
		{{- else}}<li>No Repositories Found</li>{{end}}
		</ul>
	</body>
</html>`))


func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
