package pokedex

import (
  "fmt"
  "sync"

  "pokedex/internal/pkjson"
)

type Pokedex struct {
  dex map[string]pkjson.PokemonData
  mux sync.RWMutex
}

func (pkdx *Pokedex) Add(pkmn pkjson.PokemonData) (exists bool) {
  pkdx.mux.Lock()
  defer pkdx.mux.Unlock()

  _, exists = pkdx.dex[pkmn.Name]

  pkdx.dex[pkmn.Name] = pkmn

  return
}

func (pkdx *Pokedex) Get(name string) (pkmn pkjson.PokemonData, exists bool) {
  pkdx.mux.RLock()
  defer pkdx.mux.RUnlock()

  pkmn, exists = pkdx.dex[name]

  return
}

func (pkdx *Pokedex) List() {
  pkdx.mux.RLock()
  defer pkdx.mux.RUnlock()

  fmt.Println("Your Pokedex:")
  for name, _ := range pkdx.dex {
    fmt.Printf("  - %s\n", name)
  }
}

func NewPokedex() *Pokedex {
  return &Pokedex{ dex: map[string]pkjson.PokemonData{}, mux: sync.RWMutex{} }
}
