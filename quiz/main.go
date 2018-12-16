package main

import (
    "flag"
    "fmt"
    "os"
    "bufio"
    "strings"
    "math/rand"
    "io"
    "time"
)

type Quiz struct {
    statement string
    answer string
}

type Config struct {
    timelimit int
    filename string
}

func main() {
    timelimit := 30
    filename := "problems.csv"

    config := Config{ timelimit, filename }

    csvFilename := flag.String("csv", config.filename, "a csv file in the format of 'question,answer'")
    nSecond := flag.Int("time", config.timelimit, "Timer limit for answer the question. argument should be passed in seconds.")
    flag.Parse()
    file, err := os.Open(*csvFilename)
    rand.Seed(time.Now().UnixNano())
    limit := 10 // number of quizzes

    if err != nil {
        exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
    }

    records := readQuiz(file)

    startTimer(*nSecond)

    interactQuiz(records, limit)
}

func exit(msg string) {
    fmt.Printf(msg)
    os.Exit(1)
}

func startTimer(timelimit int) {
    timer := time.NewTimer(time.Second * time.Duration(timelimit))
    go func() {
        <-timer.C
        fmt.Printf("\n時間切れ！")
        os.Exit(0)
    }()
}

func readQuiz(file io.Reader) []Quiz {
    var records []Quiz
    sc := bufio.NewScanner(file)

    for sc.Scan() {
        line := strings.Split(sc.Text(), ",")
        x := line[0]
        y := line[1]
        records = append(records, Quiz{x, y})
    }

    return records
}

func interactQuiz(records []Quiz, limits int) {
    n := int(len(records))
    nCorrect := 0
    reader := bufio.NewReader(os.Stdin)

    for _, n := range rand.Perm(n)[0:limits]{
        record := records[n]
        fmt.Printf("%s: ", record.statement)
        answer, _ := reader.ReadString('\n')
        answer = strings.Trim(answer, "\n")
        if answer == "exit" {
            break
        }
        if answer == record.answer {
            fmt.Printf("正解！\n")
            nCorrect += 1
        } else {
            fmt.Printf("不正解！ 答えは%sでした。\n", record.answer)
        }
        fmt.Printf("\n")
    }

    fmt.Printf("クイズ終了！ あなたの正解数は%d/%dで、正解率は%d%%", nCorrect, limits, nCorrect*100/limits)
}
