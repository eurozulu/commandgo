package misc

type StringSet interface {
	Get() []string
	Add(sz ...string)
}

type stringSet struct {
	ss map[string]bool
}

func (s stringSet) Get() []string {
	sz := make([]string, len(s.ss))
	var i int
	for s := range s.ss {
		sz[i] = s
		i++
	}
	return sz
}

func (s *stringSet) Add(sz ...string) {
	for _, ks := range sz {
		if s.ss[ks] {
			continue
		}
		s.ss[ks] = true
	}
}

func NewStringSet(sz ...string) *stringSet {
	ss := &stringSet{ss: map[string]bool{}}
	ss.Add(sz...)
	return ss
}
