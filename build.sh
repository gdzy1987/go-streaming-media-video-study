#! /bin/bash
case $1 in
    mac )
        # 构建web和其他services
        echo "start build mac ..."

        cd $GOPATH/src/go-streaming-media-video-study/api
        go build -o ../bin/api

        cd $GOPATH/src/go-streaming-media-video-study/scheduler
        go build -o ../bin/scheduler

        cd $GOPATH/src/go-streaming-media-video-study/streamserver
        go build -o ../bin/streamserver

        cd $GOPATH/src/go-streaming-media-video-study/web
        go build -o ../bin/web

        cp $GOPATH/src/go-streaming-media-video-study/config/conf.json     $GOPATH/src/go-streaming-media-video-study/bin/
    ;;

    linux  )
        # 构建web和其他services
         echo "start build linux ..."

        cd $GOPATH/src/go-streaming-media-video-study/api
        env GOOS=linux GOARCH=amd64 go build -o ../bin/api

        cd $GOPATH/src/go-streaming-media-video-study/scheduler
        env GOOS=linux GOARCH=amd64 go build -o ../bin/scheduler

        cd $GOPATH/src/go-streaming-media-video-study/streamserver
        env GOOS=linux GOARCH=amd64 go build -o ../bin/streamserver

        cd $GOPATH/src/go-streaming-media-video-study/web
        env GOOS=linux GOARCH=amd64 go build -o ../bin/web

        cp $GOPATH/src/go-streaming-media-video-study/config/conf.json     $GOPATH/src/go-streaming-media-video-study/bin/
    ;;
    *)

    echo "usage: build[linux|mac]"
esac