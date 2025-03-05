package main

import (
  "encoding/json"
  "fmt"
  "io"
  "net/http"
  "os"
  "math/rand"
  "strings"
  "time"

  "pokedex/internal/pkcache"
  "pokedex/internal/pkjson"
  "pokedex/internal/pokedex"
)

func cleanInput(text string) (split []string) {
  splitstr := strings.Fields(text)
  for _, sps := range splitstr {
    split = append(split, sps)
  }

  return
}

func commandExit([]string) error {
  fmt.Println("Closing the Pokedex... Goodbye!")

  os.Exit(0)

  return nil // Lol
}

func initCommandHelp(mep *map[string]cliCommand) (func([]string) error) {
  cmds := mep

  return func(argv []string) error {
    fmt.Println("Usage:")
    if cmds == nil {
      return fmt.Errorf("No command map provided")
    } else if len(*cmds) == 0 {
      return fmt.Errorf("No commands to print")
    }

    for _, cmd := range *cmds {
      fmt.Printf("%s: %s\n", cmd.name, cmd.description)
    }

    return nil
  }
}

func initCommandMap() (func([]string) error, func([]string) error) {
  next       := "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
  prev       := ""
  cache      := pkcache.NewPkCache(30 * time.Second)

  var locs pkjson.LocationGroupResponse

  print_locs := func() {
    for _, loc := range locs.Results {
      fmt.Println(loc.Name)
    }
  }

  nav := func(where string) error {
    var err error
    if bytes, ok := cache.Get(where); !ok {
      res, err := http.Get(where)
      if err != nil { return err }
      defer res.Body.Close()

      body, err := io.ReadAll(res.Body)
      if err != nil { return err }
      cache.Add(next, body)

      err = json.Unmarshal(body, &locs)
    } else {
      err = json.Unmarshal(bytes, &locs)
    }

    if err != nil { return err }

    next = locs.Next
    prev = locs.Previous
    print_locs()
    return nil
  }

  fwd := func(argv []string) error {
    if next == "" {
      return fmt.Errorf("You're on the last page. Use `mapb`.")
    }

    return nav(next)
  }

  bak := func(argv []string) error {
    if prev == "" {
      return fmt.Errorf("You're on the first page. Use `map`.")
    }

    return nav(prev)
  }

  return fwd, bak
}

func initExploreCmd() (func(argv []string) error) {
  cache := pkcache.NewPkCache(30 * time.Second)

  return func(argv []string) error {
    if len(argv) < 1 { return fmt.Errorf("No place to explore!") }
    area := argv[0]
    fmt.Println("Exploring:", area)
    var loc pkjson.Location
    var err error

    if data, ok := cache.Get(area); !ok {
      res, err := http.Get("https://pokeapi.co/api/v2/location-area/" + area)
      if err != nil { return err }
      defer res.Body.Close()

      body, err := io.ReadAll(res.Body)
      if err != nil { return err }
      cache.Add(area, body)

      err = json.Unmarshal([]byte(body), &loc)
    } else {
      err = json.Unmarshal(data, &loc)
    }

    if err != nil { return err }

    for _, enc := range loc.Encounters {
      fmt.Println(enc.Pokemon.Name)
    }

    return nil
  }
}

func initCommandCatch(dex *pokedex.Pokedex) (func([]string) error) {
  return func (argv []string) error {
    if len(argv) < 1 { return fmt.Errorf("Must supply pokemon name to catch") }

    res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + argv[0])
    if err != nil { return err }
    defer res.Body.Close()

    data, err := io.ReadAll(res.Body)
    if err != nil { return err }

    var pkmn pkjson.PokemonData
    const (
      base_offset = 36
      base_high   = 635 - 36 // Blissey - Sunkern
    )

    err = json.Unmarshal(data, &pkmn)
    if err != nil { return err }

    fmt.Printf("Throwing a Pokeball at %s...\n", pkmn.Name)

    randval := (rand.Intn(base_high) + base_offset)
    fmt.Println("Random val:", randval, " need:", pkmn.BaseExp)
    if randval >= pkmn.BaseExp {
      if already_caught := dex.Add(pkmn); already_caught {
        fmt.Printf("Another %s was caught!\n", pkmn.Name)
      } else {
        fmt.Printf("%s was caught!\n", pkmn.Name)
      }
    } else {
      fmt.Printf("%s escaped!\n", pkmn.Name)
    }

    return nil
  }
}

func initCommandInspect(dex *pokedex.Pokedex) (func([]string) error) {
  return func(argv []string) error {
    if len(argv) < 1 { return fmt.Errorf("No name given.") }

    name := argv[0]
    if pkmn, have := dex.Get(name); have {
      fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", pkmn.Name, pkmn.Height, pkmn.Weight)
      for _, stat := range pkmn.Stats {
        fmt.Printf("  - %s: %d\n", stat.Stat.Name, stat.Value)
      }
      fmt.Printf("Types:\n")
      for _, typ := range pkmn.Types {
        fmt.Printf("  - %s\n", typ.Type.Name)
      }
    } else {
      fmt.Println("You have not yet caught", name)
    }
    return nil
  }
}

func initCommandPokedex(dex *pokedex.Pokedex) (func([]string) error) {
  return func(argv []string) error {
    dex.List()
    return nil
  }
}
