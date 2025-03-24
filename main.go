package main

import (
	"flag"
	"zlatolas/projectManager/tui"
)


func main() {
  var repo = flag.String("repo", "Default", "the name of the repo")
  var user = flag.String("user", "Default", "your username on github")
  flag.Parse()

  tui.InitTui(*repo, *user);

}
