package job

import (
	"encoding/json"
	"fmt"
	"math"
	"oracle-go/pkg/entity"
	"oracle-go/pkg/logger"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)


type Job struct {
	Cfg *ini.Section
	DB *gorm.DB
	timeDiff int64
}

// LogEntry represents a log entry in the JSON format
type LogEntry struct {
    Level     string `json:"level"`
    Message   string `json:"message"`
    Timestamp string `json:"timestamp"`
}


func (j *Job) Init(cron *cron.Cron) {
	j.timeDiff = j.Cfg.Key("timeDiffinSecond").MustInt64() * 1000
	cron.AddFunc(j.Cfg.Key("cronExpress").String(), j.SyncPrice)
}

func (j *Job) SyncPrice() {
	logger.Info("Syncing price")

	defer func() {
        logger.Info("Syncing price Done")
    }()

	priceLogPath := j.Cfg.Key("priceLog").String()

	lines, _, err := readLastLines(priceLogPath, 200)
    if err != nil {
        fmt.Printf("Error reading last lines from log file: %v\n", err)
        return
    }


	tokenId := "ton"
	currentTime := time.Now()
	priceMap := make(map[string]float32)

	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		var entry LogEntry
		err := json.Unmarshal([]byte(line), &entry)
		if err != nil {
			panic(err)
		}

		layout := time.RFC3339
		parsedTime, err := time.Parse(layout, entry.Timestamp)

		if err != nil {
			fmt.Println("Error parsing timestamp:", err)
			return
		}
		diff := currentTime.Sub(parsedTime)

		if diff.Milliseconds() > j.timeDiff {
			break
		}

		re := regexp.MustCompile(`^([^|]+)\|([^|]+)\|([^|]+)$`)
		matches := re.FindStringSubmatch(entry.Message)
		if (len(matches) != 4) {
			continue
		}

		value, err := strconv.ParseFloat(matches[3], 32)
        if err != nil {
            panic(value)
        }

		if _, ok := priceMap[matches[1]]; !ok {
			priceMap[matches[1]] = float32(value)
		}
	}
	
	logger.Info("Price map contents:", priceMap)

	if len(priceMap) < 3 {
		logger.Error("No price data found")
		return
	}

	calculatedPrice := calculatePrice(&priceMap)

	result := j.DB.Create(&entity.Price{
		TokenID: tokenId,
		Price:   float64(calculatedPrice),
	})
	logger.Info("Price saved to database")
    if result.Error != nil {
        panic(result.Error)
    }
}

func calculatePrice(priceMap *map[string]float32) float32 {
	values := make([]float32, 0, len(*priceMap))
    for _, value := range *priceMap {
        values = append(values, value)
    }

	min := float32(math.MaxFloat32)
    max := float32(0)
    for _, value := range values {
        if value < min {
            min = value
        }
        if value > max {
            max = value
        }
    }

    // Remove minimum and maximum values
    var filteredValues []float32
    for _, value := range values {
        if value != min && value != max {
            filteredValues = append(filteredValues, value)
        }
    }

	var sum float32
    for _, value := range filteredValues {
        sum += value
    }

    avg := sum / float32(len(filteredValues))

	return avg
}



// readLastLines reads the last `n` lines from a log file and ensures each line is a valid JSON object.
func readLastLines(logPath string, maxLine int) ([]string, int, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	var lines []string
	var buffer []byte
	const bufferSize = 4096
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, 0, err
	}

	fileSize := fileInfo.Size()
	offset := fileSize

	var currentLine strings.Builder

	for offset > 0 {
		if offset < int64(bufferSize) {
			buffer = make([]byte, offset)
		} else {
			buffer = make([]byte, bufferSize)
		}

		offset -= int64(len(buffer))
		if offset < 0 {
			offset = 0
		}

		if _, err := file.Seek(offset, 0); err != nil {
			return nil, 0, err
		}

		n, err := file.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			return nil, 0, err
		}

		// Process the buffer from end to start
		for i := n - 1; i >= 0; i-- {
			if buffer[i] == '\n' {
				if currentLine.Len() > 0 {
					line := reverseString(currentLine.String())
					line = strings.TrimSpace(line)
					if len(line) > 0 {
						lines = append([]string{line}, lines...)
						if len(lines) >= maxLine {
							return lines[:maxLine], len(lines), nil
						}
					}
					currentLine.Reset()
				}
			} else {
				currentLine.WriteByte(buffer[i])
			}
		}
	}

	// Add the last line if there's any remaining
	if currentLine.Len() > 0 {
		line := reverseString(currentLine.String())
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			lines = append([]string{line}, lines...)
		}
	}

	// If the file has fewer lines than `n`, just return the total lines
	if len(lines) < maxLine {
		return lines, len(lines), nil
	}

	return lines[:maxLine], len(lines), nil
}

func reverseString(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		result[len(s)-1-i] = r
	}
	return string(result)
}

