package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/shadowabi/AutoDomain_rebuild/utils/Error"
	"os"
	"sort"
	"strings"
)

func WriteToFile(writeResultList []string, output string) {

	file, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	Error.HandlePanic(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	sort.Strings(writeResultList)

	if strings.HasSuffix(output, ".json") {
		jsonData := OutputJson(writeResultList)
		fmt.Fprintln(writer, string(jsonData))
	} else if len(writeResultList) != 0 {
		for _, i := range writeResultList {
			fmt.Fprintln(writer, i)
		}
	}
}

func OutputJson(writeResultList []string) (jsonData []byte) {
	data := make(map[string][]string)
	data["url"] = writeResultList

	jsonData, err := json.Marshal(data)
	Error.HandlePanic(err)
	return jsonData
}
