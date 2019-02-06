workflow "Build Container" {
  on = "push"
  resolves = ["Push container"]
}

action "Prepare modules" {
  uses = "retgits/actions/gocenter@master"
  args = ["mod vendor"]
}

action "Update project" {
  uses = "retgits/actions/sh@master"
  needs = "Prepare modules"
  args = ["mkdir function", "rm *_test.go", "mv *.go function/", "mv go.* function/", "mv vendor function/", "git clone https://github.com/openfaas-incubator/golang-http-template templates", "rm -rf templates/template/golang-http/function ", "cp -R templates/template/golang-http/ .", "rm -rf templates"]
}

action "Docker Login" {
  uses = "actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "Build container" {
  needs = ["Update project", "Docker Login"]
  uses = "actions/docker/cli@master"
  args = "build . -t retgits/openfaas-trello:$GITHUB_SHA"
}

action "Push container" {
  needs = "Build container"
  uses = "actions/docker/cli@master"
  args = "push retgits/openfaas-trello:$GITHUB_SHA"
}
