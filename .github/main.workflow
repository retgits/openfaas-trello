workflow "Validate Actions" {
  on = "push"
  resolves = ["Hello"]
}

action "Hello" {
  uses = "retgits/actions/sh@master"
  args = ["ls -alh", "echo Hello World"]
}