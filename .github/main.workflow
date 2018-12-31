workflow "Test & build" {
  on = "push"
  resolves = ["TestResult"]
}

workflow "Release a new version" {
  on = "release"
  resolves = ["ReleaseResult"]
}

action "Branch" {
  uses = "actions/bin/filter@master"
  args = "branch"
}

action "Codegen" {
  uses = "supinf/github-actions/go/codegen@master"
  args = "generate server -f spec.yaml -t app/generated"
}

action "Deps" {
  needs = ["Codegen"]
  uses = "supinf/github-actions/go/deps@master"
  env = {
    SRC_DIR = "app/"
  }
}

action "Lint" {
  needs = ["Branch", "Deps"]
  uses = "supinf/github-actions/go/lint@master"
  env = {
    SRC_DIR = "app/"
  }
}

action "Test" {
  needs = ["Deps"]
  uses = "supinf/github-actions/go/test@master"
  env = {
    SRC_DIR = "app/"
    IGNORE_DIR = "/generated/"
  }
}

action "Build" {
  needs = ["Deps"]
  uses = "supinf/github-actions/go/build@master"
  env = {
    SRC_DIR = "app/generated/cmd/return-everything-server/"
  }
}

action "TestResult" {
  needs = ["Lint", "Test", "Build"]
  uses = "actions/bin/debug@master"
}

action "Tags" {
  uses = "actions/bin/filter@master"
  args = "tag v*"
}

action "Release" {
  needs = ["Tags", "Build"]
  uses = "supinf/github-actions/github/release@master"
  env = {
    ARTIFACT_DIR = "app/generated/cmd/return-everything-server/dist/"
  }
  secrets = ["GITHUB_TOKEN"]
}

action "BuildImage" {
  needs = ["Tags", "Build"]
  uses = "supinf/github-actions/docker/build@master"
  args = "pottava/http-re:1.3"
  env = {
    DOCKERFILE = "prod/1.3/Dockerfile"
    BUILD_OPTIONS = "--no-cache"
  }
}

action "TagImage" {
  needs = ["BuildImage"]
  uses = "supinf/github-actions/docker/tag@master"
  env = {
    SRC_IMAGE = "pottava/http-re:1.3"
    DST_IMAGE = "pottava/http-re:latest"
  }
}

action "Login" {
  needs = ["BuildImage"]
  uses = "supinf/github-actions/docker/login@master"
  secrets = ["DOCKER_USERNAME", "DOCKER_PASSWORD"]
}

action "PushImage" {
  needs = ["TagImage", "Login"]
  uses = "supinf/github-actions/docker/push@master"
  args = "pottava/http-re:1.3,pottava/http-re:latest"
}

action "ReleaseResult" {
  needs = ["Release", "PushImage"]
  uses = "actions/bin/debug@master"
}
