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

func main() {
    csvFilename := flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
    flag.Parse()
    file, err := os.Open(*csvFilename)
    rand.Seed(time.Now().UnixNano())
    limit := 10 // number of quizzes

    if err != nil {
        exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
    }

    records := readQuiz(file)
    interactQuiz(records, limit)
}

func exit(msg string) {
    fmt.Printf(msg)
    os.Exit(1)
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
        text, _ := reader.ReadString('\n')
        text = strings.Trim(text, "\n")
        if text == "exit" {
            break
        }
        if text == record.answer {
            fmt.Printf("正解！\n")
            nCorrect += 1
        } else {
            fmt.Printf("不正解！ 答えは%sでした。\n", record.answer)
        }
        fmt.Printf("\n")
    }

    fmt.Printf("クイズ終了！ あなたの正解数は%d/%dで、正解率は%d%%", nCorrect, limits, nCorrect*100/limits)
}
