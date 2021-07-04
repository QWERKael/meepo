package main

import (
	"github.com/shirou/gopsutil/net"
	"meepo/task/tasks"
	"meepo/util"
	"utility-go/codec"
)

//获取网卡信息
func Net(args *tasks.Args, extraArgs []byte) ([]byte, error) {
	util.SugarLogger.Debugf("执行show net命令...")

	type netInfo struct {
		Name string `json:"name"`
		Mac  string `json:"mac"`
		Ip   string `json:"ip"`
	}

	inter, _ := net.Interfaces()
	nis := make([]netInfo, 0)

	for i := 0; i < len(inter); i++ {
		ni := netInfo{}
		ni.Name = inter[i].Name
		if inter[i].HardwareAddr == "" {
			ni.Mac = "无"
		} else {
			ni.Mac = inter[i].HardwareAddr
		}
		//fmt.Println("IP：", inter[i].Addrs)
		var ipStr string
		for _, ip := range inter[i].Addrs {
			ipStr += ip.Addr + "   "
		}
		ni.Ip = ipStr
		nis = append(nis, ni)
	}
	b, err := codec.EncodeJson(nis)
	if err != nil {
		util.SugarLogger.Errorf("编码时出错：%s", err.Error())
		return nil, err
	}
	return b, nil
}

func Info(args *tasks.Args, extraArgs []byte) ([]byte, error) {
	info := struct {
		ServerHost string `json:"server host"`
		ServerPort int    `json:"server port"`
		ServerName string `json:"actor name"`
	}{util.Config.ListenHost,
		util.Config.ListenPort,
		util.Config.ActorName,
	}

	b, err := codec.EncodeJson(info)
	if err != nil {
		util.SugarLogger.Errorf("编码时出错：%s", err.Error())
		return nil, err
	}
	return b, nil
}

//获取监听信息
//func Net(task tasks.Task) (string, error) {
//	log.Debugf("task信息: %#v", task)
//	l := utils.Lines{}
//	inter, _ := net.stat
//
//	for i := 0; i < len(inter); i++ {
//		l.LineAppend("----------")
//		l.LineAppend("网卡名：%s", inter[i].Name)
//		if inter[i].HardwareAddr == "" {
//			l.LineAppend("MAC：%s", "无")
//		} else {
//			l.LineAppend("MAC：%s", inter[i].HardwareAddr)
//		}
//		//fmt.Println("IP：", inter[i].Addrs)
//		var ipStr string
//		for _, ip := range inter[i].Addrs {
//			ipStr += ip.Addr + "   "
//		}
//		l.LineAppend("IP: %s", ipStr)
//	}
//	return l.String(), nil
//}

//获取负载信息
//func Load(task tasks.Task) (string, error) {
//	log.Debugf("task信息: %#v", task)
//	l := utils.Lines{}
//	info, _ := load.Avg()
//	l.LineAppend(" 1分钟内负载\t%6.3f", info.Load1)
//	l.LineAppend(" 5分钟内负载\t%6.3f", info.Load5)
//	l.LineAppend("15分钟内负载\t%6.3f", info.Load15)
//	return l.String(), nil
//}
//
//type Process struct {
//	Pid        int32
//	Name       string
//	UserName   string
//	CPUPercent float64
//	MemPercent float32
//}
//
//func (p Process) String() string {
//	s := fmt.Sprintf(" %6d | %20s | %20s | %7.4f%% | %7.4f%% |",
//		p.Pid, p.Name, p.UserName, p.CPUPercent, p.MemPercent)
//	//fmt.Println(s)
//	return s
//}
//
//type Processes struct {
//	ProcList []Process
//	TotalCPU float64
//	TotalMem float32
//	by       func(p, q *Process) bool
//}
//
//func (p Processes) Len() int {
//	return len(p.ProcList)
//}
//
//func (p Processes) Swap(i, j int) {
//	p.ProcList[i], p.ProcList[j] = p.ProcList[j], p.ProcList[i]
//}
//
//func (p Processes) Less(i, j int) bool {
//	return p.by(&p.ProcList[i], &p.ProcList[j])
//}
//
//func (p *Processes) New() {
//	log.Debugf("获取进程信息")
//	processes, err := process.Processes()
//	utils.CheckErrorPanic(err)
//	p.TotalCPU = 0.0
//	p.TotalMem = float32(0.0)
//	for _, proc := range processes {
//		name, _ := proc.Name()
//		user, _ := proc.Username()
//		cpuPercent, _ := proc.CPUPercent()
//		p.TotalCPU += cpuPercent
//		memPercent, _ := proc.MemoryPercent()
//		p.TotalMem += memPercent
//		p.ProcList = append(p.ProcList, Process{Pid: proc.Pid, Name: name, UserName: user, CPUPercent: cpuPercent, MemPercent: memPercent})
//	}
//}
//
//func (p *Processes) SortBy(sortField string) error {
//	log.Debugf("开始排序")
//	switch sortField {
//	case "mem":
//		p.by = func(p, q *Process) bool {
//			return p.MemPercent > q.MemPercent
//		}
//	case "cpu":
//		p.by = func(p, q *Process) bool {
//			return p.CPUPercent > q.CPUPercent
//		}
//	default:
//		return errors.New(fmt.Sprintf("不支持的排序规则: %s", sortField))
//	}
//	sort.Sort(*p)
//	return nil
//}
//
//func (p *Processes) String(limit int) string {
//	l := utils.Lines{}
//	l.LineAppend(" %6s | %20s | %20s | %8s | %8s |", "PID", "NAME", "USER", "CPU", "MEM")
//	for i, proc := range p.ProcList {
//		if i >= limit && limit != 0 {
//			break
//		}
//		l.LineAppend(proc.String())
//	}
//	l.LineAppend("总CPU使用量: %7.4f%%", p.TotalCPU)
//	l.LineAppend("总Mem使用量: %7.4f%%", p.TotalMem)
//	return l.String()
//}
//
//func (p *Processes) toTable(limit int) *pb.CommonCmdReply {
//	reply := &pb.CommonCmdReply{
//		ResultTable: &pb.Table{
//			Header: &pb.Row{},
//			Footer: &pb.Row{},
//			Body:   make([]*pb.Row, limit),
//		},
//		Status: pb.CommonCmdReply_Ok,
//	}
//	reply.ResultTable.Header.Row = []string{"PID", "NAME", "USER", "CPU", "MEM"}
//	reply.ResultTable.Footer.Row = []string{"", "", "", fmt.Sprintf("%f", p.TotalCPU), fmt.Sprintf("%f", p.TotalMem)}
//
//	for i, _ := range reply.ResultTable.Body {
//		log.Debugf("开始填充")
//		reply.ResultTable.Body[i] = &pb.Row{Row: []string{
//			fmt.Sprintf("%d", p.ProcList[i].Pid),
//			p.ProcList[i].Name,
//			p.ProcList[i].UserName,
//			fmt.Sprintf("%f", p.ProcList[i].CPUPercent),
//			fmt.Sprintf("%f", p.ProcList[i].MemPercent),
//		}}
//	}
//	return reply
//}
//
//// 查看进程信息
//func Processlist(task tasks.Task) (*pb.CommonCmdReply, error) {
//	log.Debugf("task信息: %#v", task)
//	//l := utils.Lines{}
//	//processes, _ := process.Processes()
//	//totalCPU := 0.0
//	//totalMem := float32(0.0)
//	//for _, p := range processes {
//	//	name, _ := p.Name()
//	//	user, _ := p.Username()
//	//	cpuPercent, _ := p.CPUPercent()
//	//	totalCPU += cpuPercent
//	//	memPercent, _ := p.MemoryPercent()
//	//	totalMem += memPercent
//	//	l.LineAppend("pid:%6d\tname: %12s\tuser: %20s\tcpu: %7.4f%%\tmem: %7.4f%%\n",
//	//		p.Pid, name, user, cpuPercent, memPercent)
//	//}
//	//l.LineAppend("总CPU使用量: %7.4f%%", totalCPU)
//	//l.LineAppend("总Mem使用量: %7.4f%%", totalMem)
//	//return l.String(), nil
//	limit := utils.StringDefaultInt(task.Args["limit"], 10)
//	sortby := utils.StringDefault(task.Args["sortby"], "mem")
//	p := Processes{}
//	p.New()
//	err := p.SortBy(sortby)
//	if err != nil {
//		return nil, err
//	}
//	log.Debugln("processlist 已排序")
//	return p.toTable(limit), nil
//}
