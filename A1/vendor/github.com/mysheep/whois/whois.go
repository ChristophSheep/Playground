package whois

import (
	"bytes"
	"os/exec"
	"strings"
)

var cache = map[string]string{
	"0.0.0.0": "empty",
}

func Get(ip string) (map[string]string, error) {
	str, err := get(ip)
	if err != nil {
		return nil, err
	}
	return parseMap(str), nil
}

func parseMap(allLines string) map[string]string {

	var name_values = map[string]string{}

	lines := filterComments(strings.Split(allLines, "\n"))

	for _, line := range lines {
		xs := strings.Split(line, ":")

		if len(xs) == 2 {
			name := strings.TrimSpace(xs[0])
			value := strings.TrimSpace(xs[1])

			name_values[name] = value
		}
	}

	return name_values

}

func filter(line string) bool {

	// % A comment
	// # Another comment
	//
	if strings.HasPrefix(line, "%") || strings.HasPrefix(line, "#") {
		return false
	}

	// Empty line
	//
	if len(line) == 0 {
		return false
	}

	return true
}

func filterComments(lines []string) []string {

	// Create a zero-length slice with the same underlying array
	tmp := lines[:0]

	for _, line := range lines {

		if filter(line) {
			tmp = append(tmp, line)
		}

	}

	return tmp
}

func get(ip string) (string, error) {

	if str, found := cache[ip]; found {
		return str, nil
	}

	cmd := exec.Command("whois", ip)

	var out bytes.Buffer
	var str = ""

	cmd.Stdout = &out
	err := cmd.Run()

	if err == nil {
		str = out.String()
		cache[ip] = str
	}

	return str, err
}

/*
func main() {
	str, err := whois("199.232.18.133")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Whois: %s\n", str)
}
*/
