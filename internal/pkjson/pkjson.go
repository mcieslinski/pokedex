package pkjson

type PokemonStat struct {
  Stat  StatData `json:"stat"`
  Value int      `json:"base_stat"`
}

type PokemonType struct {
  Slot int      `json:"slot"`
  Type TypeData `json:"type"`
}

type PokemonData struct {
  Name    string        `json:"name"`
  BaseExp int           `json:"base_experience"`
  Height  int           `json:"height"`
  Weight  int           `json:"weight"`
  Stats   []PokemonStat `json:"stats"`
  Types   []PokemonType `json:"types"`
}

type TypeData struct {
  Name string `json:"name"`
}

type StatData struct {
  Name string `json:"name"`
}

type Encounter struct {
  Pokemon PokemonData `json:"pokemon"`
}

type Location struct {
  Name       string      `json:"name"`
  Encounters []Encounter `json:"pokemon_encounters"`
}

type LocationGroupResponse struct {
  Count    int        `json:"count"`
  Next     string     `json:"next"`
  Previous string     `json:"previous"`
  Results  []Location `json:"results"`
}

