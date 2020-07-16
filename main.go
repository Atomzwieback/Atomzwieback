package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type badges struct {
	Bronze int `json:"bronze"`
	Silver int `json:"silver"`
	Gold   int `json:"gold"`
}

type user struct {
	Badges           badges `json:"badge_counts"`
	AccountID        int    `json:"account_id"`
	Employee         bool   `json:"is_employee"`
	LastModifiedDate int    `json:"last_modified_date"`
	LastAccessDate   int    `json:"last_access_date"`
	RepYear          int    `json:"reputation_change_year"`
	RepQuater        int    `json:"reputation_change_quarter"`
	RepMonth         int    `json:"reputation_change_month"`
	RepWeek          int    `json:"reputation_change_week"`
	RepDay           int    `json:"reputation_change_day"`
	RepOverall       int    `json:"reputation"`
	CreationDate     int    `json:"creation_date"`
	UserType         string `json:"user_type"`
	UserId           int    `json:"user_id"`
	Location         string `json:"location"`
	WebsiteUrl       string `json:"website_url"`
	ProfileUrl       string `json:"link"`
	ProfileImage     string `json:"profile_image"`
	ProfileName      string `json:"display_name"`
}

type stackExchangeAnswer struct {
	Items          []user `json:"items"`
	More           bool   `json:"has_more"`
	MaxQuota       int    `json:"quota_max"`
	RemainingQuota int    `json:"quota_remaining"`
}

func main() {
	response, err := http.Get("https://api.stackexchange.com/2.2/users?order=desc&sort=reputation&inname=atomzwieback&site=stackoverflow")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatal(err)
	}

	var answer stackExchangeAnswer
	err = json.Unmarshal([]byte(responseData), &answer)

	if err != nil {
		fmt.Println(err)
	}

	writeStackoverflowUser(answer.Items[0])
}

func writeStackoverflowUser(userData user) {
	l, err := os.Create("README.md")

	if err != nil {
		deleteFile("README.md")
		writeStackoverflowUser(userData)
	}

	_, _ = l.WriteString("### Stackoverflow Profile:")
	_, _ = l.WriteString("\n##### General ProfileData")
	_, _ = l.WriteString("\n<!-- profile starts -->")
	_, _ = l.WriteString(fmt.Sprintf("\n ![](%s)", userData.ProfileImage))
	_, _ = l.WriteString("\n")
	_, _ = l.WriteString(fmt.Sprintf("\n| ** Username: **  [Atomzwieback](%s \"%s\")", userData.ProfileUrl, userData.ProfileName))
	_, _ = l.WriteString(fmt.Sprintf("\n| ** Last seen: **  %s ", time.Unix(int64(userData.LastAccessDate), 0)))
	_, _ = l.WriteString(fmt.Sprintf("\n| ** User since: **  %s ", time.Unix(int64(userData.CreationDate), 0)))
	_, _ = l.WriteString("\n")
	_, _ = l.WriteString("\n| ** Badges **")
	_, _ = l.WriteString("\n")
	_, _ = l.WriteString("\n| Gold | Silver | Bronze |")
	_, _ = l.WriteString("\n| :------------ | :------------ | :------------ |")
	_, _ = l.WriteString(fmt.Sprintf("\n| %d | %d | %d |", userData.Badges.Gold, userData.Badges.Silver, userData.Badges.Bronze ))
	_, _ = l.WriteString("\n")
	_, _ = l.WriteString("\n| ** Reputation Change **")
	_, _ = l.WriteString("\n")
	_, _ = l.WriteString("\n| Overall | Year | Quarter | Month |  Week | Day |")
	_, _ = l.WriteString("\n| :------------ | :------------ | :------------ | :------------ | :------------ | :------------ |")
	_, _ = l.WriteString(fmt.Sprintf("\n| %d | %d | %d | %d | %d | %d |", userData.RepOverall, userData.RepYear, userData.RepQuater, userData.RepMonth, userData.RepWeek, userData.RepDay ))
	_, _ = l.WriteString("\n<!-- profile ends -->")
}

func writeLine(f os.File, text string) bool {
	_, err := f.WriteString(text)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func generateHeading(headingSize int, headingText string) string {
	sizes := [6]int{1, 2, 3, 4, 5, 6}
	inSlice, err := sizeInSlice(headingSize, sizes)

	if err != nil {
		log.Fatal("Selected heading size is not in range of 1 - 6.")
	}

	if inSlice {
		return strings.Repeat("#", headingSize) + " " + headingText
	}
	return ""
}

func sizeInSlice(size int, sizes [6]int) (bool, error) {
	for _, availableSize := range sizes {
		if availableSize == size {
			return true, nil
		}
	}
	return false, errors.New("given size is not in slice")
}

func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) {
		return
	}

	fmt.Println("File Deleted")
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
