package registry

import (
	"fmt"
	"net/http"
	"net/url"

	"crypto/tls"
	"flag"
	"encoding/json"
	"html/template"
	"sort"
	"strings"
	//"strconv"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v1 "github.com/google/go-containerregistry/pkg/v1"

	// "github.com/distribution/distribution/v3/configuration"
	configuration "cn.dev.docker-registry/conf"
)

var (
	//addr = flag.String("addr", ":8080", "ui address")
	// apis = flag.String("api", "https://localhost:8143", "api address")
	user = flag.String("user", "admin", "registry user")
	pass = flag.String("pass", "admin123", "registry pass")
	size = flag.Bool("size", true, "registry size") //or: req's param
	//tout = flag.Duration("tout", time.Second, "api cache timeout")
	apis, host string // host:port of api server
)
func parseArgs(config *configuration.Configuration) {
	flag.Parse()
	//apis= config.List.Apis 
	// apis= "http://localhost:8143"

	// TODO extend Config类?
	/* *user= "admin" //config.List.User
	*pass= "admin123" //config.List.Pass
	*size= true //config.List.Size */
	*user= config.List.User
	*pass= config.List.Pass
	*size= config.List.Size

	//set apis
	schema:= "http://"
	if ""!=config.HTTP.TLS.Certificate {
		schema= "https://"
	}
	addr:= config.HTTP.Addr
	strs:= strings.Split(addr,":")
	if ""==strs[0] {
		strs[0]="127.0.0.1"
	}
	apis=schema+strs[0]+":"+strs[1] 
	//fmt.Println("apis="+apis) //apis=http://127.0.0.1:8143
	
	api, err := url.Parse(apis)
	// api, err := url.Parse(config.HTTP.Addr) //url_err
	if err != nil {
		panic(err)
	}
	host = api.Host
}

func getClient(config *configuration.Configuration) (*http.Client) {
	client := &http.Client{}

	//TODO: 判断registry的参数是否为https模式；
	//https://blog.csdn.net/weixin_43064185/article/details/125300410
	if ""!=config.HTTP.TLS.Certificate {
		// 忽略 https 证书验证
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: transport}
	}
	return client
}

func doGet(client *http.Client, uri string, config *configuration.Configuration)(*http.Response, error){
	url := apis + uri
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(*user, *pass)//设置需要认证的username和password
	if err != nil {
		panic(err)
	}

	//处理返回结果
	return client.Do(req)
}
func imageList(w http.ResponseWriter, r *http.Request, config *configuration.Configuration) {
	// fmt.Println("===config.List.Apis:  "+config.List.Apis)
	parseArgs(config)
	// fmt.Println("==============host:  "+host+config.HTTP.Addr)

	//生成client 参数为默认
	client := getClient(config) //&http.Client{}
	
	//catalog获取仓库列表
	uri:="/v2/_catalog"
	rsp, err := doGet(client, uri, config)
	if err != nil {
		// panic(err)
		w.Write([]byte("RequestError, uri:"+uri))
		return
	}
	//fmt.Println("rsp.StatusCode "+rsp.StatusCode)
	if 200!=rsp.StatusCode {
		w.Write([]byte("RequestError, Status:"+rsp.Status)) //strconv.Itoa(rsp.StatusCode)
		return
	}
	// TODO: check status
	var catalog struct {
		Repos []string `json:"repositories"`
	}
	if err := json.NewDecoder(rsp.Body).Decode(&catalog); err != nil {
		return 
	}
	if err := rsp.Body.Close(); err != nil {
		return 
	}

	//tags/list获取各img的标签列表
	imgMap := make(map[string][]string, len(catalog.Repos))
	for _, name := range catalog.Repos {
		//生成client 参数为默认
		//client := &http.Client{}
		uri= "/v2/" + name + "/tags/list"
		rsp, err := doGet(client, uri, config)
		if err != nil {
			w.Write([]byte("RequestError, uri:"+uri))
			return
		}

		/* defer rsp.Body.Close()
		body, err := ioutil.ReadAll(rsp.Body)
		fmt.Printf("%s\n",body)	 */	

		// TODO: check status
		var tags struct {
			Name string   `json:"name"`
			Tags []string `json:"tags"`
		}
		if err := json.NewDecoder(rsp.Body).Decode(&tags); err != nil {
			return 
		}
		imgMap[name] = tags.Tags
		if err := rsp.Body.Close(); err != nil {
			return 
		}
	}

	//auth remote.Option
	auth := remote.WithAuth(authn.FromConfig(authn.AuthConfig{
		Username: *user,
		Password: *pass,
	}))
	// https, skip_key_validate
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	trans := remote.WithTransport(transport)
	//遍历各img所有tags的layer层信息（获取size）
	detailTags := make([]string, 0, len(catalog.Repos))
	for oneImg, tags := range imgMap {
		for _, tag := range tags {
			if true!=*size {
				detailTags = append(detailTags, oneImg+":"+tag)
			} else {
				tagName, _ := name.NewTag(host+"/"+oneImg+":"+tag) //host > apis
				//remote获取layer层信息: http or https
				var img v1.Image
				if ""!=config.HTTP.TLS.Certificate {
					img, err = remote.Image(tagName, auth, trans) //https
				} else {
					img, err = remote.Image(tagName, auth)
				}
				if err != nil {
					println("causing failure " + err.Error())
					//return nil
					detailTags = append(detailTags, oneImg+":"+tag+" (errSize)")
				} else {
					//digest, _ := img.Digest()
					layers, _ := img.Layers()

					var imgSize int64
					for _, layer := range layers {
						size, _ := layer.Size()
						imgSize += size
					}
					detailTags = append(detailTags, oneImg+":"+tag+" ("+ByteCountSI(imgSize)+")")
				}
			}
		}
	}
	sort.Strings(detailTags)

	if err := htmlTPL.Execute(w, detailTags); err != nil {
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
		API:
		<ul>
			<li><a target="_blank" href="https://docs.docker.com/registry/spec/api/">/v2/api</a></li>
			<li><a target="_blank" href="/v2/_catalog">/v2/_catalog</a></li>
		</ul>
		PUBLIC:
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
