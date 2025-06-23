package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	EmailPostfix = [4]string{"@gmail.com", "@hotmail.com", "@yahoo.com", "@zoho.com"}
)

type userInfo struct {
	id       int
	username string
	email    string
}

type users struct{ users []userInfo }

func (u users) getRandomUserName(id int, random *randomNumbers) string {

	return fmt.Sprintf("%s%s%s", strconv.Itoa(id)[5:], strconv.Itoa(int(time.Now().UnixMilli()))[:6], EmailPostfix[random.random.Intn(len(EmailPostfix))])
}

func (u users) getRandomEmail(id int, random *randomNumbers) string {

	return fmt.Sprintf("%d%s", id, EmailPostfix[random.random.Intn(len(EmailPostfix))])
}

func (u *users) getUsers(random *randomNumbers) {
	// Initiate Random lib Extentions we wrote
	var n int = rand.Int()
	var seed int = random.getSeed()
	random.getRand(seed)

	id := random.random.Intn(n)
	user := userInfo{
		id:       id,
		username: u.getRandomUserName(id, random),
		email:    u.getRandomEmail(id, random),
	}
	u.users = append(u.users, user)
	fmt.Printf("new User Created \nID: %d\nUsername: %s\nEmail: %s\n\n%s\n\n", user.id, user.username, user.email, repeateString("-", 50))

	memoryUsed, err := getMemoryUsageMB()
	if err != nil {
		errMessage := fmt.Sprintf("Could not get Memory Usage: %s\n", err)
		fmt.Println(errMessage)
	}
	if memoryUsed != -1 {
		fmt.Printf("Memory Usage while generating user with ID: %d is: %.2fMB\n", id, memoryUsed)
	}
}

func repeateString(s string, num int) string {
	var finalString []string
	for range num {
		finalString = append(finalString, s)
	}
	return strings.Join(finalString, "")
}

func getProgramProcessID(name string) (int, error) {
	var processName string = name

	if len(name) == 0 {
		fmt.Println("ProcessName is required, Since it was not Provided will use ./main as Process Name!")
		processName = "main"
	}

	cmd := exec.Command("pgrep", processName)
	output, err := cmd.Output()

	if err != nil {
		if err.Error() == "exit status 1" {
			errMessage := fmt.Sprintf("Process '%s' not found", processName)
			panic(errMessage)
		} else {
			errMessage := fmt.Sprintf("Error executing pgrep command: %s", err)
			return -1, errors.New(errMessage)
		}
	}

	pidStr := strings.TrimSpace(string(output))
	pid, err := strconv.Atoi(pidStr)

	if err != nil {
		return -1, errors.New("Error Converting PID value")
	}
	return pid, nil
}

func getMemoryUsageMB() (float32, error) {
	pid, err := getProgramProcessID("main")
	if err != nil {
		errMessage := fmt.Sprintf("Error getting PID: %s\n", err)
		return -1, errors.New(errMessage)
	}

	cmdStr := fmt.Sprintf("/proc/%d/status", pid)
	cmd := exec.Command("grep", "VmRSS", cmdStr)
	output, err := cmd.Output()
	if err != nil {
		errMessage := fmt.Sprintf("Error retriving memory usage via cmd: %s\n", err)
		return -1, errors.New(errMessage)
	}
	outputStr := string(output)

	fields := strings.Fields(outputStr)

	if len(fields) < 2 {
		errMessage := fmt.Sprintf("Error unexpected format: %s\n", outputStr)
		return -1, errors.New(errMessage)
	}

	memoryUsageKb, err := strconv.ParseFloat(fields[1], 32)
	if err != nil {
		errMessage := fmt.Sprintf("Error Converting cmd output from str to int: %s\n", err)
		return -1, errors.New(errMessage)
	}

	return float32(memoryUsageKb / 1024), nil
}

func processUsers(totalNum int) {
	random := randomNumbers{}
	users := users{}
	for i := range totalNum {
		users.getUsers(&random)
		if i == totalNum-1 {
			memoryUsed, err := getMemoryUsageMB()
			if err != nil {
				errMessage := fmt.Sprintf("Could not get Memory Usage: %s\n", err)
				fmt.Println(errMessage)
			}
			if memoryUsed != -1 {
				fmt.Printf("Memory Usage while generating last user: %.2fMB\n", memoryUsed)
			}
		}
	}
	fmt.Printf("Total number of random numbers generated: %d \n", len(users.users))
}

func main() {
	start := time.Now()

	totalNum := 50000000
	processUsers(totalNum)

	fmt.Printf("Total Time: %.5f seconds \n", time.Since(start).Seconds())
}
