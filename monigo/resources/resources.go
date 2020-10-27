package resources

type resourceBox struct {
	storage map[string][]byte
}

func newResourceBox() *resourceBox {
	return &resourceBox{storage: make(map[string][]byte)}
}

func (r *resourceBox) Has(file string) bool {
	_, ok := r.storage[file]
	return ok
}

func (r *resourceBox) Get(file string) []byte {
	val, ok := r.storage[file]
	if !ok {
		return []byte{}
	}

	return val
}

func (r *resourceBox) Add(file string, content []byte) {
	r.storage[file] = content
}

var resources = newResourceBox()

func Get(file string) []byte {
	return resources.Get(file)
}

func Add(file string, content []byte) {
	resources.Add(file, content)
}

func Has(file string) bool {
	return resources.Has(file)
}