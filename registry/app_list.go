package registry

import (
	"fmt"
	"net/http"
	"net/url"
	// "errors"
	// "strconv"

	"crypto/tls"
	"encoding/json"
	// "flag"
	"html/template"
	"sort"
	"strings"
	"time"

	//"strconv"
	"github.com/google/go-containerregistry/pkg/authn"
	// v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"

	// "github.com/distribution/distribution/v3/configuration"
	configuration "gitee.com/g-devops/docker-registry/conf"
	cmap "github.com/orcaman/concurrent-map"
)

var (
	// //addr = flag.String("addr", ":8080", "ui address")
	// // apis = flag.String("api", "https://localhost:8143", "api address")
	// user = flag.String("user", "admin", "registry user")
	// pass = flag.String("pass", "admin123", "registry pass")
	// size = flag.Bool("size", true, "registry size") //or: req's param
	// //tout = flag.Duration("tout", time.Second, "api cache timeout")
	
	user, pass string
	size bool
	apis, host string // host:port of api server
	conf2 *configuration.Configuration
	tagsizeMap  cmap.ConcurrentMap
	countMap  cmap.ConcurrentMap
)

// 后台同步imgListSize
//  1.刷入所有列表img_ver；
//  2.判断one.size是否就绪，无则读取;
func doRefreshTagsize() {
	for {
		_, err:= getRepoTags(conf2) //repoTags
		if nil!=err {
			// w.Write([]byte("getRepoTags err:"+err.Error()))
			// return
			time.Sleep(2*time.Second) //avoid cpu 100%
			continue
		}

		countImgSize(conf2)

		time.Sleep(10*time.Second) //60>3>20 //TODO loop间隔可配
	}

	// tagsizeMap.Set(key, tunnel)
	/* if item, ok := service.tagsizeMap.Get(key); ok {
		tunnelDetails := item.(*portainer.TunnelDetails)
		return tunnelDetails
	} */

	/* for item := range service.tagsizeMap.IterBuffered() {
		tunnel := item.Val.(*portainer.TunnelDetails)
		if tunnel.Port == port {
			return service.getUnusedPort()
		}
	} */

	// return
}

func parseArgs(config *configuration.Configuration) {
	// flag.Parse()
	//apis= config.List.Apis 
	// apis= "http://localhost:8143"

	// TODO extend Config类?
	/* *user= "admin" //config.List.User
	*pass= "admin123" //config.List.Pass
	*size= true //config.List.Size */
	user= config.Extend.List.User
	pass= config.Extend.List.Pass
	size= config.Extend.List.Size

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

func doGet(client *http.Client, uri string)(*http.Response, error){
	url := apis + uri
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(user, pass)//设置需要认证的username和password
	if err != nil {
		panic(err)
	}

	//处理返回结果
	return client.Do(req)
}

var client *http.Client
func getRepoTags(config *configuration.Configuration) ([]string, error) {
	// fmt.Println("===config.List.Apis:  "+config.List.Apis)
	parseArgs(config)
	// fmt.Println("==============host:  "+host+config.HTTP.Addr)

	//生成client 参数为默认
	if nil==client {
		client = getClient(config) //&http.Client{}
	}
	
	//catalog获取仓库列表
	uri:="/v2/_catalog"
	rsp, err := doGet(client, uri)
	
	// 保证关闭连接. 不关闭连接将导致close-wait累积，最终占满端口。监控将报错:cannot assign requested address
	//   在return前，定义defer做rsp的关闭,免mem未释放 占满内存
	if rsp != nil { // 当请求失败，resp为nil时，直接defer会导致panic，因此需要先判断
		defer rsp.Body.Close()
	}
	/* if err := rsp.Body.Close(); err != nil {
		return nil, err
	} */

	if err != nil {
		// panic(err)
		// w.Write([]byte("RequestError, uri:"+uri))
		return nil, errors.New("RequestError, uri:"+uri) //+err.Error()
	}
	//fmt.Println("rsp.StatusCode "+rsp.StatusCode)
	if 200!=rsp.StatusCode {
		// w.Write([]byte("RequestError, Status:"+rsp.Status)) //strconv.Itoa(rsp.StatusCode)
		return nil, errors.New("RequestError, Status:"+rsp.Status)
	}
	// TODO: check status
	var catalog struct {
		Repos []string `json:"repositories"`
	}
	if err := json.NewDecoder(rsp.Body).Decode(&catalog); err != nil {
		return nil, err
	}
	

	//tags/list获取各img的标签列表
	imgMap := make(map[string] []string, len(catalog.Repos)) //[]string >> map[string] []string
	for _, name := range catalog.Repos {
		//生成client 参数为默认
		//client := &http.Client{}
		uri= "/v2/" + name + "/tags/list"
		rsp, err := doGet(client, uri)
		if err != nil {
			// w.Write([]byte("RequestError, uri:"+uri))
			return nil, errors.New("RequestError, uri:"+uri)
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
			return nil, err
		}
		imgMap[name] = tags.Tags
		if rsp != nil {
			defer rsp.Body.Close()
		}
		/* if err := rsp.Body.Close(); err != nil {
			return nil, err
		} */
	}

	// detailTags := make([]string, 0, len(catalog.Repos))
	repoTags := make([]string, 0) //[]string
	for oneImg, tags := range imgMap {
		for _, tag := range tags {
			repoTags = append(repoTags, oneImg+":"+tag)
		}
	}
	sort.Strings(repoTags)

	/* if nil==tagsizeMap {
		tagsizeMap= cmap.New()
	} */
	/* return &Service{
		tagsizeMap: cmap.New(), //ref pt's
		dataStore:        dataStore,
		shutdownCtx:      shutdownCtx,
	} */

	for _, imgtag := range repoTags {
		if !tagsizeMap.Has(imgtag) { //不存在时才put, 免val覆盖
			tagsizeMap.Set(imgtag, "") //nil> init ""
		}

		// TODO; if img deleted > loop gMap drop non-exist
	}

	// client.CloseIdleConnections() //each free> use global cli

	return repoTags, nil
}

var transport *http.Transport
func goContainerregistryImageSize(imageTag, tlsCert string) (string, error) {
	//auth remote.Option
	auth := remote.WithAuth(authn.FromConfig(authn.AuthConfig{
		Username: user,
		Password: pass,
	}))
	// https, skip_key_validate
	if nil==transport { //
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	trans := remote.WithTransport(transport)


	/* var err error
	var imgSize int64
	imgSize= 0

	tagName, _ := name.NewTag(imageTag) //host > apis

	// ref
	// https://github.com/docker/index-cli-plugin/blob/56dd3635cd47826d82f4ee90a2b73883c2549735/registry/save.go#L304
	ref, err := name.ParseReference(imageTag)
	if err != nil {
		return "", nil //errors.Wrapf(err, "failed to parse reference: %s", imageTag)
	}

	//remote获取layer层信息: http or https
	var img v1.Image
	if ""!=tlsCert {
		img, err = remote.Image(ref, auth, trans) //https
	} else {
		img, err = remote.Image(tagName, auth)
	}
	if err != nil {
		println("causing failure " + err.Error())
		//return nil
		// detailTags = append(detailTags, imgtag+" (errSize)")
	} else {
		//digest, _ := img.Digest()
		layers, _ := img.Layers() //TODO: 按multiArch,计算各arch的layers汇总Size

		// //multi: 无layers? || layers汇总？
		for _, layer := range layers {
			// layer.MediaType() //string
			size, _ := layer.Size()
			imgSize += size
		}
		// detailTags = append(detailTags, imgtag+" ("+ByteCountSI(imgSize)+")")
	} */

	imageRefs, err := GetImageReferences(imageTag)
	if err != nil {
		fmt.Printf("\nerror in getImageReferences: %s", err.Error())
		return "", errors.Wrap(err, "getting image references from container registry")
	}
	if len(imageRefs) == 0 {
		return "", fmt.Errorf("\n%d image references found for %q", len(imageRefs), imageTag)
	}
	// refData := imageRefs[0]
	var multiSize string
	// DO archs sort排序
	sizeArr:= []string{}
	for _, imageRef:= range imageRefs {
		ref, err:= name.ParseReference(imageRef.Digest) //refData.Digest
		if err != nil {
			return "", err
		}
		
		img, err:= remote.Image(ref, auth, trans)
		if err != nil {
			return "", err
		}
		var imgSize int64
		layers, err := img.Layers()
		if err != nil {
			return "", err
		}
		for _, layer := range layers {
			size, _:= layer.Size()
			imgSize += size
		}
		// // multiSize+=fmt.Sprintf(imageRef.OS+"/"+imageRef.Arch+" size: "+ByteCountSI(imgSize))
		// // multiSize+=fmt.Sprintf("(%s/%s %s)", imageRef.OS, imageRef.Arch, ByteCountSI(imgSize))
		// multiSize+=fmt.Sprintf(" %s:%s |", imageRef.Arch, ByteCountSI(imgSize))
		sizeArr= append(sizeArr, fmt.Sprintf(" %s:%s", imageRef.Arch, ByteCountSI(imgSize)))
	}
	// 对字符串数组进行排序
	sort.Strings(sizeArr)
	// 打印排序后的数组
	fmt.Println(sizeArr)

	// multiSize= multiSize[1:len(multiSize)-2] //drop last "| "
	// multiSize= fmt.Sprintf(" %s", sizeArr)
	multiSize= strings.Join(sizeArr, " | ")
	return multiSize, nil
}

func countImgSize(config *configuration.Configuration) (error) {
	//遍历各img所有tags的layer层信息（获取size）
	// detailTags := make([]string, 0) //[]string
	// for _, imgtag := range repoTags {
	for item := range tagsizeMap.IterBuffered() {
		imgtag:= item.Key
		size0:= item.Val.(string)
		
		// 判断遍历计数:免每次都计算size;
		if ""!=size0 { //size设置过，才做重复判断; 否则直接重算
			cnt, exist:= countMap.Get(imgtag)
			if !exist {
				countMap.Set(imgtag, 1)
				// continue //首次需要计算
			} else {
				// if ""!=val {
				cnt2:= cnt.(int)
				if cnt2<=5 { //10s * 5; 不超过5次 则continue跳过计算;(免imgtag重推)
					// DO: ttl alive 5min?>> counts
					countMap.Set(imgtag, cnt2 +1)
					continue //跳过计算size
				} else {
					// fmt.Printf("==reCountSize: %s\n", imgtag)
					countMap.Set(imgtag, 1) //重置计数; 并进入后续:重计算imgtag尺寸
				}
			}
		}

		// msg:= fmt.Sprintf("==init cntSize: %s, val: %s", imgtag, val)
		// println(msg)

		// detailTags = append(detailTags, imgtag)
		// tunnel := item.Val.(*portainer.TunnelDetails)
		if true!=size {
			// detailTags = append(detailTags, imgtag)
		} else {
			var imgSize string
			// imgSize= 0

			// imgSize, err:= undockInspectImage(host, imgtag)
			imgSize, err:= goContainerregistryImageSize(host+"/"+imgtag, config.HTTP.TLS.Certificate)
			if nil!=err {
				fmt.Println("err: "+err.Error())
			}
			tagsizeMap.Set(imgtag, imgSize) //just replace the same key.
		}
	}

	return nil
}

func imageList(w http.ResponseWriter, r *http.Request) {
	// 展示列表前，把tags清单刷一遍(之后由loop取size)
	_, err:= getRepoTags(conf2) //repoTags
	if nil!=err {
		w.Write([]byte("getRepoTags err:"+err.Error()))
		return
	}

	/* // config:= conf2
	countImgSize(conf2) */
	
	//遍历各img所有tags的layer层信息（获取size）
	detailTags := make([]string, 0) //[]string
	// for _, imgtag := range repoTags {
	for item := range tagsizeMap.IterBuffered() {
		imgtag:= item.Key
		imgsize:= item.Val.(string)
		if ""==imgsize {
			imgsize= "errSize"
		}
		// format: img:tag # arch:imgsize | arch2:imgsize2
		one:= fmt.Sprintf("%s # %s", imgtag, imgsize)
		if true!=size {
			one= imgtag
		}
		detailTags = append(detailTags, one)
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
	const unit = 1024 //1000
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
