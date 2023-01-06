package registry

//服务信息抽象
type Service struct {
	Name  string  `json:"name"`  //服务名
	Nodes []*Node `json:"nodes"` //服务节点
}

//服务节点抽象
type Node struct {
	Id     string `json:"id"`     //节点号
	IP     string `json:"ip"`     //节点ip
	Port   int    `json:"port"`   // 端口
	Weight int    `json:"weight"` //权重
}
