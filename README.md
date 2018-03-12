# htags
因為隔壁棚的 [App](https://github.com/DaidoujiChen/Dai-Hentai) 在顯示 tag 的時候一直都是英文的, 覺得很不開心, 
所以從別人家借了一個 [對照表](https://github.com/Mapaler/EhTagTranslator) 來用, 為了方便, 所以寫成一個小小的 service 放在 app engine 上,
根據我自己的需求, 我只需要中英文對照而已, 所以就把他整理成一個 json 吐回來.

## Service Link
[點我取得列表](https://alltesthere-186305.appspot.com/)

## Deploy
這邊是用簡單的 [Buddy](https://app.buddy.works) 來做自動化, 雖然這中間也遇到很多有的沒有的問題就是, 這邊就一次整理起來, 方便未來查找

![](https://s3-ap-northeast-1.amazonaws.com/daidoujiminecraft/Daidouji/%E8%9E%A2%E5%B9%95%E5%BF%AB%E7%85%A7+2018-03-12+%E4%B8%8B%E5%8D%8811.06.34.png)

Action 的部分只有一個步驟

![](https://s3-ap-northeast-1.amazonaws.com/daidoujiminecraft/Daidouji/%E8%9E%A2%E5%B9%95%E5%BF%AB%E7%85%A7+2018-03-12+%E4%B8%8B%E5%8D%8811.09.55.png)

就是建立一個 Golang 的 Docker =w=, 版本選擇 `1.8`, 為了要用 App Engine `Standard`, 而不是用 `Flexible` 的呦, 後面會再說到他的特性, 總之,
在 Run Command 的部分設定

![](https://s3-ap-northeast-1.amazonaws.com/daidoujiminecraft/Daidouji/%E8%9E%A2%E5%B9%95%E5%BF%AB%E7%85%A7+2018-03-12+%E4%B8%8B%E5%8D%8811.15.02.png)

```
export GOPATH=/
export GO15VENDOREXPERIMENT=1
go get -u google.golang.org/appengine
go test
go build

#deploy
# 這邊這個 serviceAccount.json 需要從 gcp 的主控台那邊生長出來
gcloud auth activate-service-account --key-file ./config/serviceAccount.json
export CLOUDSDK_CORE_PROJECT="alltesthere-186305"

#移除掉多餘的版本, 一開始就使用的話會有錯誤, 需要有兩個版本以上時才可以正常運作
gcloud app versions delete --service default `gcloud app versions list | awk '{if ($3 == "0.00") {print $2}}'` --quiet
gcloud app deploy app.yaml --quiet --project alltesthere-186305
```

另外, 在 Environment Customization 的部分

![](https://s3-ap-northeast-1.amazonaws.com/daidoujiminecraft/Daidouji/%E8%9E%A2%E5%B9%95%E5%BF%AB%E7%85%A7+2018-03-12+%E4%B8%8B%E5%8D%8811.20.06.png)

設定預先需要安裝在這個 docker image 裡面的檔案們

```
#install glcoud
apt-get update
apt-get install -y apt-transport-https
echo "deb https://packages.cloud.google.com/apt cloud-sdk-jessie main" | tee -a /etc/apt/sources.list.d/google-cloud.list
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
apt-get update
apt-get install -y google-cloud-sdk-app-engine-go
```
Buddy 的部分就大概這樣了

## App Engine
前面說到要用 Standard 而不用 Flexible 的原因在於, 如果我們有把 scale 設定好的話, 只要沒有 request 打進來, instance 的數量可以縮減到 0 台, 
就不用負擔長期開機的成本 O3Ob

![](https://s3-ap-northeast-1.amazonaws.com/daidoujiminecraft/Daidouji/%E8%9E%A2%E5%B9%95%E5%BF%AB%E7%85%A7+2018-03-12+%E4%B8%8B%E5%8D%8811.24.46.png)

一開始 instance 很多是因為沒有設好, 後面就可以掉到 0 台囉
