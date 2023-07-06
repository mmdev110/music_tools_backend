#開発用
FROM golang:1.19.2

#install gopls for auto-formatting
#RUN go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest
#RUN go install github.com/ramya-rao-a/go-outline@latest
#RUN go install github.com/cweill/gotests/gotests@latest
#RUN go install github.com/fatih/gomodifytags@latest
#RUN go install github.com/josharian/impl@latest
#RUN go install github.com/haya14busa/goplay/cmd/goplay@latest
#RUN go install honnef.co/go/tools/cmd/staticcheck@latest
RUN go install golang.org/x/tools/gopls@latest
#install air for hot-reloading
RUN go install github.com/cosmtrek/air@latest

WORKDIR /go/src/
ADD . /go/src/
RUN go mod tidy
#CMD ["go","run","app/main.go"]
CMD ["air","-c",".air.toml"]