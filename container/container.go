package container

type Map map[string]interface{}

func (m Map) Add(key string, value interface{}) {
	m[key] = value
}

func (m Map) GetString(key string) (string, bool) {
	v, ok := m[key]
	if !ok {
		return "", false
	}
	if s, ok := v.(string); ok {
		return s, true
	}
	return "", false
}

func (m Map) GetSliceString(key string) ([]string, bool) {
	v, ok := m[key]
	if !ok {
		return []string{}, false
	}
	if sliceString, ok := v.([]string); ok {
		return sliceString, true
	}
	return []string{}, false
}

func (m Map) Exclude(keys []string) Map {
	hash := make(map[string]struct{})
	for _, k := range keys {
		hash[k] = struct{}{}
	}
	res := make(Map)
	for k, v := range m {
		if _, ok := hash[k]; !ok {
			res[k] = v
		}
	}
	return res
}

func (m Map) Include(keys []string) Map {
	hash := make(map[string]struct{})
	for _, k := range keys {
		hash[k] = struct{}{}
	}
	res := make(Map)
	for k, v := range m {
		if _, ok := hash[k]; ok {
			res[k] = v
		}
	}
	return res
}

func (m Map) AppendSliceString(key string, value []string) Map {
	v, ok := m[key]
	if !ok {
		m[key] = value
		return m
	}
	if sliceString, ok := v.([]string); ok {
		m[key] = append(sliceString, value...)
		return m
	}
	return m
}

func (m Map) Keys() []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m Map) IsEmpty() bool {
	return len(m.Keys()) == 0
}
