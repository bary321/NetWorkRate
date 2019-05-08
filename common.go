package NetWorkRate

import (
	"encoding/json"
)

//rpc使用
type Args struct {
	Interval int
}

type Common int

func (c *Common) GetRate(args *Args, rates *IORates) error {
	r, err := FastGet(false, nil, args.Interval)
	if err != nil {
		return err
	}

	d, _ := json.Marshal(r)
	json.Unmarshal(d, rates)
	return nil
}
