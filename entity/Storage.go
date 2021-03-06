//provide interface of the operation on data
/* warning: without fully test */
/* remember to tell me if there are bugs */
package entity

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var (
	current_user   string
	total_user     []User
	total_meeting  []Meeting
	userLib        string = "user.json"
	meetingLib     string = "meeting.json"
	currentUserLib string = "curUser.txt"
)

func Init() {
	ReadUserFile()
	ReadMeetingFile()
	ReadCurrentUser()
}

func UpdateLib() {
	WriteUserFile()
	WriteMeetingFile()
	WriteCurrentUserFile()
}

func GetCurrentUser() string {
	return current_user
}

func SetCurrentUser(n string) {
	current_user = n
}

func GetAllUser() []User {
	return total_user
}

func GetAllMeeting() []Meeting {
	return total_meeting
}

func ReadUserFile() {
	file, err := os.Open(userLib)
	if err != nil {
		return
	}
	state, _ := file.Stat()
	if state.Size() == 0 {
		return
	}
	buffer := make([]byte, state.Size())
	_, err = file.Read(buffer)
	if err != nil {
		return
	}
	buffer = []byte(os.ExpandEnv(string(buffer)))
	err = json.Unmarshal(buffer, &total_user)
	if err != nil {
		return
	}
}

func ReadMeetingFile() {
	file, err := os.Open(meetingLib)
	if err != nil {
		return
	}
	state, _ := file.Stat()
	if state.Size() == 0 {
		return
	}
	buffer := make([]byte, state.Size())
	_, err = file.Read(buffer)
	if err != nil {
		return
	}
	buffer = []byte(os.ExpandEnv(string(buffer)))
	err = json.Unmarshal(buffer, &total_meeting)
	if err != nil {
		return
	}
}

func ReadCurrentUser() {
	file, err := os.Open(currentUserLib)
	if err != nil {
		return
	}
	state, _ := file.Stat()
	if state.Size() == 0 {
		return
	}
	buf := bufio.NewReader(file)
	line, err := buf.ReadString('\n')
	if err != nil {
		return
	}
	current_user = strings.TrimSpace(line)
}

func WriteUserFile() {
	userRec, err := json.Marshal(total_user)
	if err != nil {
		fmt.Println(err)
	}
	f, _ := os.Create(userLib)
	defer f.Close()
	f.WriteString(string(userRec))
}

func WriteMeetingFile() {
	meetingRec, err := json.Marshal(total_meeting)
	if err != nil {
		fmt.Println(err)
	}
	f, _ := os.Create(meetingLib)
	defer f.Close()
	f.WriteString(string(meetingRec))
}

func WriteCurrentUserFile() {
	f, _ := os.Create(currentUserLib)
	defer f.Close()
	f.WriteString(current_user + "\n")
}

func IsUserExist_Login(name string, psw string) bool {
	for i := 0; i < len(total_user); i++ {
		if total_user[i].Username == name && total_user[i].Password == psw {
			return true
		}
	}
	return false
}

func UserCheck(new_user User) bool {
	for i := 0; i < len(total_user); i++ {
		if total_user[i].Username == new_user.Username {
			return false
		}
	}
	return true
}

func UsernameCheck(name string) bool {
	for i := 0; i < len(total_user); i++ {
		if total_user[i].Username == name {
			return true
		}
	}
	return false
}

func CreateUser(name string, psw string, ma string, ph string) int {
	user := User{
		Username: name,
		Password: psw,
		Mail:     ma,
		Phone:    ph,
	}
	if UserCheck(user) {
		total_user = append(total_user, user)
		return 0
	} else {
		return 2
	}
}

func DeleteUser(name string) bool {
	var i int
	for i = 0; i < len(total_user); i++ {
		if total_user[i].Username == name {
			break
		}
	}
	if i == len(total_user) {
		return false
	}
	total_user[i] = total_user[len(total_user)-1]
	total_user = total_user[0 : len(total_user)-1]
	return true
}

func MeetingCheck(t string) int {
	for i := 0; i < len(total_meeting); i++ {
		if total_meeting[i].Title == t {
			return i
		}
	}
	return -1
}

func QueryMeetings(name, st, et string) []Meeting {
	var rec []Meeting
	for i := 0; i < len(total_meeting); i++ {
		if !(total_meeting[i].EndTime <= st || total_meeting[i].StartTime >= et) && (total_meeting[i].Sponsor == name || total_meeting[i].isParticipator(name) != -1) {
			rec = append(rec, total_meeting[i])
		}
	}
	return rec
}

func CheckUserFreeTime(name, st, et string) bool {
	var rec []Meeting
	for i := 0; i < len(total_meeting); i++ {
		if !(total_meeting[i].EndTime <= st || total_meeting[i].StartTime >= et) {
			rec = append(rec, total_meeting[i])
		}
	}
	if len(rec) != 0 {
		for j := 0; j < len(rec); j++ {
			if rec[j].Sponsor == name || rec[j].isParticipator(name) != -1 {
				return false
			}
		}
	}
	return true
}

func CreateMeeting(t string, s string, st string, et string, p []string) bool {
	if len(p) == 0 {
		return false
	}

	meeting := Meeting{
		Title:         t,
		Sponsor:       s,
		StartTime:     st,
		EndTime:       et,
		Participators: p,
	}
	total_meeting = append(total_meeting, meeting)
	return true
}

func DeleteMeeting(t string, name string) int {
	if UsernameCheck(name) {
		pos := MeetingCheck(t)
		if pos != -1 && total_meeting[pos].Sponsor == name {
			total_meeting[pos] = total_meeting[len(total_meeting)-1]
			total_meeting = total_meeting[0 : len(total_meeting)-1]
			return 0
		} else {
			return 1
		}
	} else {
		return 2
	}
}

func AddMeetingParticipators(title string, participator string) int {
	pos := MeetingCheck(title)
	if pos == -1 || total_meeting[pos].getSponsor() != current_user {
		return 1
	}
	if UsernameCheck(participator) && CheckUserFreeTime(participator, total_meeting[pos].getStartTime(), total_meeting[pos].getEndTime()) {
		total_meeting[pos].addParticipator(participator)
		return 0
	} else {
		return 2
	}
}

func DeleteMeetingParticipators(t string, player string) int {
	if UsernameCheck(player) {
		pos := MeetingCheck(t)
		if pos == -1  {
			return 1
		}
		if total_meeting[pos].removeParticipator(player) {
			if len(total_meeting[pos].Participators) == 0 {
				DeleteMeeting(total_meeting[pos].Title, current_user)
			}
			return 0
		} else {
			return 1
		}
	} else {
		return 2
	}
}

func ClearMeetings(name string) int {
	if UsernameCheck(name) {
		var rec []Meeting
		for i := 0; i < len(total_meeting); i++ {
			if total_meeting[i].Sponsor != name {
				rec = append(rec, total_meeting[i])
			}
		}
		total_meeting = rec
		return 0
	} else {
		return 2
	}
}

func DeleteAllMeeting(name string) int {
	if UsernameCheck(name) {
		var rec []string
		ClearMeetings(name)
		for i := 0; i < len(total_meeting); i++ {
			if total_meeting[i].isParticipator(name) != -1 {
				rec = append(rec, total_meeting[i].Title)
			}
		}
		for j := 0; j < len(rec); j++ {
			DeleteMeetingParticipators(rec[j], name)
		}
		return 0
	} else {
		return 2
	}
}
