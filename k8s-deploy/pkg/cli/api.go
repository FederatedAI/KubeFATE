package cli

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	LIST int = iota
	INFO
	ERROR
	MSG
	JOB
)

type Item interface {
	getRequestPath() (Path string)
	addArgs() (Args string)
	getResult(Type int) (result interface{}, err error)
	outPut(result interface{}, Type int) error
}

func PostItem(i Item, Body []byte) error {
	req := &Request{
		Type: "POST",
		Path: i.getRequestPath(),
		Body: Body,
	}

	rep, err := Send(req)
	if err != nil {
		return err
	}
	if rep.Code != 200 {
		msg, err := i.getResult(ERROR)
		if err != nil {
			return err
		}
		err = json.Unmarshal(rep.Body, &msg)
		if err != nil {
			return err
		}
		err = i.outPut(msg, ERROR)
		if err != nil {
			return err
		}
		return nil
	}
	msg, err := i.getResult(JOB)

	err = json.Unmarshal(rep.Body, &msg)
	if err != nil {
		return err
	}

	err = i.outPut(msg, JOB)
	if err != nil {
		return err
	}
	return nil
}
func PutItem(i Item, Body []byte) error {
	req := &Request{
		Type: "PUT",
		Path: i.getRequestPath(),
		Body: Body,
	}

	rep, err := Send(req)
	if err != nil {
		return err
	}
	if rep.Code != 200 {
		msg, err := i.getResult(ERROR)
		if err != nil {
			return err
		}
		err = json.Unmarshal(rep.Body, &msg)
		if err != nil {
			return err
		}
		err = i.outPut(msg, ERROR)
		if err != nil {
			return err
		}
		return nil
	}
	msg, err := i.getResult(JOB)

	err = json.Unmarshal(rep.Body, &msg)
	if err != nil {
		return err
	}

	err = i.outPut(msg, JOB)
	if err != nil {
		return err
	}
	return nil
}

func GetItem(i Item, UUID string) error {
	req := &Request{
		Type: "GET",
		Path: i.getRequestPath() + UUID,
		Body: nil,
	}

	rep, err := Send(req)
	if err != nil {
		return err
	}

	if rep.Code != 200 {
		msg, err := i.getResult(ERROR)
		if err != nil {
			return err
		}
		err = json.Unmarshal(rep.Body, &msg)
		if err != nil {
			return err
		}
		err = i.outPut(msg, ERROR)
		if err != nil {
			return err
		}
		return nil
	}

	msg, err := i.getResult(INFO)

	err = json.Unmarshal(rep.Body, &msg)
	if err != nil {
		return err
	}

	err = i.outPut(msg, INFO)
	if err != nil {
		return err
	}
	return nil
}

func GetItemList(i Item) error {
	req := &Request{
		Type: "GET",
		Path: i.getRequestPath() + i.addArgs(),
		Body: nil,
	}

	rep, err := Send(req)
	if err != nil {
		return err
	}

	if rep.Code != 200 {
		msg, err := i.getResult(ERROR)
		if err != nil {
			return err
		}
		err = json.Unmarshal(rep.Body, &msg)
		if err != nil {
			return err
		}
		err = i.outPut(msg, ERROR)
		if err != nil {
			return err
		}
		return nil
	}

	msg, err := i.getResult(LIST)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rep.Body, &msg)
	if err != nil {
		return err
	}

	err = i.outPut(msg, LIST)
	if err != nil {
		return err
	}
	return nil
}

func DeleteItem(i Item, UUID string) error {
	req := &Request{
		Type: "DELETE",
		Path: i.getRequestPath() + UUID,
		Body: nil,
	}

	rep, err := Send(req)
	if err != nil {
		return err
	}

	if rep.Code != 200 {
		msg, err := i.getResult(ERROR)
		if err != nil {
			return err
		}
		err = json.Unmarshal(rep.Body, &msg)
		if err != nil {
			return err
		}
		err = i.outPut(msg, ERROR)
		if err != nil {
			return err
		}
		return nil
	}

	msg, err := i.getResult(JOB)
	if err != nil {
		return err
	}
	err = json.Unmarshal(rep.Body, &msg)
	if err != nil {
		return err
	}

	err = i.outPut(msg, JOB)
	if err != nil {
		return err
	}
	return nil
}

func ErrOutPut(err error) {
	out := os.Stdout
	_, _ = fmt.Fprintf(out, "%s", err)
}
