package MapController

type MapManager struct {
	depth 	int 		//节点深度
	isLeaf 	bool 		//是否是叶子节点
	region  Region		//区域范围
	LU		*MapManager //左上子节点指针
	LB 		*MapManager //左下子节点指针
	RU      *MapManager //右上子节点指针
	RB 		*MapManager //右下子节点指针
	objNum  int 		//单位元素的数量
	object 	MapObj		  //单位元素
}

type Region struct {
	Up  float64
	Bottom float64
	Left   float64
	Right  float64
}

type MapObj struct {
	row int
	col int
	obj interface{}
}