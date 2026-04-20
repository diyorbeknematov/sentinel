package collector

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/diyorbek/sentinel/agent/internal/config"
	"github.com/diyorbek/sentinel/agent/internal/models"
)

var re = regexp.MustCompile(`(?P<time>\\S+) (?P<level>\\S+) (?P<event>\\S+) (?P<message>.*)`)

func StartAppLogCollector(cfg *config.Config) {
	file, err := os.Open(cfg.AppLog)
	if err != nil {
		fmt.Println("Log file topilmadi:", cfg.AppLog)
		return
	}
	defer file.Close()

	// move to end (tail -f)
	file.Seek(0, 2)

	reader := bufio.NewReader(file)

	fmt.Println("App log kuzatilmoqda:", cfg.AppLog)

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

		appLog, err := parseAppLog(line)
		if err != nil {
			fmt.Println(err)
		}
		// TODO: send log
		fmt.Println("LOG:", appLog)
	}
}

func parseText(line string) (*models.AppLogPayload, error) {
	match := re.FindStringSubmatch(line)
	if match == nil {
		return nil, fmt.Errorf("invalid log")
	}

	result := map[string]string{}
	for i, name := range re.SubexpNames() {
		if i > 0 && name != "" {
			result[name] = match[i]
		}
	}

	// parse time
	logTime, err := time.Parse("02/Jan/2006:15:04:05 -0700", result["time"])
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %v", err)
	}

	return &models.AppLogPayload{
		Event:   result["event"],
		Level:   result["level"],
		Message: result["message"],
		LogTime: logTime,
	}, nil
}

func parseAppLog(line string) (*models.AppLogPayload, error) {

	if isJSON(line) {
		return parseJSON(line)
	}

	return parseText(line)
}

func isJSON(line string) bool {
	return json.Valid([]byte(line))
}

func parseJSON(line string) (*models.AppLogPayload, error) {
	var log models.AppLogPayload
	err := json.Unmarshal([]byte(line), &log)
	return &log, err
}
