package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func CheckError(err error) {
	if err != nil {
		//fmt.Println(err)
	}
}

func ReadQuiz(path string) [][]string {
	file, err := os.Open(path)
	CheckError(err)

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 0

	quiz, err := reader.ReadAll()
	CheckError(err)

	return quiz
}

func ReadQuestion(quiz [][]string, order int) []string {
	for ord, question := range quiz {
		if ord == order {
			return question
		}
	}

	return nil
}

func ReadAnswer() string {
	scanner := bufio.NewReader(os.Stdin)
	line, _, _ := scanner.ReadLine()

	return string(line)
}

func ReadAnswer1(answer chan<- string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _, _ := reader.ReadLine()
		answer <- string(line)
	}
}

func ReadPath() string {
	file, err := os.OpenFile("path.txt", os.O_RDWR|os.O_CREATE, 0755)
	CheckError(err)

	defer file.Close()

	var path []byte
	file.Read(path)

	if path == nil {
		file.WriteString("quiz.csv")
		return "quiz.csv"
	}

	return string(path)
}

func WritePath(path string) {
	file, err := os.OpenFile("path.txt", os.O_WRONLY, 0644)
	CheckError(err)

	defer file.Close()

	_, err = file.WriteString(path)
	CheckError(err)
}

func ShowMenu() {
	path := ReadPath()
	timer := 5
	answerchan := make(chan string)
	opt := ""
	answer := ""
	shuffle := false
	change := false

	for {
		if ReadQuiz(path) == nil {
			fmt.Println("REQUESTED CREATED FILE: quiz.csv")
			fmt.Println("PROGRAM TERMINATED")
			break
		}

		fmt.Println("= = = = = = = = = = = = QUIZ MACHINE = = = = = = = = = = = =")
		fmt.Println("")
		fmt.Println("Welcome to the quiz machine, where you can test your knowledge!")
		fmt.Println("")
		fmt.Println("E: Start Quiz")
		fmt.Println("S: Shuffle")
		fmt.Println("R: Rename File")
		fmt.Println("N: Load New Quiz")
		fmt.Println("T: Set Timer")
		fmt.Println("X: Exit")
		fmt.Println("")

		if opt == "" {
			fmt.Scanln(&opt)
		} else {
			change = true
			opt = <-answerchan
		}

		if opt == "E" {
			fmt.Println("")
			fmt.Println("You have", timer, "seconds. Ready to start?")
			fmt.Println("")
			fmt.Print("Ready (Y or N): ")

			for {
				if answer == "" && change == false {
					fmt.Scanln(&answer)
				}

				if change == true {
					answer = <-answerchan
					fmt.Println(answer)
				}

				if answer != "Y" && answer != "N" {
					fmt.Println("")
					fmt.Println("Not suported!")
					fmt.Println("")
					fmt.Print("Ready (Y or N): ")

					if change == false {
						answer = ""
					}

					continue
				} else {
					break
				}
			}

			if answer == "N" && change == false {
				opt = ""
				answer = ""
				fmt.Println("")
				continue
			} else if answer == "N" && change == true {
				fmt.Println("")
				continue
			}

			fmt.Println("")
			fmt.Println(">> QUIZ <<")
			fmt.Println("")

			answers, quiz := ShowQuiz(path, timer, answerchan, shuffle)
			ShowResult(answers, quiz)

			fmt.Println("")
			fmt.Println("E: Start Again")
			fmt.Println("X: Exit")
			fmt.Println("")

			brk := false
			for {
				choice := <-answerchan

				if choice == "X" {
					brk = true
					break
				} else if choice == "E" {
					fmt.Println("")
					break
				} else {
					fmt.Println("Invalid Input!")
					fmt.Println("")
					fmt.Println("E: Start Again")
					fmt.Println("X: Exit")
					fmt.Println("")

					continue
				}
			}

			if brk == true {
				break
			}

		} else if opt == "R" {
			if change != true {
				opt = ""
			}

			old := path
			for {
				fmt.Println("")
				fmt.Print("New Name: ")

				if change != true {
					path = ReadAnswer()
				} else {
					path = <-answerchan
				}

				if path != "" {
					os.Rename(old, path)

					WritePath(path)

					fmt.Println("")
					fmt.Print("File Renamed!")
					fmt.Println("")
					fmt.Println("")
					break
				}

				if path == "" {
					fmt.Println("")
					fmt.Print("Invalid Input!")
					fmt.Println("")
					continue
				}
			}
		} else if opt == "N" {
			if change != true {
				opt = ""
			}

			for {
				fmt.Println("")
				fmt.Print("New Quiz Name: ")

				var newQuiz string

				if change != true {
					newQuiz = ReadAnswer()
				} else {
					newQuiz = <-answerchan
				}

				if newQuiz != "" {
					path = newQuiz

					fmt.Println("")
					fmt.Print("New Quiz Loaded!")
					fmt.Println("")
					fmt.Println("")
					break
				} else {
					fmt.Println("")
					fmt.Print("Invalid Input!")
					fmt.Println("")
					continue
				}
			}
		} else if opt == "T" {
			if change != true {
				opt = ""
			}

			for {
				fmt.Println("")
				fmt.Print("New Timer: ")

				var input string

				if change != true {
					input = ReadAnswer()
				} else {
					input = <-answerchan
				}

				crv, _ := strconv.Atoi(input)

				if crv != 0 {
					timer = crv
					fmt.Println("")
					fmt.Print("New Timer Set!")
					fmt.Println("")
					fmt.Println("")
					break
				} else {
					fmt.Println("")
					fmt.Print("Invalid Input!")
					fmt.Println("")
					continue
				}
			}
		} else if opt == "S" {
			if shuffle {
				shuffle = false

				fmt.Println("")
				fmt.Print("Shuffle OFF!")
				fmt.Println("")
				fmt.Println("")
			} else {
				shuffle = true

				fmt.Println("")
				fmt.Print("Shuffle ON!")
				fmt.Println("")
				fmt.Println("")
			}

			if change != true {
				opt = ""
			}
		} else if opt == "X" {
			opt = ""
			break
		} else {
			fmt.Println("")
			fmt.Println("Not suported!")
			fmt.Println("")

			if change != true {
				opt = ""
			}
		}
	}
}

func ShowResult(answers []bool, quiz [][]string) {
	fmt.Println("")
	fmt.Println(">>> RESULTS <<<<")
	fmt.Println("")

	for i := range answers {
		fmt.Print("(", answers[i], ") --> ", quiz[i][0], " <-- (", quiz[i][1], ")", "\n")
	}

	right := 0
	for _, answer := range answers {
		if answer == true {
			right++
		}
	}
	fmt.Println("")
	fmt.Println("Total Answers:", len(answers), "Total Right:", right, "Total Wrong:", len(answers)-right)
}

func Shuffle(quiz [][]string) [][]string {
	quizShuffled := make([][]string, len(quiz))
	ind := make([]int, len(quiz))
	var indShuffled []int

	for i := range ind {
		ind[i] = i
	}

	for ind != nil { //init shuffle
		if len(ind) == 1 {
			indShuffled = append(indShuffled, ind[0])
			break
		}

		random := rand.Intn(len(ind))
		indShuffled = append(indShuffled, ind[random])

		ind = slices.Delete(ind, random, random+1)
	}

	for i := range quizShuffled {
		quizShuffled[i] = quiz[indShuffled[i]]
	}

	return quizShuffled
}

func ShowQuiz(path string, timer int, achan chan string, shuffle bool) ([]bool, [][]string) {
	quiz := ReadQuiz(path)
	trueList := make([]bool, len(quiz))
	timechan := make(chan int)
	var answer string
	var time int

	if shuffle {
		quiz = Shuffle(quiz)
	}

	ord := 0
	for {
		question := ReadQuestion(quiz, ord)

		if question != nil {
			fmt.Print(question[0], " ")
		} else {
			break
		}

		if ord == 0 {
			go SetTimer(timechan)
			go ReadAnswer1(achan)
		}

	L:
		for {
			select {
			case time = <-timechan:
				if time == timer {
					fmt.Print("\n\nTime's Up!\n")

					return trueList, quiz
				}

			case answer = <-achan:
				answer = TreatAnswer(answer)

				if question[1] == answer {
					trueList[ord] = true
					break L
				} else {
					trueList[ord] = false
					break L
				}
			}
		}

		ord++
	}

	return trueList, quiz
}

func TreatAnswer(answer string) string {
	answer = strings.ToLower(answer)

	for _, l := range answer {
		if l == ' ' {
			answer = answer[1:]
		} else {
			break
		}
	}

	for answer != "" {
		//dont be afraid of the bug, receive him.
		if answer[len(answer)-1] != ' ' {
			break
		} else {
			answer = answer[:len(answer)-1]
		} //[b][e][f][o][r][e][_]
	}

	return answer
}

func SetTimer(timer chan<- int) {
	count := 0

	for {
		time.Sleep(time.Second)
		count++

		timer <- count
	}
}

func main() {
	ShowMenu()
}
