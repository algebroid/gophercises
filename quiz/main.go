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
    isShuffled bool
}

func main() {
    defaultTimeLimit := 30
    filename := "problems.csv"
    shuffle := true

    defaultConfig := Config{ defaultTimeLimit, filename, shuffle }

    csvFilename := flag.String("csv", defaultConfig.filename, "a csv file in the format of 'question,answer'")
    timeLimit := flag.Int("time", defaultConfig.timelimit, "Timer limit for answer the question. argument should be passed in seconds.")
    isShuffled := flag.Bool("shuffle", defaultConfig.isShuffled, "Flags for whether Quizzes are shuffled or unshuffled.")
    flag.Parse()
    file, err := os.Open(*csvFilename)
    rand.Seed(time.Now().UnixNano())
    limit := 10 // number of quizzes

    if err != nil {
        exit(fmt.Sprintf("Failed to open the CSV file: %s\n", *csvFilename))
    }

    records := readQuiz(file)

    interactQuiz(records, limit, *timeLimit, *isShuffled)
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

func interactQuiz(records []Quiz, limits int, timeLimit int, isShuffled bool) {
    n := int(len(records))
    nCorrect := 0
    reader := bufio.NewReader(os.Stdin)
    var quizNumbers []int

    if isShuffled {
        quizNumbers = rand.Perm(n)[0:limits]
    } else {
        quizNumbers = makeRange(limits)
    }

    timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
    quizloop:
    for _, n := range quizNumbers {
        record := records[int(n)]
        fmt.Printf("%s: ", record.statement)

        answerCh := make(chan string)
        go func() {
            answer, _ := reader.ReadString('\n')
            answer = strings.Trim(answer, "\n")
            answerCh <- answer
        }()

        select {
        case <- timer.C:
            fmt.Println()
            break quizloop
        case answer := <-answerCh:
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
    }

    fmt.Printf("クイズ終了！ あなたの正解数は%d/%dで、正解率は%d%%", nCorrect, limits, nCorrect*100/limits)
}

func makeRange(n int) []int {
    array := make([]int, n)
    for i := range array {
        array[i] = i
    }
    return array
}
