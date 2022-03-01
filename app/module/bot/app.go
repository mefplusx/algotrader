package bot

import (
	"encoding/json"
	"fmt"
	"os"
	commonEntity "robot/common/entity"
	"robot/module/bot/interfaces"
)

type App struct {
	Adviser          interfaces.Adviser      `json:"-"`
	BotFileStatePath string                  `json:"-"`
	StartBotFromMP   commonEntity.MomentPath `json:"-"`

	InMP *commonEntity.MomentPath `json:"in"`
}

func (a *App) saveState() {
	content, _ := json.Marshal(a)
	file, _ := os.Create(a.BotFileStatePath + "/state.json")
	file.Write(content)
	file.Close()
}

func (a *App) loadState() {
	b, err := os.ReadFile(a.BotFileStatePath + "/state.json")
	if err == nil {
		json.Unmarshal(b, a)
	}
}

func (a *App) Init() {
	a.loadState()
}

// проверить, или нужно что-то делать
func (a *App) Do() {
	transactionType, momentPath := a.Adviser.GetLastIn()
	fmt.Println(transactionType, momentPath)
}
