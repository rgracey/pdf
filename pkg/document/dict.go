package document

type Dict map[string]interface{}

func (d Dict) Get(key string) interface{} {
	if _, ok := d[key]; !ok {
		return nil
	}

	return d[key]
}
