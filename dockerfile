FROM <docker_url>/core/golang-buildtools:1.13.8 as GoBuildStage
WORKDIR /srv/package

USER root
COPY . ./src/github.com/diptamay/go-commons
RUN chown -R go ./src/github.com/diptamay/go-commons

USER go

RUN \
    export GOFLAGS=-mod=vendor \
    project_path=/srv/package/src/github.com/diptamay/go-commons \
    && cd $project_path \
    && go build

RUN go test ./... -short -v -cover -failfast -timeout 300s

LABEL "Description"="Go commons"