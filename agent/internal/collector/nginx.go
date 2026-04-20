package collector

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/diyorbek/sentinel/agent/internal/models"
)

var nginxRegex = regexp.MustCompile(
	`(?P<ip>\\S+) - - \\[(?P<time>.*?)\\] \"(?P<method>\\S+) (?P<path>\\S+) \\S+\" (?P<status>\\d+) (?P<bytes>\\d+) \"[^\"]*\" \"(?P<ua>.*?)\"`,
)

func StartNginxLogCollector(cfg *config.Config) {
	file, err := os.Open(cfg.NginxLog)
	if err != nil {
		fmt.Println("Nginx log topilmadi:", cfg.NginxLog)
		return
	}
	defer file.Close()

	// tail -f (start from end)
	file.Seek(0, 2)

	reader := bufio.NewReader(file)

	fmt.Println("Nginx log kuzatilmoqda:", cfg.NginxLog)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Logni o'qishda xatolik:", err)
			time.Sleep(200 * time.Millisecond)
			continue
		}

		if len(line) == 0 {
			continue
		}

		nginxLog, err := parseNginxLog(line)
		if err != nil {
			fmt.Println("Error in parsing Log:", err)
		}
		// TODO: parse qilamiz keyin
		fmt.Println("NGINX:", nginxLog)
	}
}

func parseNginxLog(line string) (*models.NginxLogPayload, error) {
	match := nginxRegex.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("invalid nginx log format")
	}

	result := make(map[string]string)
	for i, name := range nginxRegex.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	// parse status
	status, _ := strconv.Atoi(result["status"])

	// parse bytes
	bytes, _ := strconv.Atoi(result["bytes"])

	// parse time
	logTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", result["time"])
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %v", err)
	}

	return &models.NginxLogPayload{
		IP:        result["ip"],
		Method:    result["method"],
		Path:      result["path"],
		Status:    status,
		Bytes:     bytes,
		UserAgent: result["ua"],
		LogTime:   logTime,
	}, nil
}
