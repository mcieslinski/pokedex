package main

import (
  "bufio"
  "fmt"
  "os"
  "strings"

  "pokedex/internal/pokedex"
)

type cliCommand struct {
  name        string
  description string
  callback    func([]string) error
}

func main() {
  scn := bufio.NewScanner(os.Stdin)

  dex := pokedex.NewPokedex()

  // Map of commands
  cmdlist := map[string] cliCommand {
    "exit": {
      name:        "exit",
      description: "Exit the Pokedex",
      callback:    commandExit,
    },
  }

  // Register help command
  helpcmd := initCommandHelp(&cmdlist)
  cmdlist["help"] = cliCommand {
    name:        "help",
    description: "Displays a help message",
    callback:    helpcmd,
  }

  // Register map command
  fwdcmd, bakcmd := initCommandMap()
  cmdlist["map"] = cliCommand {
    name:        "map",
    description: "Display map areas going forward. 20 at a time.",
    callback:    fwdcmd,
  }
  cmdlist["mapb"] = cliCommand {
    name:        "mapb",
    description: "Display map areas going backward. 20 at a time.",
    callback:    bakcmd,
  }

  explorecmd := initExploreCmd()
  cmdlist["explore"] = cliCommand {
    name:        "explore",
    description: "Checks pokemon in a given area, usage: explore <area>",
    callback:    explorecmd,
  }

  catchcmd := initCommandCatch(dex)
  cmdlist["catch"] = cliCommand {
    name:        "catch",
    description: "Attempts to catch a given pokemon",
    callback:    catchcmd,
  }
  
  inspectcmd := initCommandInspect(dex)
  cmdlist["inspect"] = cliCommand {
    name:        "inspect",
    description: "Shows data of a caught pokemon",
    callback:    inspectcmd,
  }

  pokedexcmd := initCommandPokedex(dex)
  cmdlist["pokedex"] = cliCommand {
    name:        "pokedex",
    description: "Shows all caught pokemon",
    callback:    pokedexcmd,
  }

  fmt.Println("Welcome to the Pokedex!")
  
  for {
    fmt.Print("Pokedex > ")
    scn.Scan()
    cmdstr := make([]string, len(scn.Text()))
    if (len(scn.Text()) > 0) {
      for idx, str := range cleanInput(scn.Text()) {
        cmdstr[idx] = strings.ToLower(str)
      }
    } else {
      continue
    }

    cmd, exists := cmdlist[cmdstr[0]]
    if !exists {
      fmt.Println("Unknown command")
    } else {
      if err := cmd.callback(cmdstr[1:]); err != nil {
        fmt.Println(err.Error())
      }
    }
  }
}
