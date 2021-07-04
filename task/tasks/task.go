package tasks

import (
	"utility-go/codec"
)

type Task interface {
	UUID() string
	EncodeTask() ([]byte, error)
	DecodeTask([]byte) error
	SetResult(rst []byte)
	GetResult() []byte
	Exec() error
}

type MeepoTask struct {
	uuid       string
	PluginName string
	Command    string
	Args       *Args  // 任务参数
	ExtraArgs  []byte // 扩展的参数，一个json格式的[]byte对象
	Result     []byte // 任务结果，一个json格式的[]byte对象
}

type Args struct {
	SubCommands []string
	Flags       []string
	KVs         map[string][]string
}

func NewMeepoTask(pluginName string, command string, args *Args) *MeepoTask {
	return &MeepoTask{
		PluginName: pluginName,
		Command:    command,
		Args:       args,
		Result:     nil,
	}
}

func (mt MeepoTask) UUID() string {
	return mt.uuid
}

func (mt *MeepoTask) EncodeTask() ([]byte, error) {
	json, err := codec.EncodeJson(*mt)
	if err != nil {
		return nil, err
	}
	return json, nil
}

func (mt *MeepoTask) DecodeTask(json []byte) error {
	if err := codec.DecodeJson(json, mt); err != nil {
		return err
	}
	return nil
}

func (mt *MeepoTask) SetResult(rst []byte) {
	mt.Result = rst
}

func (mt *MeepoTask) GetResult() []byte {
	return mt.Result
}

//func (mt *MeepoTask) Exec() error {
//	var err error
//	mt.Result, err = executor.Exec(mt)
//	if err != nil {
//		return err
//	}
//	return nil
//}

//type LuaCommandTask struct {
//	uuid    string
//	Command string
//	Result  []byte // 任务结果，一个json格式的[]byte对象
//}
//
//func NewLuaCommandTask() LuaCommandTask {
//	return LuaCommandTask{
//		uuid:    uuid.NewString(),
//		Command: "",
//		Result:  nil,
//	}
//}
//
//func (lct LuaCommandTask) UUID() string {
//	return lct.uuid
//}
//
//func (lct *LuaCommandTask) EncodeTask() ([]byte, error) {
//	json, err := codec.EncodeJson(*lct)
//	if err != nil {
//		return nil, err
//	}
//	return json, nil
//}
//
//func (lct *LuaCommandTask) DecodeTask(json []byte) error {
//	if err := codec.DecodeJson(json, lct); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (lct *LuaCommandTask) SetResult(rst []byte) {
//	lct.Result = rst
//}
//
//func (lct *LuaCommandTask) GetResult() []byte {
//	return lct.Result
//}
//
//func (lct *LuaCommandTask) Exec() error {
//	golua.InitLoad()
//	err := golua.ExecuteLuaCommand(golua.L, lct.Command)
//	if err != nil {
//		lct.Result = []byte("执行失败：" + err.Error())
//		return err
//	}
//	lct.Result = []byte("执行成功")
//	return nil
//}
//
//type LuaScriptTask struct {
//	uuid       string
//	ScriptName string
//	Result     []byte // 任务结果，一个json格式的[]byte对象
//}
//
//func NewLuaScriptTask() LuaScriptTask {
//	return LuaScriptTask{
//		uuid:       uuid.NewString(),
//		ScriptName: "",
//		Result:     nil,
//	}
//}
//
//func (lst LuaScriptTask) UUID() string {
//	return lst.uuid
//}
//
//func (lst *LuaScriptTask) EncodeTask() ([]byte, error) {
//	json, err := codec.EncodeJson(*lst)
//	if err != nil {
//		return nil, err
//	}
//	return json, nil
//}
//
//func (lst *LuaScriptTask) DecodeTask(json []byte) error {
//	if err := codec.DecodeJson(json, lst); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (lst *LuaScriptTask) SetResult(rst []byte) {
//	lst.Result = rst
//}
//
//func (lst *LuaScriptTask) GetResult() []byte {
//	return lst.Result
//}
//
//func (lst *LuaScriptTask) Exec() error {
//	golua.InitLoad()
//	path := filepath.Join(util.Config.LuaScriptsPathAbs, lst.ScriptName)
//	rst, err := golua.ExecuteLuaScriptWithResult(golua.L, path)
//	if err != nil {
//		return err
//	}
//	lst.Result = []byte(fmt.Sprintf("%#v", rst))
//	return nil
//}
