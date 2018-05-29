package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

type hackCategory struct {
	ProjectId    int
	CategoryId   int
	Name         string
	Icon         string
	Color        string
	DisplayOrder int
	CreationDate string
}

type hackStage struct {
	ProjectId    int
	StageId      int
	Name         string
	Icon         string
	Color        string
	Status       int
	IsUnblocker  bool
	DisplayOrder int
}

type hackUser struct {
	UserId    int
	Username  string
	Name      string
	AvatarUrl string
	IsAdmin   bool
	IsActive  bool
}

type hackTask struct {
	ProjectId          int
	TaskId             int
	Name               string
	Description        string
	Category           hackCategory
	SubCategory        hackCategory
	Stage              hackStage
	PlannedCost        float64
	TotalCost          float64
	Priority           int
	GlobalDisplayOrder int
	DesignDisplayOrder int
	ModelNode          interface{}
	UpdateDate         string
	CreationDate       string
	Creator            interface{}
	Milestone          interface{}
	AssignedUsers      []hackUser
	CompletedSubTasks  int
	TotalSubTasks      int
	ImportanceLevel    interface{}
	HasDependencies    bool
	IsBlocked          bool
}

var apiKey string

func apiRequest(url string) []byte {
	BearerToken := "Bearer " + apiKey
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", BearerToken)
	req.Header.Add("X-AppVersion", "110")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	respJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return respJSON
}

// getTasks returns a list of tasks from hack n' plan
func getTasks() []hackTask {
	tasksJSON := apiRequest("https://app.hacknplan.com/api/v1/projects/31094/milestones/80780/tasks?categoryId=0")
	var tasks []hackTask
	err := json.Unmarshal(tasksJSON, &tasks)
	if err != nil {
		panic(err)
	}
	return tasks
}

// https://app.hacknplan.com/api/v1/projects/31094/users?includeInactive=true
func getUsers() []hackUser {
	tasksJSON := apiRequest("https://app.hacknplan.com/api/v1/projects/31094/users?includeInactive=true")
	var tasks []hackUser
	err := json.Unmarshal(tasksJSON, &tasks)
	if err != nil {
		panic(err)
	}
	return tasks
}

type hackUserCounted struct {
	UserId             int
	Username           string
	Name               string
	AvatarUrl          string
	IsAdmin            bool
	IsActive           bool
	CompletedTaskCount int
	TaskCount          int
}

func countUserTasks(users []hackUser) []hackUserCounted {
	countedUsersMap := make(map[string]hackUserCounted, len(users))
	for _, user := range users {
		userCounted := hackUserCounted{
			UserId:             user.UserId,
			Username:           user.Username,
			Name:               user.Name,
			AvatarUrl:          user.AvatarUrl,
			IsAdmin:            user.IsAdmin,
			IsActive:           user.IsActive,
			CompletedTaskCount: 0,
			TaskCount:          0,
		}
		countedUsersMap[user.Username] = userCounted
	}
	tasks := getTasks()
	for _, task := range tasks {
		for _, user := range task.AssignedUsers {
			newCU := countedUsersMap[user.Username]
			newCU.TaskCount++
			if task.Stage.Name == "Completed" {
				newCU.CompletedTaskCount++
			}
			countedUsersMap[user.Username] = newCU
		}
	}
	countedUsers := make([]hackUserCounted, 0, len(countedUsersMap))
	for _, countedUser := range countedUsersMap {
		countedUsers = append(countedUsers, countedUser)
	}
	return countedUsers
}

// BootstrapAPI connects the API to the main app's router
func BootstrapAPI(r **httprouter.Router, _apiKey string) {
	apiKey = _apiKey
	(*r).GET("/api/users", userHandler)
	(*r).GET("/api/my_tasks/:username", rootHandler)
}

// rootHandler handles the root api path (which has no functionality)
func rootHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Add("Content-Type", "application/json")
	myUsername := p.ByName("username")
	if len(myUsername) > 0 {
		tasks := getTasks()
		newTasks := []hackTask{}
		for _, task := range tasks {
			forMe := false
			for _, user := range task.AssignedUsers {
				if strings.ToLower(user.Username) == strings.ToLower(myUsername) && !forMe {
					forMe = true
				}
			}
			if forMe {
				newTasks = append(newTasks, task)
			}
		}
		data, err := json.Marshal(newTasks)
		if err != nil {
			panic(err)
		}
		fmt.Fprintln(rw, string(data))
	}
}

func userHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	rw.Header().Add("Content-Type", "application/json")
	users := getUsers()
	countedUsers := countUserTasks(users)
	data, err := json.Marshal(countedUsers)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(rw, string(data))
}
