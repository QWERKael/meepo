package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/memberlist"
	"github.com/pborman/uuid"
	"meepo/task/tasks"
	"meepo/util"
	"os"
	"strings"
	"sync"
	"utility-go/codec"
	"utility-go/gostruct"
)

var (
	mtx        sync.RWMutex
	items      = map[string]string{}
	broadcasts *memberlist.TransmitLimitedQueue
	mList      *memberlist.Memberlist
	members    map[string]member
)

type pkg struct {
	items   map[string]string
	members map[string]member
}

type Role uint8

const (
	UnknownNode    Role = 0
	NormalServer   Role = 1
	NormalClient   Role = 2
	ResourceServer Role = 3
)

type member struct {
	nodeName   string
	serverAddr string
	serverName string
	role       Role
	groups     gostruct.Set
	labels     gostruct.Set
}

func newMember() member {
	return member{
		serverAddr: fmt.Sprintf("%s:%d", util.Config.ListenHost, util.Config.ListenPort),
		serverName: util.Config.ActorName,
		role:       NormalServer,
		groups:     gostruct.Set{},
		labels:     gostruct.Set{},
	}
}

type broadcast struct {
	msg    []byte
	notify chan<- struct{}
}

type delegate struct{}

type update struct {
	Action string // add, del
	Data   map[string]string
}

func (b *broadcast) Invalidates(other memberlist.Broadcast) bool {
	return false
}

func (b *broadcast) Message() []byte {
	return b.msg
}

func (b *broadcast) Finished() {
	if b.notify != nil {
		close(b.notify)
	}
}

func (d *delegate) NodeMeta(limit int) []byte {
	return []byte{}
}

func (d *delegate) NotifyMsg(b []byte) {
	if len(b) == 0 {
		return
	}

	switch b[0] {
	case 'd': // data
		var updates []*update
		if err := codec.DecodeJson(b[1:], &updates); err != nil {

			return
		}
		mtx.Lock()
		for _, u := range updates {
			for k, v := range u.Data {
				switch u.Action {
				case "add":
					items[k] = v
				case "del":
					delete(items, k)
				}
			}
		}
		mtx.Unlock()
	case 'm': // 接收到member信息，将其加入到members中
		var member member
		if err := codec.DecodeJson(b[1:], &member); err != nil {
			return
		}
		mtx.Lock()
		members[member.nodeName] = member
		mtx.Unlock()
	}
}

func (d *delegate) GetBroadcasts(overhead, limit int) [][]byte {
	return broadcasts.GetBroadcasts(overhead, limit)
}

func (d *delegate) LocalState(join bool) []byte {
	var pkg pkg
	mtx.RLock()
	pkg.items = items
	pkg.members = members
	mtx.RUnlock()
	b, _ := codec.EncodeJson(pkg)
	return b
}

func (d *delegate) MergeRemoteState(buf []byte, join bool) {
	if len(buf) == 0 {
		return
	}
	if !join {
		return
	}
	var pkg pkg
	if err := codec.DecodeJson(buf, &pkg); err != nil {
		return
	}
	mtx.Lock()
	for k, v := range pkg.items {
		items[k] = v
	}
	for k, v := range pkg.members {
		members[k] = v
	}
	mtx.Unlock()
}

type eventDelegate struct{}

func (ed *eventDelegate) NotifyJoin(node *memberlist.Node) {
	util.SugarLogger.Infof("A node has joined: " + node.String())
}

func (ed *eventDelegate) NotifyLeave(node *memberlist.Node) {
	delete(members, node.Name)
	util.SugarLogger.Infof("A node has left: " + node.String())
}

func (ed *eventDelegate) NotifyUpdate(node *memberlist.Node) {
	util.SugarLogger.Infof("A node was updated: " + node.String())
}

func Set(args *tasks.Args, extraArgs []byte) ([]byte, error) {
	key := args.SubCommands[0]
	val := args.SubCommands[1]
	mtx.Lock()
	items[key] = val
	mtx.Unlock()

	b, err := codec.EncodeJson([]*update{
		{
			Action: "add",
			Data: map[string]string{
				key: val,
			},
		},
	})

	if err != nil {
		errStr := fmt.Sprintf("更新【%s : %s】失败", key, val)
		util.SugarLogger.Error(errStr)
		return []byte(errStr), err
	}

	broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})
	return []byte(key + " 已更新"), nil
}

func Del(args *tasks.Args, extraArgs []byte) ([]byte, error) {
	key := args.SubCommands[0]
	mtx.Lock()
	delete(items, key)
	mtx.Unlock()

	b, err := codec.EncodeJson([]*update{{
		Action: "del",
		Data: map[string]string{
			key: "",
		},
	}})

	if err != nil {
		errStr := fmt.Sprintf("删除【%s】失败", key)
		util.SugarLogger.Error(errStr)
		return []byte(errStr), err
	}

	broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("d"), b...),
		notify: nil,
	})
	return []byte(key + " 已删除"), nil
}

func Get(args *tasks.Args, extraArgs []byte) ([]byte, error) {
	mtx.RLock()
	val := items[args.SubCommands[0]]
	mtx.RUnlock()
	return []byte(val), nil
}

func Show(args *tasks.Args, extraArgs []byte) ([]byte, error) {
	switch strings.ToLower(args.SubCommands[0]) {
	case "nodes":
		return showNodes()
	case "members":
		return showMembers()
	case "me":
		return showMe()
	}
	return nil, errors.New("子命令不存在！")
}

func showNodes() ([]byte, error) {
	nodes := mList.Members()
	addrs := make([]memberlist.Address, 0)
	for _, n := range nodes {
		addrs = append(addrs, n.FullAddress())
	}
	b, err := codec.EncodeJson(addrs)
	if err != nil {
		util.SugarLogger.Errorf("编码时出错：%s", err.Error())
		return nil, err
	}
	return b, nil
}

func showMembers() ([]byte, error) {
	b, err := codec.EncodeJson(members)
	if err != nil {
		util.SugarLogger.Errorf("编码时出错：%s", err.Error())
		return nil, err
	}
	return b, nil
}

func showMe() ([]byte, error) {
	me := mList.LocalNode()
	b, err := codec.EncodeJson(me.FullAddress())
	if err != nil {
		util.SugarLogger.Errorf("编码时出错：%s", err.Error())
		return nil, err
	}
	return b, nil
}

func start(joinTo string) (string, error) {
	hostname, _ := os.Hostname()
	member := newMember()
	member.nodeName = hostname + "-" + uuid.NewUUID().String()
	c := memberlist.DefaultLocalConfig()
	c.Logger = util.StdLogger
	c.Events = &eventDelegate{}
	c.Delegate = &delegate{}
	c.BindPort = 0
	c.Name = member.nodeName
	var err error
	mList, err = memberlist.Create(c)
	if err != nil {
		return "", err
	}
	if len(joinTo) > 0 {
		parts := strings.Split(joinTo, ",")
		_, err := mList.Join(parts)
		if err != nil {
			return "", err
		}
	}
	broadcasts = &memberlist.TransmitLimitedQueue{
		NumNodes: func() int {
			return mList.NumMembers()
		},
		RetransmitMult: 3,
	}

	// 广播当前节点的member信息
	mtx.Lock()
	members[member.nodeName] = member
	mtx.Unlock()

	var b []byte
	b, err = codec.EncodeJson(member)

	broadcasts.QueueBroadcast(&broadcast{
		msg:    append([]byte("m"), b...),
		notify: nil,
	})

	// 返回当前节点地址
	node := mList.LocalNode()
	addr := fmt.Sprintf("%s:%d", node.Addr, node.Port)
	util.SugarLogger.Debugf("Local member %s\n", addr)
	return addr, nil
}

func Run(args *tasks.Args, extraArgs []byte) ([]byte, error) {
	var joinTo string
	if len(args.SubCommands) < 1 {
		joinTo = ""
	} else {
		joinTo = args.SubCommands[0]
	}
	//fmt.Printf("port: %d\njoinTo: %s\n", port, joinTo)
	util.SugarLogger.Debugf("启动gossip，加入到members：%s\n", joinTo)
	addr, err := start(joinTo)
	if err != nil {
		util.SugarLogger.Error(err)
		return nil, err
	}
	return []byte(addr), nil
}

//func main() {
//	runtime.GOMAXPROCS(runtime.NumCPU())
//
//	n := 2
//	chAddr <- ""
//	for i := 1; i <= n; i++ {
//		addr := <-chAddr
//		fmt.Printf("running #%d: %s\n", i, addr)
//		go Run(5000+i, addr)
//		time.Sleep(2 * time.Second)
//	}
//
//	//go Run(5002, "10.253.121.179:61740")
//	_, _ = console.ReadLine()
//}
