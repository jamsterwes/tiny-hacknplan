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

// getTasks returns a list of tasks from hack n' plan
func getTasks() []hackTask {
	BearerToken := "Bearer " + apiKey
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://app.hacknplan.com/api/v1/projects/31094/milestones/80780/tasks?categoryId=0", nil)
	req.Header.Add("Authorization", BearerToken)
	req.Header.Add("X-AppVersion", "110")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	tasksJSON, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var tasks []hackTask
	err = json.Unmarshal(tasksJSON, &tasks)
	if err != nil {
		panic(err)
	}
	return tasks
}

// BootstrapAPI connects the API to the main app's router
func BootstrapAPI(r **httprouter.Router, _apiKey string) {
	apiKey = _apiKey
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
