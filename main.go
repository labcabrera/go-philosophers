package main

import (
  "bufio"
  "fmt"
  "os"
  "log"
  "strings"
  "flag"
  "strconv"
  "math/rand"
  "time"
)

type Philosopher struct {
  name string
  rules []PhilosopherRule
  adapter GeneticRuleAdapter
}

type PhilosopherRule struct {

  /** 0: empty
    * 1: X
    * 2: Y
    * 3: any
    *
    * 012
    * 345
    * 678
    */
  value [9]byte
  result byte
  /** weight
    */
  w float64
}

type PhilosopherGame struct {

  /** 0: empty
    * 1: player 1
    * 2: player 2
    */
  board [9]byte
  turn byte
}

type PhilosopherSimulationOptions struct {

  ruleFile1 string
  ruleFile2 string
  countRules1 int
  countRules2 int
  count int
}

type GeneticRuleAdapter interface {
  match(r *PhilosopherRule)
  unmatch(r *PhilosopherRule)
  invalid(r *PhilosopherRule)
}

// -- default adapter implementation ------------------------------------------

type DefaultGeneticRuleAdapter struct {
}

func (t DefaultGeneticRuleAdapter) match(r *PhilosopherRule) {
  r.w++
}

func (t DefaultGeneticRuleAdapter) unmatch(r *PhilosopherRule) {
  r.w--
}

func (t DefaultGeneticRuleAdapter) invalid(r *PhilosopherRule) {
  log.Println("Detected invalid gen ", r.value, r.result, r.w)
  r.w -= 100
}

// ------------------------ lets see it ---------------------------------------

func main() {
  var option = flag.String("option", "run", "Options: run | gen-rules | process-rules")
  var count = flag.Int("count", 1, "count")
  var ruleFile1 = flag.String("ruleFile1", "", "file rules player 1")
  var ruleFile2 = flag.String("ruleFile2", "", "file rules player 2")
  var filename = flag.String("filename", "", "filename")  
  flag.Parse()

  log.Println("Running option [", *option, "]")

  // Use fixed seed to test the same game
  var seed = time.Now().UTC().UnixNano()
  log.Println("Using random seed ", seed)
  rand.Seed(seed)

  switch *option {
  case "run":
    options := new(PhilosopherSimulationOptions)
    options.count = *count
    options.ruleFile1 = *ruleFile1
    options.ruleFile2 = *ruleFile2
    options.countRules1 = 10000
    options.countRules2 = 10000
    run(options)
    break
  case "gen-rules":
    genRules(*filename, *count)
    break
  case "process-rules":
    processFileRules(*filename)
    break
  }
}

func genRules(filename string, count int) {
  log.Println("Gen rules")
  p := new(Philosopher)
  initializeRandomRules(p, count)
  writeRulesToFile(p, filename)
}

func processFileRules(filename string) {
  log.Println("Processing rule files")
}

func run(options *PhilosopherSimulationOptions) {
  player01 := new(Philosopher)
  player01.name = "Nietzsche"
  player01.adapter = new(DefaultGeneticRuleAdapter)
  readRulesFromFile(player01, options.ruleFile1)

  player02 := new(Philosopher)
  player02.name = "Nietzsche"
  player02.adapter = new(DefaultGeneticRuleAdapter)
  readRulesFromFile(player02, options.ruleFile2)

  game := new(PhilosopherGame)
  for i := 0 ; i < options.count ; i++ {
    initializeGame(game)
    simulateGame(game, player01, player02)
  }
  writeRulesToFile(player01, options.ruleFile1 + ".new")
  writeRulesToFile(player02, options.ruleFile2 + ".new")
}

func initializeRandomRules(p *Philosopher, count int) {
  tmp := make([]PhilosopherRule, count, count)
  for i := 0 ; i < count ; i++ {
    inititializeRandomRule(&tmp[i])
  }
  p.rules = tmp
  log.Println("Generated random rules. Count: ", count)
}

func inititializeRandomRule(r *PhilosopherRule) {
  for i :=0 ; i < 9; i++ {
    r.value[i] = byte(rand.Intn(4))
  }
  r.w = 0.0
  r.result = byte(rand.Intn(9))
}

func initializeGame(g *PhilosopherGame) {
  for i :=0 ; i < 9; i++ {
    g.board[i] = 0
  }
  g.turn = 0
}

func simulateGame(g *PhilosopherGame, p01 *Philosopher, p02 *Philosopher) {
  log.Println("Simulating game " + p01.name + " vs " + p02.name)

  var winner byte
  var currentPlayer *Philosopher
  var playerValue byte

  for (winner == 0 || winner == 3) && g.turn < 10 {
    g.turn++

    if(g.turn % 2 == 1) {
      currentPlayer = p01
      playerValue = 1
    } else {
      currentPlayer = p02
      playerValue = 2
    }

    log.Println("Current player ", currentPlayer.name)
  
    simulateMove(g, currentPlayer, playerValue)

    log.Println("Status after move: ", g.board)
    winner = evalWinner(g)
  }  
}

func simulateMove(g *PhilosopherGame, p *Philosopher, playerValue byte) {
  var ruleResult byte
  var rule *PhilosopherRule
  for i := 0 ; i < len(p.rules) ; i++ {
    rule = &p.rules[i]
    ruleResult = evalRule(g, rule)
    if(ruleResult <= 8) {
      log.Println("Rule matches ", rule.value, rule.result, rule.w)
      g.board[ruleResult] = playerValue
      p.adapter.match(rule)
      return
    } else if(ruleResult == 128) {
      p.adapter.unmatch(rule)
    } else if(ruleResult == 129) {
      p.adapter.invalid(rule)
    }
  }
}

/** 0: no winner
  * 1: player 1
  * 2: player 2
  * 3: draw
  */
func evalWinner(g *PhilosopherGame) byte {
  v := [][]uint8 {
    { 0, 1, 2},
    { 3, 4, 5},
    { 6, 7, 8},
    { 0, 3, 6},
    { 1, 4 ,7},
    { 2, 5, 8},
    { 0, 4, 8},
    { 2, 4, 6},
  }
  var i byte
  var t byte
  log.Println("Eval board winner ", g)
  for i = 0 ; i < 8 ; i++ {
    t = g.board[v[i][0]];
    log.Println(" ... ", t, g.board[v[i][1]], g.board[v[i][2]])
    if(t > 0 && t == g.board[v[i][1]] && t == g.board[v[i][2]]) {
      log.Println("We have a winner ", t)
      return t
    }
  }
  return 0
}

/**
  * 128 rule unmatched
  * 129 inconsistent
  */
func evalRule(g *PhilosopherGame, r *PhilosopherRule) byte {
  for i := 0 ; i < 9 ; i++ {
    if(r.value[i] != 3 && r.value[i] != g.board[i]) {
      return 128
    }
  }
  if(g.board[r.result] != 0) {
    return 129
  }
  return r.result
}

// -- io utils ----------------------------------------------------------------

func writeRulesToFile(p *Philosopher, fileName string) {
  f, err := os.Create(fileName)
  check(err)
  w := bufio.NewWriter(f)  
  for i := 0 ; i < len(p.rules) ; i++ {
    for z := 0 ; z < len(p.rules[i].value) ; z++ {
      _, err = w.WriteString(fmt.Sprint(p.rules[i].value[z]))
    }
    w.WriteString("\t")
    w.WriteString(fmt.Sprint(p.rules[i].result))
    w.WriteString("\t")
    w.WriteString(fmt.Sprint(p.rules[i].w))
    w.WriteString("\n")
  }
  w.Flush()
  check(err)
}

func readRulesFromFile(p *Philosopher, fileName string) {
  log.Println("Reading file ", fileName)
  file, err := os.Open(fileName)
  check(err)
  var str string
  var rules []PhilosopherRule
  var rule *PhilosopherRule
  defer file.Close()
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    str = scanner.Text()
    rule = parseRule(str)
    rules = append(rules, *rule)
  }
  p.rules = rules
  log.Println("Readed file ", fileName)
}

func parseRule(value string) *PhilosopherRule {
  tmp := strings.Split(value, "\t")
  rule := new(PhilosopherRule)
  for i := 0 ; i < 9 ; i++ {
    rule.value[i] = tmp[0][i] - 48 // ascii 48 = 0
  }
  rule.result = tmp[1][0] - 48
  rule.w, _ = strconv.ParseFloat(tmp[2], 64)
  return rule
}

func check(e error) {
  if e != nil {
      panic(e)
  }
}
