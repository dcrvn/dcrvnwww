package helper

type (
	IdManager struct {
		//sys *sync.Mutex
		ids []int
		check	map[int]bool
	}
)

func NewIdManager() *IdManager {
	return &IdManager{
		check: make(map[int]bool),
		ids: []int{},
	}
}

func (i *IdManager) Set(ids ...int)  {
	/*i.sys.Lock()
	defer i.sys.Unlock()*/
	for _,id := range ids {
		if _,ok := i.check[id];!ok {
			i.check[id] = true
			i.ids = append(i.ids, id)
		}
	}
}

func (i *IdManager) Get() []int {
	return i.ids
}