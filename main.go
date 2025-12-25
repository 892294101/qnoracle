package main

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/892294101/qnoracle/log"
	opt "github.com/892294101/qnoracle/options"
	"github.com/sijms/go-ora/v2"
	_ "github.com/sijms/go-ora/v2"
	"github.com/sijms/go-ora/v2/configurations"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ConnectStr struct {
	User     string
	Password string
	Host     string
	Port     int
	Sid      string
}

func main() {
	opts, err := opt.ParseOptions(os.Args)
	if err != nil {
		log.Logvf(log.Always, "error parsing command line options: %s\n", err.Error())
		os.Exit(1)
	}

	var oraOpts configurations.ConnectionConfig
	oraOpts.Timeout = time.Duration(*opts.GlobalCommand.Timeout)
	oraOpts.ConnectTimeout = time.Duration(*opts.GlobalCommand.ConnectTimeout)
	oraOpts.Location = "Asia/Shanghai"

	c, _ := GetConnectString(*opts.GlobalCommand.Url)
	url := go_ora.BuildUrl(c.Host, c.Port, c.Sid, c.User, c.Password, nil)

	client, err := go_ora.NewConnection(url, &oraOpts)
	if err != nil {
		log.Logvf(log.Always, "establish connection error: %s\n", Trim(err.Error()))
		os.Exit(2)
	}
	defer client.Close()

	err = client.Open()
	if err != nil {
		log.Logvf(log.Always, "open connection error: %s\n", Trim(err.Error()))
		os.Exit(3)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*opts.GlobalCommand.Timeout)*time.Second)
	defer cancel()

	rowSet, err := client.QueryContext(ctx, *opts.GlobalCommand.Query, nil)
	if err != nil {
		log.Logvf(log.Always, "query sql error: %v\n", Trim(fmt.Sprintf("%v. sqlText: %v", err.Error(), *opts.GlobalCommand.Query)))
		os.Exit(4)
	}

	defer rowSet.Close()
	var i int

	for {
		row := make([]driver.Value, len(rowSet.Columns()))
		err := rowSet.Next(row)
		if err != nil && err == io.EOF {
			break
		}

		if *opts.GlobalCommand.ConvertJson {
			json, err := ConvertToJson(row, rowSet.Columns())
			if err != nil {
				log.Logvf(log.Always, "convert to json error: %s\n", Trim(err.Error()))
			}
			fmt.Fprintf(os.Stdout, "%s\n", json)
		} else {
			if i == 0 {
				PrintColumn(rowSet.Columns())
			}
			for _, value := range row {
				fmt.Fprintf(os.Stdout, "%v\t", value)
			}
			fmt.Fprintf(os.Stdout, "\n")
		}

		i++
		if i >= *opts.GlobalCommand.Limit {
			break
		}
	}

}

func PrintColumn(columns []string) {
	for _, value := range columns {
		fmt.Fprintf(os.Stdout, "%v\t", value)
	}
	fmt.Fprintf(os.Stdout, "\n")
}
func GetConnectString(str string) (c ConnectStr, err error) {
	re := regexp.MustCompile(`^([^/]+)\/([^@]+)@([^:]+):(\d+)\/([^/]+)$`)
	matches := re.FindStringSubmatch(str)
	if len(matches) > 0 {
		c.User = matches[1]
		c.Password = matches[2]
		c.Host = matches[3]

		portStr := matches[4]
		portInt, err := strconv.Atoi(portStr)

		if err != nil {
			return c, fmt.Errorf("error parsing port (%v): %v", portStr, err)
		}

		c.Port = portInt
		c.Sid = matches[5]
	} else {
		return c, fmt.Errorf("invalid connection string: %s", str)
	}
	return c, nil
}

func Trim(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "  ", " ")
	return s
}

func ConvertToJson(row []driver.Value, colName []string) (string, error) {
	if len(row) != len(colName) {
		return "", fmt.Errorf("length of the slice for column value, column name, and column type does not match")
	}

	data := make(map[string]interface{}, len(row))
	for i, value := range row {
		data[colName[i]] = value
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("json encoding failed:", err)
	}
	return string(jsonData), nil
}
